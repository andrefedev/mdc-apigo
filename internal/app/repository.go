package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[Code])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, WrapCodeNotFound(err))
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

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, WrapSessionNotFound(err))
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

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = WrapSessionNotFound(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return new(result), nil
}

// USER___

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
					clauses = append(clauses, fmt.Sprintf(`u.name ILIKE $%d`, len(values)))
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
