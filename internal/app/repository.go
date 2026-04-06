package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"apigo/internal/modules/postgres"
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

// SALES__

// ORDER__

// # ORDER #

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

func (r Repository) OrderSelect(ctx context.Context, ref string) (*Order, error) {
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
	slot.id AS slot_ref, slot.code AS slot_code, days.work_date AS slot_work_date
	-- FROM --
	FROM orders AS o
	-- JOIN CLIENT --
	JOIN users AS u ON u.id = o.uid
	-- JOIN SHIPPING --
	JOIN users_addrs AS addr ON addr.id = o.shid
	-- JOIN DELIVERY SLOT --
	JOIN deliveries_slots AS slot ON slot.id = o.sloid
	-- JOIN DELIVERY DAYS --  
	JOIN deliveries_days AS days ON days.id = slot.wid 
	-- WHERE --
	WHERE o.id = $1
	`

	rows, err := r.db.Query(ctx, qry, ref)
	if err != nil {
		return nil, fmt.Errorf("Repository.OrderSelect: [db query] [%w]", err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapOrderNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (r Repository) OrderSelectAll(ctx context.Context, filter *OrderFilterData, paging *OrderPagingData) ([]*Order, error) {
	const op = "App.Repository.OrderDelete"

	qry := `
	SELECT
	o.id, o.number, o.status, o.date_created, o.date_updated, o.payment_status, o.payment_method,
	-- USER --
	u.id AS user_ref, u.name AS user_name, u.phone AS user_phone,
	-- ADDR --
	addr.id AS addr_ref, addr.lat AS addr_lat, addr.lng AS addr_lng, addr.name AS addr_name, addr.cmna AS addr_cmna, addr.route AS addr_route,
	addr.street AS addr_street, addr.neighb AS addr_neighb, addr.locality AS addr_locality, addr.sublocal AS addr_sublocal, addr.address1 AS addr_address1, addr.address2 AS addr_address2,
	-- SLOT --
	slot.id AS slot_ref, slot.code AS slot_code, days.work_date AS slot_work_date
	-- FROM --
	FROM orders AS o
	-- JOIN CLIENT --
	JOIN users AS u ON u.id = o.uid
	-- JOIN SHIPPING --
	JOIN users_addrs AS addr ON addr.id = o.shid
	-- JOIN DELIVERY SLOT --
	JOIN deliveries_slots AS slot ON slot.id = o.sloid
	-- JOIN DELIVERY DAYS --
	JOIN deliveries_days AS days ON days.id = slot.wid
	-- ORDER BY --
	-- ORDER BY p.date_created DESC
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
				clauses = append(clauses, fmt.Sprintf(`days.work_date = $%d`, len(values)))
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

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[Order])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return results, nil
}
