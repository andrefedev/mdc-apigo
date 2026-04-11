package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"apigo/internal/modules/postgres"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Pgdb
}

func NewRepository(db *postgres.Pgdb) *Repository {
	return &Repository{db: db}
}

func (r Repository) CodeInsert(ctx context.Context, data *CodeInsertData) (string, error) {
	op := "App.Repository.CodeInsert"
	qry := `INSERT INTO users_codes_original (code, phone) VALUES (@code, @phone) RETURNING id;`

	var ref string
	if err := r.db.QueryRow(
		ctx,
		qry,
		pgx.NamedArgs{
			"code":  data.Code,
			"phone": data.Phone,
		},
	).Scan(&ref); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ref, nil
}

func (r Repository) CodeDelete(ctx context.Context, ref string) (int64, error) {
	op := "App.Repository.CodeDelete"
	qry := `DELETE FROM users_codes_original WHERE id = $1`

	res, err := r.db.Exec(ctx, qry, ref)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.RowsAffected(), nil
}

func (r Repository) CodeSelect(ctx context.Context, ref string) (*Code, error) {
	const op = "App.Repository.CodeSelect"
	qry := `SELECT id, code, phone, date_created, date_expired FROM users_codes_original WHERE id = $1`

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[Code])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapCodeNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (r Repository) SessionInsert(ctx context.Context, data *SessionInsertData) (string, error) {
	const op = "App.Repository.SessionInsert"

	qry := `
	INSERT INTO auth_sessions (uid, token_hash)
	VALUES (@uid, @token_hash)
	RETURNING id
	`

	var ref string
	if err := r.db.QueryRow(
		ctx,
		qry,
		pgx.NamedArgs{
			"uid":        data.UserRef,
			"token_hash": data.TokenHash,
			// "date_expired": data.DateExpired,
		},
	).Scan(&ref); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ref, nil
}

func (r Repository) SessionSelect(ctx context.Context, ref string) (*Session, error) {
	const op = "App.Repository.SessionSelect"

	qry := `
	SELECT id, uid, token_hash, last_used_at, date_expired, date_created, date_revoked
	FROM auth_sessions
	WHERE id = $1
	`

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapSessionNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return new(result), nil
}

func (r Repository) SessionSelectByToken(ctx context.Context, token string) (*Session, error) {
	const op = "App.Repository.SessionSelectByToken"

	qry := `
	SELECT s.id, s.uid, s.token_hash, s.date_expired, s.date_revoked, s.date_created,
	u.is_super AS is_super, u.is_staff AS is_staff, u.is_active AS is_active -- ok --
	FROM auth_sessions AS s JOIN users AS u ON u.id = s.uid WHERE token_hash = $1 LIMIT 1
	`

	rows, err := r.db.Query(ctx, qry, token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapSessionNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return new(result), nil
}

// USER___

func (r Repository) UserInsert(ctx context.Context, data *UserInsertData) (string, error) {
	const op = "App.Repository.UserInsert"

	qry := `INSERT INTO users (name, phone, is_super, is_staff, is_active) VALUES (@name, @phone, @is_super, @is_staff, @is_active) RETURNING id;`

	var ref string

	if err := r.db.QueryRow(
		ctx,
		qry,
		pgx.NamedArgs{
			"name":      data.Name,
			"phone":     data.Phone,
			"is_super":  data.IsSuper,
			"is_staff":  data.IsStaff,
			"is_active": data.IsActive,
		},
	).Scan(&ref); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ref, nil
}

func (r Repository) UserUpdate(ctx context.Context, ref string, paths []string, data *UserUpdateData) (int64, error) {
	const op = "App.Repository.UserUpdate"

	var values []any
	var clauses []string
	for _, path := range paths {
		switch path {
		case "name":
			values = append(values, data.Name)
			clauses = append(clauses, fmt.Sprintf(`name = $%d`, len(values)))
		case "phone":
			values = append(values, data.Phone)
			clauses = append(clauses, fmt.Sprintf(`phone = $%d`, len(values)))
		case "is_super":
			values = append(values, data.IsSuper)
			clauses = append(clauses, fmt.Sprintf(`is_super = $%d`, len(values)))
		case "is_staff":
			values = append(values, data.IsStaff)
			clauses = append(clauses, fmt.Sprintf(`is_staff = $%d`, len(values)))
		case "is_active":
			values = append(values, data.IsActive)
			clauses = append(clauses, fmt.Sprintf(`is_active = $%d`, len(values)))
		case "last_login":
			values = append(values, data.LastLogin)
			clauses = append(clauses, fmt.Sprintf(`last_login = $%d`, len(values)))
		}
	}
	if len(clauses) == 0 {
		return 0, nil // nada que actualizar
	}

	// Update sets
	xquery := "UPDATE users"
	xquery += " SET " + strings.Join(clauses, ", ")

	// Where
	values = append(values, ref)
	xquery += fmt.Sprintf(" WHERE id = $%d", len(values))

	// exec
	res, err := r.db.Exec(ctx, xquery, values...)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.RowsAffected(), nil
}

func (r Repository) UserSelect(ctx context.Context, ref string) (*User, error) {
	const op = "App.Repository.UserSelect"

	qry := `
	SELECT
	id, name, phone, is_staff, is_super, is_active, last_login, date_joined
	FROM users WHERE id = $1
	`

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapUserNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return new(res), nil
}

func (r Repository) UserRefByPhone(ctx context.Context, phone string) (string, error) {
	const op = "App.Repository.UserRefByPhone"

	qry := `SELECT id FROM users WHERE phone = $1 LIMIT 1`

	var userRef string
	if err := r.db.QueryRow(ctx, qry, phone).Scan(&userRef); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, WrapUserNotFound(err))
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userRef, nil
}

