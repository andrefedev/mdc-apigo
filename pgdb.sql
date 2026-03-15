CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE
TEXT SEARCH CONFIGURATION es ( COPY = pg_catalog.spanish );
ALTER
TEXT SEARCH CONFIGURATION es ALTER
MAPPING FOR hword, hword_part, word WITH unaccent, spanish_stem;

ALTER TABLE users
    ADD COLUMN search_tsvector tsvector GENERATED ALWAYS AS (setweight(to_tsvector('es', coalesce(name, '')), 'A') ||
                                                             setweight(to_tsvector('simple', coalesce(phone, '')), 'B')) STORED;

CREATE INDEX IF NOT EXISTS users_search_gin ON users USING gin (tsvector_search);

CREATE INDEX IF NOT EXISTS users_search_gin ON users USING gin (tsvector_search);

CREATE INDEX IF NOT EXISTS idx_users_ref_date_joined_desc ON users (date_joined DESC, ref DESC);



CREATE
EXTENSION IF NOT EXISTS postgis;

-- Una sola dirección por defecto por usuario (ignorando soft delete)
CREATE UNIQUE INDEX IF NOT EXISTS users_addrs_is_default_uindex ON users_addrs(user_id) WHERE is_default;

-- Evita duplicar la misma place_id para el mismo usuario
CREATE UNIQUE INDEX IF NOT EXISTS ux_users_addrs_user_place ON users_addrs(user_id, google_place_id) WHERE deleted_at IS NULL AND google_place_id IS NOT NULL;


ALTER TABLE orders
    ADD COLUMN total_price integer GENERATED ALWAYS AS (base_price - disc_price) STORED;

-- ORDERS
DROP TYPE order_status_enum;
CREATE TYPE payment_method_enum AS ENUM ('bank','cash','breb', 'qrcode');
CREATE TYPE payment_status_enum AS ENUM ('pending','refunded', 'authorized');
CREATE TYPE order_status_enum AS ENUM ('draft', 'pending', 'acepted', 'canceled','dispatched','successfully');

ALTER TABLE orders_lines ADD COLUMN total_price integer GENERATED ALWAYS AS (base_price - disc_price) STORED;


ALTER TABLE orders
    ADD COLUMN order_tsvector tsvector GENERATED ALWAYS AS (setweight(to_tsvector('es', coalesce(number, '')), 'A') ||
                                                             setweight(to_tsvector('simple', coalesce(phone, '')), 'B')) STORED;

INSERT INTO deliveries_slots (label, start_min, end_min)
VALUES ('07:00 – 09:00', 7 * 60, 9 * 60),
       ('09:00 - 11:00', 9 * 60, 11 * 60),
       ('11:00 – 13:00', 11 * 60, 13 * 60),
       ('13:00 - 15:00', 13 * 60, 15 * 60),
       ('15:00 - 17:00', 15 * 60, 17 * 60),
       ('17:00 – 19:00', 17 * 60, 19 * 60);


--

CREATE OR REPLACE FUNCTION seed_deliveries_year(
    p_year         int,
    p_capacity     int    DEFAULT 20,               -- capacidad por slot
    p_time_zone    text   DEFAULT 'America/Bogota', -- hora local
    p_only_if_open boolean DEFAULT false,           -- si true, solo días abiertos
    p_reserved0    int    DEFAULT 0                 -- reservado inicial
)
    RETURNS TABLE(day_date date, inserted_slots int)  -- <— renombrado para evitar ambigüedad
    LANGUAGE plpgsql
AS $$
DECLARE
v_start date;
    v_end   date;
BEGIN
    IF p_year IS NULL OR p_year < 1 THEN
        RAISE EXCEPTION 'Año inválido: %', p_year;
END IF;

    v_start := make_date(p_year, 1, 1);
    v_end   := (v_start + INTERVAL '1 year - 1 day')::date;

    -- 1) Crear días del año (idempotente)
INSERT INTO deliveries_days (work_date, is_open)
SELECT gs.d::date, TRUE
FROM generate_series(v_start::timestamp, v_end::timestamp, INTERVAL '1 day') AS gs(d)
    ON CONFLICT (work_date) DO NOTHING;

-- 2) Insertar franjas fijas por cada día (idempotente)
RETURN QUERY
    WITH days AS (
            SELECT dd.id AS wid, dd.work_date, dd.is_open
            FROM deliveries_days AS dd
            WHERE dd.work_date BETWEEN v_start AND v_end
              AND (NOT p_only_if_open OR dd.is_open)
        ),
             slots AS (
                 VALUES
                     ('07:00-10:00','07:00'::time,'10:00'::time),
                     ('11:00-13:00','11:00'::time,'13:00'::time),
                     ('14:00-16:00','14:00'::time,'16:00'::time),
                     ('17:00-19:00','17:00'::time,'19:00'::time)
             ),
             ins AS (
                 INSERT INTO deliveries_days_slots (wid, code, from_time, until_time, capacity, reserved)
                     SELECT
                         d.wid,
                         s.column1::varchar(128) AS code,
                         ((d.work_date::timestamp + s.column2::time) AT TIME ZONE p_time_zone) AS from_time,
                         (
                             (
                                 d.work_date::timestamp
                                     + s.column3::time
                                     + CASE WHEN s.column3::time <= s.column2::time THEN INTERVAL '1 day' ELSE INTERVAL '0 day' END
                                 ) AT TIME ZONE p_time_zone
                             ) AS until_time,
                         p_capacity,
                         p_reserved0
                     FROM days d
                              CROSS JOIN slots s
                     ON CONFLICT (wid, code) DO NOTHING
                     RETURNING wid
             )
