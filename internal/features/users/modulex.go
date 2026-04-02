package users

//type UserRaw struct {
//	Ref        string     `db:"id"`
//	Idk        *string    `db:"idk"`
//	Rank       *float32   `db:"rank"`
//	Name       string     `db:"name"`
//	Phone      string     `db:"lookups"`
//	IsStaff    bool       `db:"is_staff"`
//	IsSuper    bool       `db:"is_super"`
//	IsActive   bool       `db:"is_active"`
//	LastLogin  *time.Time `db:"last_login"`
//	DateJoined time.Time  `db:"date_joined"`
//}
//
//func (raw *UserRaw) ToModel() *User {
//	return &User{
//		Ref:        raw.Ref,
//		Idk:        raw.Idk,
//		Name:       raw.Name,
//		Phone:      raw.Phone,
//		IsStaff:    raw.IsStaff,
//		IsSuper:    raw.IsSuper,
//		IsActive:   raw.IsActive,
//		LastLogin:  raw.LastLogin,
//		DateJoined: raw.DateJoined,
//	}
//}