func (r Repository) UserSelectByPhone(ctx context.Context, phone string) (*User, error) {
	const op = "App.Repository.UserSelectByPhone"

	qry := `
	SELECT
	id, name, phone, is_staff, is_super, is_active, last_login, date_joined
	FROM users WHERE phone = $1
	`

	rows, err := r.db.Query(ctx, qry, phone)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapUserNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return new(res), nil
}

func (r Repository) UserSelectAll(ctx context.Context, filter *UserFilterData, paging *UserPagingData) ([]*User, error) {
	const op = "App.Repository.UserSelectAll"

	qry := `
	SELECT
	u.id, u.name, u.phone, u.is_super, u.is_staff, u.is_active, u.last_login, u.date_joined
	FROM users AS u
	`

	// # BEGIN FILTER #
	var values []any
	var clauses []string

	if filter != nil {
		// FLAT_QUERY
		if filter.FlatQuery != nil {
			q := strings.TrimSpace(*filter.FlatQuery)
			if q != "" {
				_, err := strconv.ParseInt(q, 10, 64)
				if err == nil {
					values = append(values, "%"+q+"%")
					clauses = append(clauses, fmt.Sprintf(`u.phone ILIKE $%d`, len(values)))
				} else {
					values = append(values, "%"+q+"%")
					clauses = append(clauses, fmt.Sprintf(`unaccent(lower(u.name)) ILIKE unaccent(lower($%d))`, len(values)))
				}
			}
		}

		// IS_SUPER
		if filter.IsSuper != nil {
			values = append(values, *filter.IsSuper)
			clauses = append(clauses, fmt.Sprintf(`u.is_super = $%d`, len(values)))
		}

		// IS_STAFF
		if filter.IsStaff != nil {
			values = append(values, *filter.IsStaff)
			clauses = append(clauses, fmt.Sprintf(`u.is_staff = $%d`, len(values)))
		}

		// IS_ACTIVE
		if filter.IsActive != nil {
			values = append(values, *filter.IsActive)
			clauses = append(clauses, fmt.Sprintf(`u.is_active = $%d`, len(values)))
		}
	}

	// # CLASUSES SEP #
	if len(clauses) > 0 {
		qry += " WHERE " + strings.Join(clauses, " AND ")
	}

	// # ORDER BY #
	qry += " ORDER BY u.date_joined DESC, u.id DESC"

	// # PAGINATION #
	if paging != nil {
		qry += fmt.Sprintf(` LIMIT %d `, paging.Limit)
		qry += fmt.Sprintf(` OFFSET %d `, paging.Offset)
	}

	// # END DEFAULT FILTER #

	rows, err := r.db.Query(ctx, qry, values...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[User])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

// USER_ADDR__

func (r Repository) UserAddrInsert(ctx context.Context, uid string, data *UserAddrInsertData) (string, error) {
	const op = "App.Repository.UserAddrInsert"

	qry := `
	INSERT INTO users_addrs (uid, pid, lat, lng, name, cmna, route, street, neighb, locality, sublocal, address1, address2)
	VALUES (@uid, @pid, @lat, @lng, @name, @cmna, @route, @street, @neighb, @locality, @sublocal, @address1, @address2) RETURNING id;`

	var ref string
	if err := r.db.QueryRow(
		ctx,
		qry,
		pgx.NamedArgs{
			"uid":      uid,
			"pid":      data.Pid,
			"lat":      data.Lat,
			"lng":      data.Lng,
			"name":     data.Name,
			"cmna":     data.Cmna,
			"route":    data.Route,
			"street":   data.Street,
			"neighb":   data.Neighb,
			"locality": data.Locality,
			"sublocal": data.Sublocal,
			"address1": data.Address1,
			"address2": data.Address2,
		},
	).Scan(&ref); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ref, nil
}

func (r Repository) UserAddrUpdate(ctx context.Context, ref string, paths []string, data *UserAddrUpdateData) (int64, error) {
	const op = "App.Repository.UserAddrUpdate"

	var values []any
	var clauses []string

	for _, path := range paths {
		switch path {
		case "lat":
			values = append(values, data.Lat)
			clauses = append(clauses, fmt.Sprintf(`lat = $%d`, len(values)))
		case "lng":
			values = append(values, data.Lng)
			clauses = append(clauses, fmt.Sprintf(`lng = $%d`, len(values)))
		case "cmna":
			values = append(values, data.Cmna)
			clauses = append(clauses, fmt.Sprintf(`cmna = $%d`, len(values)))
		case "route":
			values = append(values, data.Route)
			clauses = append(clauses, fmt.Sprintf(`route = $%d`, len(values)))
		case "street":
			values = append(values, data.Street)
			clauses = append(clauses, fmt.Sprintf(`street = $%d`, len(values)))
		case "neighb":
			values = append(values, data.Neighb)
			clauses = append(clauses, fmt.Sprintf(`neighb = $%d`, len(values)))
		case "locality":
			values = append(values, data.Locality)
			clauses = append(clauses, fmt.Sprintf(`locality = $%d`, len(values)))
		case "sublocal":
			values = append(values, data.Sublocal)
			clauses = append(clauses, fmt.Sprintf(`sublocal = $%d`, len(values)))
		case "address1":
			values = append(values, data.Address1)
			clauses = append(clauses, fmt.Sprintf(`address1 = $%d`, len(values)))
		case "address2":
			values = append(values, data.Address2)
			clauses = append(clauses, fmt.Sprintf(`address2 = $%d`, len(values)))
		case "is_default":
			values = append(values, data.IsDefault)
			clauses = append(clauses, fmt.Sprintf(`is_default = $%d`, len(values)))
		}
	}

	if len(clauses) == 0 {
		return 0, nil // nada que actualizar
	}

	// Update sets
	xquery := "UPDATE users_addrs"
	xquery += " SET " + strings.Join(clauses, ", ")

	// Where
	values = append(values, ref)
	xquery += fmt.Sprintf(" WHERE id = $%d", len(values))

	// exec
	res, err := r.db.Exec(ctx, xquery, values...)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.RowsAffected(), nil
}

func (r Repository) UserAddrSelect(ctx context.Context, ref string) (*UserAddr, error) {
	const op = "App.Repository.UserAddrSelect"

	qry := `
	SELECT
	id, pid, lat, lng, name, cmna, route, street, neighb, locality, sublocal, address1, address2, is_default, date_created, date_updated
	FROM users_addrs WHERE id = $1
	`

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[UserAddr])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapUserAddrNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (r Repository) UserAddrSelectAll(ctx context.Context, uid string) ([]*UserAddr, error) {
	const op = "App.Repository.UserAddrSelectAll"

	qry := `
	SELECT
	id, pid, lat, lng, name, cmna, route, street, neighb, locality, sublocal, address1, address2, is_default, date_created, date_updated
	FROM users_addrs WHERE uid = $1 ORDER BY is_default DESC, date_created DESC
	`

	rows, err := r.db.Query(ctx, qry, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[UserAddr])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return results, nil
}

// CATLG__

// GENRE__

func (r Repository) GenreSelect(ctx context.Context, ref string) (*Genre, error) {
	const op = "App.Repository.GenreSelect"

	qry := `
	SELECT
	g.id, g.name, g.descr, g.imurl, g.display, g.is_public, g.date_created
	-- FROM GENRES && WHERE --
	FROM genres AS g WHERE g.id = $1
	`

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[Genre])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapGenreNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (r Repository) GenreSelectAll(ctx context.Context) ([]*Genre, error) {
	const op = "App.Repository.GenreSelectAll"

	qry := `
	SELECT
	g.id, g.name, g.descr, g.imurl, g.display, g.is_public, g.date_created
	-- COUNT(g.id) AS product_count -- COUNT PRODUCTS --
	-- FROM GENRES --
	FROM genres AS g
	-- LEFT JOIN PRODUCTS --
	-- LEFT JOIN products AS p ON p.gid = g.id
	-- GROUP BY --
	-- GROUP BY g.id, g.name, g.imurl, g.descr, g.display, g.is_public, g.date_created
	-- ORDER BY DESC --
	ORDER BY g.display, g.date_created DESC;
	`

	rows, err := r.db.Query(ctx, qry)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[Genre])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return results, nil
}

