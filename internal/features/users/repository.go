package users

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"apigo/internal/modules/postgres"

	"github.com/jackc/pgx/v5"
)

// # USER #

type Repository struct {
	db *postgres.Pgdb
}

func NewRepository(db *postgres.Pgdb) *Repository {
	return &Repository{db: db}
}

//func (r Repository) Insert(ctx context.Context, data *Data) (string, error) {
//	query := `INSERT INTO users (name, lookups, is_super, is_staff, is_active) VALUES (@name, @lookups, @is_super, @is_staff, @is_active) RETURNING id;`
//
//	var ref string
//	if err := r.db.QueryRow(
//		ctx,
//		query,
//		pgx.NamedArgs{
//			"name":      data.Name,
//			"lookups":     data.Phone,
//			"is_super":  data.IsSuper,
//			"is_staff":  data.IsStaff,
//			"is_active": data.IsActive,
//		},
//	).Scan(&ref); err != nil {
//		return "", fmt.Errorf("User.Repository.Insert: [db query row]: [%w]", err)
//	}
//
//	return ref, nil
//}
//
//func (r Repository) Update(ctx context.Context, ref string, paths []string, data *Data) (int64, error) {
//	var values []any
//	var clauses []string
//	for _, path := range paths {
//		switch path {
//		case "idk":
//			values = append(values, data.Idk)
//			clauses = append(clauses, fmt.Sprintf(`idk = $%d`, len(values)))
//		case "name":
//			values = append(values, data.Name)
//			clauses = append(clauses, fmt.Sprintf(`name = $%d`, len(values)))
//		case "lookups":
//			values = append(values, data.Phone)
//			clauses = append(clauses, fmt.Sprintf(`lookups = $%d`, len(values)))
//		case "is_super":
//			values = append(values, data.IsSuper)
//			clauses = append(clauses, fmt.Sprintf(`is_super = $%d`, len(values)))
//		case "is_staff":
//			values = append(values, data.IsStaff)
//			clauses = append(clauses, fmt.Sprintf(`is_staff = $%d`, len(values)))
//		case "is_active":
//			values = append(values, data.IsActive)
//			clauses = append(clauses, fmt.Sprintf(`is_active = $%d`, len(values)))
//		case "last_login":
//			values = append(values, data.LastLogin)
//			clauses = append(clauses, fmt.Sprintf(`last_login = $%d`, len(values)))
//		}
//	}
//	if len(clauses) == 0 {
//		return 0, nil // nada que actualizar
//	}
//
//	// Update sets
//	qry := "UPDATE users"
//	qry += " SET " + strings.Join(clauses, ", ")
//
//	// Where
//	values = append(values, ref)
//	qry += fmt.Sprintf(" WHERE id = $%d", len(values))
//
//	// exec
//	res, err := r.db.Exec(ctx, qry, values...)
//	if err != nil {
//		return 0, fmt.Errorf("User.Repository.Update: [db query exec] [%w]", err)
//	}
//
//	return res.RowsAffected(), nil
//}

func (r Repository) Select(ctx context.Context, ref string) (*User, error) {
	const op = "User.Repository.Select"

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

func (r Repository) SelectByPhone(ctx context.Context, phone string) (*User, error) {
	const op = "User.Repository.SelectByPhone"

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

func (r Repository) SelectAll(ctx context.Context, filter *FilterData, paging *PagingData) ([]*User, error) {
	const op = "User.Repository.SelectAll"

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

func (r Repository) Insert(ctx context.Context, data *InserData) (string, error) {
	query := `INSERT INTO users (name, phone, is_super, is_staff, is_active) VALUES (@name, @phone, @is_super, @is_staff, @is_active) RETURNING id;`

	var ref string
	if err := x.queryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"name":      data.Name,
			"phone":     data.Phone,
			"is_super":  data.IsSuper,
			"is_staff":  data.IsStaff,
			"is_active": data.IsActive,
		},
	).Scan(&ref); err != nil {
		return "", fmt.Errorf("Repository.UserInsert: [db query row]: [%w]", err)
	}

	return ref, nil
}

func (r Repository) Update(ctx context.Context, ref string, paths []string, data *domain.UserData) (int64, error) {
	var values []any
	var clauses []string
	for _, path := range paths {
		switch path {
		case "idk":
			values = append(values, data.Idk)
			clauses = append(clauses, fmt.Sprintf(`idk = $%d`, len(values)))
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
	res, err := x.exec(ctx, xquery, values...)
	if err != nil {
		return 0, fmt.Errorf("Repository.UserUpdate: [db query exec] [%w]", err)
	}

	return res.RowsAffected(), nil
}

//func (r Repository) ExistsByPhone(ctx context.Context, lookups string) (bool, error) {
//	qry := `SELECT EXISTS (SELECT 1 FROM users WHERE lookups = $1);`
//
//	var exists bool
//	err := r.db.QueryRow(ctx, qry, lookups).Scan(&exists)
//	if err != nil {
//		return false, fmt.Errorf("User.Repository.ExistsByPhone: [db query row] %w", err)
//	}
//
//	return exists, nil
//}
