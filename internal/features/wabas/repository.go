package wabas

import (
	"apigo/internal/modules/postgres"
)

type Repository struct {
	db *postgres.Pgdb
}

func NewRepository(db *postgres.Pgdb) *Repository {
	return &Repository{db: db}
}

//func (r Repository) MessageInsert(ctx context.Context, data *CodeInsertData) (string, error) {
//	op := "Auth.Repository.CodeInsert"
//	qry := `INSERT INTO users_codes (code, phone) VALUES(@code, @phone) returning id;`
//
//	// CodeData
//	var ref string
//	if err := r.db.QueryRow(
//		ctx,
//		qry,
//		pgx.NamedArgs{
//			"code":  data.Code,
//			"phone": data.Phone,
//		},
//	).Scan(&ref); err != nil {
//		return "", apperr.Internal(op, err)
//	}
//
//	return ref, nil
//}