// ITEMS__

func (r Repository) ProductSelect(ctx context.Context, ref string, forUpdate bool) (*Product, error) {
	const op = "App.Repository.ProductSelect"

	qry := `
	SELECT
	p.id, p.upc, p.code, p.name, p.descr, p.imurl, p.display, p.weight, p.unitype, p.quantity, p.is_active,
	p.is_public, p.cost_price, p.base_price, p.num_in_alloc, p.num_in_stock, p.date_created, p.date_updated,
	-- GENRE --
	g.id AS genre_ref, g.name AS genre_name, g.descr AS genre_descr, g.imurl AS genre_imurl,
	g.display AS genre_display, g.is_public AS genre_is_public, g.date_created AS genre_date_created
	-- FROM --
	FROM products AS p
	-- JOIN GENRE --
	JOIN genres AS g ON g.id = p.gid
	-- WHERE --
	WHERE p.id = $1
	`
	if forUpdate {
		qry += " FOR UPDATE"
	}

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[ProductRaw])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapProductNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result.ToProduct(), nil
}

func (r Repository) ProductSelectAll(ctx context.Context, filter *ProductFilterData) ([]*Product, error) {
	const op = "App.Repository.ProductSelectAll"

	qry := `
	SELECT
	p.id, p.upc, p.code, p.name, p.descr, p.imurl, p.display, p.weight, p.unitype, p.quantity, p.is_active,
	p.is_public, p.cost_price, p.base_price, p.num_in_alloc, p.num_in_stock, p.date_created, p.date_updated,
	-- GENRE --
	g.id AS genre_ref, g.name AS genre_name, g.descr AS genre_descr, g.imurl AS genre_imurl,
	g.display AS genre_display, g.is_public AS genre_is_public, g.date_created AS genre_date_created
	-- FROM --
	FROM products AS p
	-- JOIN GENRE --
	JOIN genres AS g ON g.id = p.gid
	`

	var values []any
	var clauses []string

	if filter != nil {
		if filter.Query != nil {
			values = append(values, *filter.Query)
			clauses = append(clauses, fmt.Sprintf(
				`(unaccent(p.name) ILIKE unaccent('%%' || $%d || '%%') OR CAST(p.code AS TEXT) ILIKE '%%' || $%d || '%%')`,
				len(values), len(values),
			))
		}
		if filter.Genre != nil {
			values = append(values, *filter.Genre)
			clauses = append(clauses, fmt.Sprintf(`p.gid = $%d`, len(values)))
		}
		if filter.IsActive != nil {
			values = append(values, *filter.IsActive)
			clauses = append(clauses, fmt.Sprintf(`p.is_active = $%d`, len(values)))
		}
		if filter.IsPublic != nil {
			values = append(values, *filter.IsPublic)
			clauses = append(clauses, fmt.Sprintf(`p.is_public = $%d`, len(values)))
		}
	}

	if len(clauses) > 0 {
		qry += " WHERE " + strings.Join(clauses, " AND ")
	}

	qry += " ORDER BY g.display, p.display, p.date_created DESC"

	rows, err := r.db.Query(ctx, qry, values...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	raws, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[ProductRaw])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results := make([]*Product, 0, len(raws))
	for _, raw := range raws {
		results = append(results, raw.ToProduct())
	}

	return results, nil
}