SELECT d.work_date AS day_date, COALESCE(COUNT(i.wid), 0)::int AS inserted_slots
FROM (SELECT dd.id, dd.work_date FROM deliveries_days dd
      WHERE dd.work_date BETWEEN v_start AND v_end) d
         LEFT JOIN ins i ON i.wid = d.id
GROUP BY d.work_date
ORDER BY d.work_date;
END;
$$;



-- Devuelve franjas para una fecha dada, marcando si son seleccionables hoy
BEGIN;
SELECT pg_advisory_xact_lock(987654, 2026);
SELECT day_date AS work_date, inserted_slots
FROM seed_deliveries_year(2026);
COMMIT;


CREATE OR REPLACE FUNCTION order_line_check_before_delete()
    RETURNS TRIGGER AS $$
DECLARE
    statusx order_status_enum;
BEGIN
    SELECT status INTO statusx FROM orders WHERE id = OLD.oid;

    IF statusx IN ('succe', 'CANCELADO') THEN
        RAISE EXCEPTION 'No se puede eliminar un detalle de una orden con estado %', estado;
    END IF;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;


SELECT
    o.id,
    -- id de cluster (NULL = ruido / punto suelto)
    ST_ClusterDBSCAN(
    ST_Transform(addr.geom, 32618),  -- metros
    eps := 3000,                     -- ~1 km de vecindad
    minpoints := 2                   -- densidad mínima
                    ) OVER (PARTITION BY days.work_date) AS cluster_id,
    addr.neighb, u.name
FROM orders AS o
-- JOIN CLIENT --
         JOIN users AS u ON u.id = o.uid
-- JOIN CLIENT --
         JOIN users_addrs AS addr ON addr.id = o.shid
    -- JOIN DELIVERY SLOT --
         JOIN deliveries_slots AS slot ON slot.id = o.sloid
-- JOIN DELIVERY DAYS --
         JOIN deliveries_days AS days ON days.id = slot.wid
WHERE days.work_date = CURRENT_DATE - 2  AND addr.geom IS NOT NULL;


SELECT
    'BEGIN:VCARD'                   || E'\r\n' ||
    'VERSION:3.0'                   || E'\r\n' ||
    'FN: 2-' || name                   || E'\r\n' ||
    'TEL;TYPE=CELL:+57' || phone    || E'\r\n' ||
    'END:VCARD'                     || E'\r\n'
        AS vcard
FROM users where is_active = TRUE LIMIT 250 OFFSET 251;

SELECT
    o.id,
    ST_ClusterDBSCAN(ST_Transform(addr.geom, 32618), eps := 3000, minpoints := 2) OVER (PARTITION BY days.work_date) AS cluster_id,
    addr.neighb, u.name
FROM orders AS o
-- JOIN CLIENT --
         JOIN users AS u ON u.id = o.uid
-- JOIN CLIENT --
         JOIN users_addrs AS addr ON addr.id = o.shid
    -- JOIN DELIVERY SLOT --
         JOIN deliveries_slots AS slot ON slot.id = o.sloid
-- JOIN DELIVERY DAYS --
         JOIN deliveries_days AS days ON days.id = slot.wid
WHERE days.work_date = CURRENT_DATE  AND addr.geom IS NOT NULL ORDER BY cluster_id NULLS LAST;


-- SELECT
-- SUM(l.quantity) AS total_vendido
-- FROM orders_lines l
-- JOIN orders o ON o.id = l.oid
-- WHERE l.pid = '68458d20-d0a1-42c0-83a9-f22e9e1e3bdc' AND o.status = 'acepted';

SELECT
    SUM(l.quantity) AS total_vendido
FROM orders_lines l
         JOIN orders o ON o.id = l.oid
         JOIN users_addrs AS addr ON addr.id = o.shid
         JOIN deliveries_slots AS slot ON slot.id = o.sloid
         JOIN deliveries_days AS days ON days.id = slot.wid
WHERE days.work_date = '2025-12-02' AND l.pid = '86510afc-9747-4d64-a43e-7aa67bc6c064';

SELECT
    COUNT(DISTINCT o.id)
FROM orders_lines l
         JOIN orders o ON o.id = l.oid
         JOIN users_addrs AS addr ON addr.id = o.shid
         JOIN deliveries_slots AS slot ON slot.id = o.sloid
         JOIN deliveries_days AS days ON days.id = slot.wid
WHERE days.work_date = '2025-11-30' AND o.status = 'acepted';

SELECT
    SUM(l.total_price) AS total_vendido
FROM orders_lines l
         JOIN orders o ON o.id = l.oid
         JOIN users_addrs AS addr ON addr.id = o.shid
         JOIN deliveries_slots AS slot ON slot.id = o.sloid
         JOIN deliveries_days AS days ON days.id = slot.wid
WHERE days.work_date = '2025-11-25' AND o.status = 'acepted';

-- A   bc00631e-5061-4fd7-8db0-e535885632e9
-- AA  86510afc-9747-4d64-a43e-7aa67bc6c064
-- AAA 68458d20-d0a1-42c0-83a9-f22e9e1e3bdc

SELECT * FROM USERS;