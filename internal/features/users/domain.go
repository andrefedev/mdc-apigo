package users

import (
	v1 "apigo/protobuf/gen/v1"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	Ref        string     `db:"id"`
	Name       string     `db:"name"`
	Phone      string     `db:"phone"`
	IsStaff    bool       `db:"is_staff"`
	IsSuper    bool       `db:"is_super"`
	IsActive   bool       `db:"is_active"`
	LastLogin  *time.Time `db:"last_login"`
	DateJoined time.Time  `db:"date_joined"`
}

func (u User) ToProto() *v1.User {
	var dateJoined *timestamppb.Timestamp
	if !u.DateJoined.IsZero() {
		dateJoined = timestamppb.New(u.DateJoined)
	}

	var lastLogin *timestamppb.Timestamp
	if u.LastLogin != nil && !u.LastLogin.IsZero() {
		lastLogin = timestamppb.New(*u.LastLogin)
	}

	return &v1.User{
		Ref:        u.Ref,
		Name:       u.Name,
		Phone:      u.Phone,
		IsSuper:    u.IsSuper,
		IsStaff:    u.IsStaff,
		IsActive:   u.IsActive,
		LastLogin:  lastLogin,
		DateJoined: dateJoined,
	}
}

//
//import (
//	v1 "apigo/protobuf/gen/v1"
//	"time"
//
//	"google.golang.org/protobuf/types/known/timestamppb"
//)
//
//type User struct {
//	Ref        string
//	Idk        *string
//	Name       string
//	Phone      string
//	IsSuper    bool
//	IsStaff    bool
//	IsActive   bool
//	IsPremium  bool
//	LastLogin  *time.Time
//	DateJoined time.Time
//}
//
//func (u *User) ToProto() *v1.User {
//	var dateJoined *timestamppb.Timestamp
//	if !u.DateJoined.IsZero() {
//		dateJoined = timestamppb.New(u.DateJoined)
//	}
//
//	var lastLogin *timestamppb.Timestamp
//	if u.LastLogin != nil && !u.LastLogin.IsZero() {
//		lastLogin = timestamppb.New(*u.LastLogin)
//	}
//
//	return &v1.User{
//		Ref:        u.Ref,
//		Name:       u.Name,
//		Phone:      u.Phone,
//		IsSuper:    u.IsSuper,
//		IsStaff:    u.IsStaff,
//		IsActive:   u.IsActive,
//		LastLogin:  lastLogin,
//		DateJoined: dateJoined,
//	}
//}
//
//// USER_ADDR
//
//type Addr struct {
//	Ref         string
//	Pid         string
//	Lat         float64
//	Lng         float64
//	Name        string
//	Cmna        string
//	Route       string
//	Street      string
//	Neighb      string
//	Locality    string
//	Sublocal    string
//	Address1    string // casa / apto complemento
//	Address2    string // instrucciones de entrega
//	IsDefault   bool
//	DateCreated time.Time
//	DateUpdated *time.Time
//}