// SALES__

// ORDER__

func (r Repository) OrderInsert(ctx context.Context, data *OrderInsertData) (string, error) {
	const op = "App.Repository.OrderInsert"

	query := `INSERT INTO orders (
	uid, shid, sloid, status, payment_status, payment_method) VALUES (
	@uid, @shid, @sloid, @status, @payment_status, @payment_method) RETURNING id;`

	var ref string
	if err := r.db.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"uid":            data.User,
			"shid":           data.Addr,
			"sloid":          data.Slot,
			"status":         data.Status,
			"payment_status": data.PaymentStatus,
			"payment_method": data.PaymentMethod,
		},
	).Scan(&ref); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ref, nil

}

func (r Repository) OrderUpdate(ctx context.Context, ref string, paths []string, data *OrderUpdateData) (int64, error) {
	const op = "App.Repository.OrderUpdate"

	var values []any
	var clauses []string

	for _, path := range paths {
		switch path {
		case "addr":
			values = append(values, data.Addr)
			clauses = append(clauses, fmt.Sprintf(`shid = $%d`, len(values)))
		case "slot":
			values = append(values, data.Slot)
			clauses = append(clauses, fmt.Sprintf(`sloid = $%d`, len(values)))
		case "status":
			values = append(values, data.Status)
			clauses = append(clauses, fmt.Sprintf(`status = $%d`, len(values)))
		case "payment_status":
			values = append(values, data.PaymentStatus)
			clauses = append(clauses, fmt.Sprintf(`payment_status = $%d`, len(values)))
		case "payment_method":
			values = append(values, data.PaymentMethod)
			clauses = append(clauses, fmt.Sprintf(`payment_method = $%d`, len(values)))
		}
	}
	if len(clauses) == 0 {
		return 0, nil
	}

	// ADD UPDATE NOW()
	clauses = append(clauses, "date_updated = NOW()")

	// Update sets
	xquery := "UPDATE orders"
	xquery += " SET " + strings.Join(clauses, ", ")

	// Where
	values = append(values, ref)
	xquery += fmt.Sprintf(" WHERE id = $%d", len(values))

	// exec
	res, err := r.db.Exec(ctx, xquery, values...)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.RowsAffected(), nil
}

func (r Repository) OrderDelete(ctx context.Context, ref string) (int64, error) {
	const op = "App.Repository.OrderDelete"

	qry := `DELETE FROM orders WHERE id = $1 AND status = 'pending';`

	tag, err := r.db.Exec(ctx, qry, ref)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return tag.RowsAffected(), nil
}

