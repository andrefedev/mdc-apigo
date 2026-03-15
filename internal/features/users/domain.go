package users

import "time"

type User struct {
	Ref        string     `json:"id"`
	Idk        *string    `json:"-"`
	Name       string     `json:"name"`
	Phone      string     `json:"lookups"`
	IsSuper    bool       `json:"is_super"`
	IsStaff    bool       `json:"is_staff"`
	IsActive   bool       `json:"is_active"`
	IsPremium  bool       `json:"is_premium"`
	LastLogin  *time.Time `json:"last_login"`
	DateJoined time.Time  `json:"date_joined"`
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