func (r Repository) OrderSelect(ctx context.Context, ref string, forUpdate bool) (*Order, error) {
	const op = "App.Repository.OrderSelect"

	qry := `
	SELECT
	o.id, o.number, o.status, o.date_created, o.date_updated, o.payment_status, o.payment_method,
	-- CLIENT --
	u.id AS user_ref, u.name AS user_name, u.phone AS user_phone,
	-- ADDR --
	addr.id AS addr_ref, addr.lat AS addr_lat, addr.lng AS addr_lng, addr.name AS addr_name, addr.cmna AS addr_cmna, addr.route AS addr_route,
	addr.street AS addr_street, addr.neighb AS addr_neighb, addr.locality AS addr_locality, addr.sublocal AS addr_sublocal, addr.address1 AS addr_address1, addr.address2 AS addr_address2,
	-- SLOT --
	slot.id AS slot_ref, slot.kind AS slot_kind, slot.note AS slot_note, slot.wday AS slot_wday, slot.is_open AS slot_is_open, slot.date_created AS slot_date_created, slot.date_updated AS slot_date_updated,
	(extract(hour from slot.cutoff_min)::int * 60 + extract(minute from slot.cutoff_min)::int) AS slot_cutoff_min,
	(extract(hour from slot.delivery_start)::int * 60 + extract(minute from slot.delivery_start)::int) AS slot_delivery_start,
	(extract(hour from slot.delivery_until)::int * 60 + extract(minute from slot.delivery_until)::int) AS slot_delivery_until
	-- FROM --
	FROM orders AS o
	-- JOIN CLIENT --
	JOIN users AS u ON u.id = o.uid
	-- JOIN SHIPPING --
	JOIN users_addrs AS addr ON addr.id = o.shid
	-- JOIN SLOT --
	JOIN deliveries_slots AS slot ON slot.id = o.sloid
	-- WHERE --
	WHERE o.id = $1
	`

	// # FOR UPDATE #
	if forUpdate {
		qry += " FOR UPDATE"
	}

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[OrderRaw])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapOrderNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result.ToOrder(), nil
}

func (r Repository) OrderSelectAll(ctx context.Context, filter *OrderFilterData, paging *OrderPagingData) ([]*Order, error) {
	const op = "App.Repository.OrderSelectAll"

	qry := `
	SELECT
	o.id, o.number, o.status, o.date_created, o.date_updated, o.payment_status, o.payment_method,
	-- USER --
	u.id AS user_ref, u.name AS user_name, u.phone AS user_phone,
	-- ADDR --
	addr.id AS addr_ref, addr.lat AS addr_lat, addr.lng AS addr_lng, addr.name AS addr_name, addr.cmna AS addr_cmna, addr.route AS addr_route,
	addr.street AS addr_street, addr.neighb AS addr_neighb, addr.locality AS addr_locality, addr.sublocal AS addr_sublocal, addr.address1 AS addr_address1, addr.address2 AS addr_address2,
	-- SLOT --
	slot.id AS slot_ref, slot.kind AS slot_kind, slot.note AS slot_note, slot.wday AS slot_wday, slot.is_open AS slot_is_open, slot.date_created AS slot_date_created, slot.date_updated AS slot_date_updated,
	(extract(hour from slot.cutoff_min)::int * 60 + extract(minute from slot.cutoff_min)::int) AS slot_cutoff_min,
	(extract(hour from slot.delivery_start)::int * 60 + extract(minute from slot.delivery_start)::int) AS slot_delivery_start,
	(extract(hour from slot.delivery_until)::int * 60 + extract(minute from slot.delivery_until)::int) AS slot_delivery_until
	-- FROM --
	FROM orders AS o
	-- JOIN CLIENT --
	JOIN users AS u ON u.id = o.uid
	-- JOIN SHIPPING --
	JOIN users_addrs AS addr ON addr.id = o.shid
	-- JOIN SLOT --
	JOIN deliveries_slots AS slot ON slot.id = o.sloid
	`

	// # BEGIN FILTER #
	var values []any
	var clauses []string

	// # OPTIONAL FILTER'S #
	if filter != nil {
		if filter.Query != nil {
			q := strings.TrimSpace(*filter.Query)
			if q != "" {
				n, err := strconv.ParseInt(q, 10, 64)
				if err == nil {
					if n >= 1000 && n <= 99999 {
						values = append(values, n)
						clauses = append(clauses, fmt.Sprintf(`o.number = $%d`, len(values)))
					} else {
						values = append(values, q)
						clauses = append(clauses, fmt.Sprintf(`u.phone ILIKE '%%' || $%d || '%%'`, len(values)))
					}
				} else {
					values = append(values, q)
					clauses = append(clauses, fmt.Sprintf(`unaccent(u.name) ILIKE unaccent('%%' || $%d || '%%')`, len(values)))
				}
			}
		}
		if filter.Status != nil {
			q := strings.TrimSpace(*filter.Status)
			values = append(values, q)
			clauses = append(clauses, fmt.Sprintf(`o.status = $%d`, len(values)))
		}
		if filter.Delivery != nil {
			q := strings.TrimSpace(*filter.Delivery)
			d, err := time.Parse("2006-01-02", q)
			if err == nil {
				values = append(values, d.Format("2006-01-02"))
				clauses = append(clauses, fmt.Sprintf(`slot.wday = $%d`, len(values)))
			}
		}
		if filter.PaymentStatus != nil {
			q := strings.TrimSpace(*filter.PaymentStatus)
			values = append(values, q)
			clauses = append(clauses, fmt.Sprintf(`o.status <> 'canceled' AND o.payment_status = $%d`, len(values)))
		}
	}

	// # CLASUSES SEP #
	if len(clauses) > 0 {
		qry += " WHERE " + strings.Join(clauses, " AND ")
	}

	// # BEGIN SORT_BY #
	qry += ` ORDER BY o.date_created DESC`

	// # PAGINATION #
	if paging != nil {
		qry += fmt.Sprintf(` LIMIT %d `, paging.Limit)
		qry += fmt.Sprintf(` OFFSET %d `, paging.Offset)
	}

	// # END DEFAULT FILTER #

	rows, err := r.db.Query(ctx, qry, values...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	crows, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[OrderRaw])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results := make([]*Order, 0, len(crows))
	for i := range crows {
		results = append(results, crows[i].ToOrder())
	}

	return results, nil
}

// ORDER_LINE__

func (r Repository) OrderLineInsert(ctx context.Context, oid string, data *OrderLineInsertData) (string, error) {
	const op = "App.Repository.OrderLineInsert"
	qry := `INSERT INTO orders_lines (oid, pid, quantity, base_price) VALUES (@oid, @pid, @quantity, @base_price) RETURNING id;`

	var ref string
	if err := r.db.QueryRow(
		ctx,
		qry,
		pgx.NamedArgs{
			"oid":        oid,
			"pid":        data.Pid,
			"quantity":   data.Quantity,
			"base_price": data.BasePrice,
		},
	).Scan(&ref); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ref, nil

}

func (r Repository) OrderLineUpdate(ctx context.Context, ref string, paths []string, data *OrderLineUpdateData) (int64, error) {
	const op = "App.Repository.OrderLineUpdate"

	var values []any
	var clauses []string

	for _, path := range paths {
		switch path {
		case "status":
			values = append(values, data.Status)
			clauses = append(clauses, fmt.Sprintf(`status = $%d`, len(values)))
		case "quantity":
			values = append(values, data.Quantity)
			clauses = append(clauses, fmt.Sprintf(`quantity = $%d`, len(values)))
		case "base_price":
			values = append(values, data.BasePrice)
			clauses = append(clauses, fmt.Sprintf(`base_price = $%d`, len(values)))
		}
	}
	if len(clauses) == 0 {
		return 0, nil
	}

	// Update sets
	xquery := "UPDATE orders_lines"
	xquery += " SET " + strings.Join(clauses, ", ")

	// Where
	values = append(values, ref)
	xquery += fmt.Sprintf(" WHERE id = $%d", len(values))

	// exec
	res, err := r.db.Exec(ctx, xquery, values...)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.RowsAffected(), nil
}

func (r Repository) OrderLineDelete(ctx context.Context, ref string) (int64, error) {
	const op = "App.Repository.OrderLineDelete"

	qry := `DELETE FROM orders_lines WHERE id = $1 AND oid IN (SELECT id FROM orders WHERE status NOT IN ('canceled', 'successfully'));`

	tag, err := r.db.Exec(ctx, qry, ref)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return tag.RowsAffected(), nil
}

func (r Repository) OrderLineSelect(ctx context.Context, ref string, forUpdate bool) (*OrderLine, error) {
	const op = "App.Repository.OrderLineSelect"

	qry := `
	SELECT
	ol.id, ol.status, ol.quantity, ol.base_price, ol.total_price,
	p.id AS product_ref, p.upc AS product_upc, p.code AS product_code, p.name AS product_name, p.imurl AS product_imurl,
	p.display AS product_display, p.weight AS product_weight, p.unitype AS product_unitype, p.quantity AS product_quantity,
	p.is_active AS product_is_active, p.is_public AS product_is_public, p.cost_price AS product_cost_price, p.base_price AS product_base_price,
	p.num_in_alloc AS product_num_in_alloc, p.num_in_stock AS product_num_in_stock, p.date_created AS product_date_created
	-- FROM --
	FROM orders_lines AS ol
	-- JOIN PRODUCT --
	INNER JOIN products AS p ON p.id = ol.pid
	-- WHERE --
	WHERE ol.id = $1
	`

	// # FOR UPDATE #
	if forUpdate {
		qry += " FOR UPDATE"
	}

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[OrderLineRaw])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapOrderLineNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result.ToOrderLine(), nil
}

func (r Repository) OrderLineSelectAll(ctx context.Context, oid string) ([]*OrderLine, error) {
	const op = "App.Repository.OrderLineSelectAll"

	qry := `
	SELECT
	l.id, l.status, l.quantity, l.base_price, l.total_price,
	p.id AS product_ref, p.upc AS product_upc, p.code AS product_code, p.name AS product_name,
	p.imurl AS product_imurl, p.display AS product_display, p.weight AS product_weight, p.unitype AS product_unitype,
	p.quantity AS product_quantity, p.is_active AS product_is_active, p.is_public AS product_is_public, p.cost_price AS product_cost_price,
	p.base_price AS product_base_price, p.num_in_alloc AS product_num_in_alloc, p.num_in_stock AS product_num_in_stock,
	p.date_created AS product_date_created, p.date_updated AS product_date_updated
	-- FROM --
	FROM orders_lines AS l
	-- JOIN PRODUCT --
	INNER JOIN products AS p ON p.id = l.pid
	-- WHERE --
	WHERE l.oid = $1 ORDER BY l.date_created DESC;
	`

	rows, err := r.db.Query(ctx, qry, oid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[OrderLineRaw])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	orderLines := make([]*OrderLine, 0, len(results))
	for i := range results {
		orderLines = append(orderLines, results[i].ToOrderLine())
	}

	return orderLines, nil
}

// DELIVERY_DAY__

func (r Repository) DeliveryDaySelect(ctx context.Context, ref string, forUpdate bool) (*DeliverySlot, error) {
	const op = "App.Repository.DeliveryDaySelect"

	qry := `
	SELECT
	id, kind, note, wday, is_open, capacity, reserved,
	(extract(hour from cutoff_min)::int * 60 + extract(minute from cutoff_min)::int) AS cutoff_min, date_created, date_updated,
	(extract(hour from delivery_start)::int * 60 + extract(minute from delivery_start)::int) AS delivery_start,
	(extract(hour from delivery_until)::int * 60 + extract(minute from delivery_until)::int) AS delivery_until
	FROM deliveries_slots
	WHERE id = $1
	`

	if forUpdate {
		qry += ` FOR UPDATE`
	}

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[DeliverySlot])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapDeliverySlotNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (r Repository) DeliveryDaySelect2(ctx context.Context, wday string, forUpdate bool) (*DeliverySlot, error) {
	const op = "App.Repository.DeliveryDaySelect2"

	qry := `
	SELECT
	id, kind, note, wday, is_open, capacity, reserved,
	(extract(hour from cutoff_min)::int * 60 + extract(minute from cutoff_min)::int) AS cutoff_min, date_created, date_updated,
	(extract(hour from delivery_start)::int * 60 + extract(minute from delivery_start)::int) AS delivery_start,
	(extract(hour from delivery_until)::int * 60 + extract(minute from delivery_until)::int) AS delivery_until
	FROM deliveries_slots
	WHERE wday = $1
	`
	if forUpdate {
		qry += ` FOR UPDATE`
	}

	rows, err := r.db.Query(ctx, qry, wday)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[DeliverySlot])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapDeliverySlotNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (r Repository) DeliveryDayReserve(ctx context.Context, ref string) (int64, error) {
	const op = "App.Repository.DeliveryDayReserve"

	qry := `
	UPDATE deliveries_slots
	SET reserved = reserved + 1, date_updated = NOW()
	WHERE id = $1 AND is_open = TRUE AND reserved < capacity
	`

	tag, err := r.db.Exec(ctx, qry, ref)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return tag.RowsAffected(), nil
}

func (r Repository) DeliveryDayRelease(ctx context.Context, ref string) (int64, error) {
	const op = "App.Repository.DeliveryDayRelease"

	qry := `
	UPDATE deliveries_slots
	SET reserved = GREATEST(reserved - 1, 0), date_updated = NOW()
	WHERE id = $1
	`

	tag, err := r.db.Exec(ctx, qry, ref)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return tag.RowsAffected(), nil
}

func (r Repository) DeliveryDaySelectAll(ctx context.Context, filter *DeliveryDayFilterData, paging *DeliveryDayPagingData) ([]*DeliverySlot, error) {
	const op = "App.Repository.DeliveryDaySelectAll"

	qry := `
	SELECT
	id, kind, note, wday, is_open, capacity, reserved,
	(extract(hour from cutoff_min)::int * 60 + extract(minute from cutoff_min)::int) AS cutoff_min, date_created, date_updated,
	(extract(hour from delivery_start)::int * 60 + extract(minute from delivery_start)::int) AS delivery_start,
	(extract(hour from delivery_until)::int * 60 + extract(minute from delivery_until)::int) AS delivery_until
	FROM deliveries_slots
	`

	var values []any
	var clauses []string

	if filter != nil {
		if filter.Kind != nil {
			values = append(values, *filter.Kind)
			clauses = append(clauses, fmt.Sprintf(`kind = $%d`, len(values)))
		}
		if filter.IsOpen != nil {
			values = append(values, *filter.IsOpen)
			clauses = append(clauses, fmt.Sprintf(`is_open = $%d`, len(values)))
		}
		if filter.FromDate != nil {
			values = append(values, filter.FromDate.Format("2006-01-02"))
			clauses = append(clauses, fmt.Sprintf(`wday >= $%d`, len(values)))
		}
		if filter.UntilDate != nil {
			values = append(values, filter.UntilDate.Format("2006-01-02"))
			clauses = append(clauses, fmt.Sprintf(`wday <= $%d`, len(values)))
		}
	}

	if len(clauses) > 0 {
		qry += ` WHERE ` + strings.Join(clauses, ` AND `)
	}

	qry += ` ORDER BY wday ASC`

	if paging != nil {
		qry += fmt.Sprintf(` LIMIT %d OFFSET %d`, paging.Limit, paging.Offset)
	}

	rows, err := r.db.Query(ctx, qry, values...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[DeliverySlot])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return results, nil
}

//func (r Repository) DeliveryDaySelectAvailable(ctx context.Context, fromDate time.Time, applyCutoff bool, limit int32) ([]*DeliveryDay, error) {
//	const op = "App.Repository.DeliveryDaySelectAvailable"
//
//	qry := deliveryDaySelectColumns + `
//	WHERE is_open = TRUE
//	  AND reserved < capacity
//	  AND (
//	    work_date > $1
//	    OR (
//	      work_date = $1
//	      AND (
//	        NOT $2
//	        OR CURRENT_TIME <= make_time(cutoff_min / 60, cutoff_min % 60, 0)
//	      )
//	    )
//	  )
//	ORDER BY work_date ASC
//	LIMIT $3
//	`
//
//	rows, err := r.db.Query(ctx, qry, fromDate.Format("2006-01-02"), applyCutoff, limit)
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//	defer rows.Close()
//
//	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[DeliveryDay])
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	return results, nil
//}

//func (r Repository) DeliveryDaySelectNextAvailable(ctx context.Context, fromDate time.Time, applyCutoff bool) (*DeliveryDay, error) {
//	const op = "App.Repository.DeliveryDaySelectNextAvailable"
//
//	qry := deliveryDaySelectColumns + `
//	WHERE is_open = TRUE
//	  AND reserved < capacity
//	  AND (
//	    work_date > $1
//	    OR (
//	      work_date = $1
//	      AND (
//	        NOT $2
//	        OR CURRENT_TIME <= make_time(cutoff_min / 60, cutoff_min % 60, 0)
//	      )
//	    )
//	  )
//	ORDER BY work_date ASC
//	LIMIT 1
//	`
//
//	rows, err := r.db.Query(ctx, qry, fromDate.Format("2006-01-02"), applyCutoff)
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//	defer rows.Close()
//
//	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[DeliveryDay])
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			err = WrapDeliverySlotNotFound(err)
//		}
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	return result, nil
//}

//func (r Repository) DeliveryDayUpdate(ctx context.Context, workDate time.Time, paths []string, data *DeliveryDayUpdateData) (int64, error) {
//	const op = "App.Repository.DeliveryDayUpdate"
//
//	var values []any
//	var clauses []string
//
//	for _, path := range paths {
//		switch path {
//		case "kind":
//			values = append(values, data.Kind)
//			clauses = append(clauses, fmt.Sprintf(`kind = $%d`, len(values)))
//		case "note":
//			values = append(values, data.Note)
//			clauses = append(clauses, fmt.Sprintf(`note = $%d`, len(values)))
//		case "is_open":
//			values = append(values, data.IsOpen)
//			clauses = append(clauses, fmt.Sprintf(`is_open = $%d`, len(values)))
//		case "capacity":
//			values = append(values, data.Capacity)
//			clauses = append(clauses, fmt.Sprintf(`capacity = $%d`, len(values)))
//		case "cutoff_min":
//			values = append(values, data.CutoffMin)
//			clauses = append(clauses, fmt.Sprintf(`cutoff_min = $%d`, len(values)))
//		case "delivery_start":
//			values = append(values, data.DeliveryStart)
//			clauses = append(clauses, fmt.Sprintf(`delivery_start = $%d`, len(values)))
//		case "delivery_until":
//			values = append(values, data.DeliveryUntil)
//			clauses = append(clauses, fmt.Sprintf(`delivery_until = $%d`, len(values)))
//		}
//	}
//
//	if len(clauses) == 0 {
//		return 0, nil
//	}
//
//	clauses = append(clauses, "date_updated = NOW()")
//
//	xquery := `UPDATE deliveries_days`
//	xquery += ` SET ` + strings.Join(clauses, ", ")
//
//	values = append(values, workDate.Format("2006-01-02"))
//	xquery += fmt.Sprintf(` WHERE work_date = $%d`, len(values))
//
//	res, err := r.db.Exec(ctx, xquery, values...)
//	if err != nil {
//		return 0, fmt.Errorf("%s: %w", op, err)
//	}
//
//	return res.RowsAffected(), nil
//}
