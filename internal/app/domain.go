package app

import (
	"time"

	v1 "apigo/protobuf/gen/v1"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AUTH__

type Code struct {
	Ref         string    `db:"id"`
	Code        string    `db:"code"`
	Phone       string    `db:"phone"`
	DateCreated time.Time `db:"date_created"`
	DateExpired time.Time `db:"date_expired"`
}

func (r *Code) ToProto() *v1.Code {
	return &v1.Code{
		Ref:         r.Ref,
		Phone:       r.Phone,
		DateCreated: timestamppb.New(r.DateCreated),
		DateExpired: timestamppb.New(r.DateExpired),
	}
}

type Session struct {
	Ref         string     `db:"id"`
	UserRef     string     `db:"uid"`
	IsSuper     bool       `db:"is_super"`
	IsStaff     bool       `db:"is_staff"`
	IsActive    bool       `db:"is_active"`
	TokenHash   string     `db:"token_hash"`
	DateExpired time.Time  `db:"date_expired"`
	DateCreated time.Time  `db:"date_created"`
	DateRevoked *time.Time `db:"date_revoked"`
}

func (i *Session) IsRoot() bool {
	return i != nil && i.IsActive && i.IsSuper
}

func (i *Session) IsEmployee() bool {
	return i != nil && i.IsActive && (i.IsStaff || i.IsSuper)
}

// USER__

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

// USER_ADDR__

type UserAddr struct {
	Ref         string     `db:"id"`
	Pid         string     `db:"pid"`
	Lat         float64    `db:"lat"`
	Lng         float64    `db:"lng"`
	Name        string     `db:"name"`
	Cmna        string     `db:"cmna"`
	Route       string     `db:"route"`
	Street      string     `db:"street"`
	Neighb      string     `db:"neighb"`
	Locality    string     `db:"locality"`
	Sublocal    string     `db:"sublocal"`
	Address1    string     `db:"address1"` // casa / apto complemento
	Address2    string     `db:"address2"` // instrucciones de entrega
	IsDefault   bool       `db:"is_default"`
	DateCreated time.Time  `db:"date_created"`
	DateUpdated *time.Time `db:"date_updated"`
}

func (u *UserAddr) ToProto() *v1.UserAddr {
	var dateCreated *timestamp.Timestamp
	if !u.DateCreated.IsZero() {
		dateCreated = timestamppb.New(u.DateCreated)
	}

	var dateUpdated *timestamp.Timestamp
	if u.DateUpdated != nil && !u.DateUpdated.IsZero() {
		dateUpdated = timestamppb.New(*u.DateUpdated)
	}

	return &v1.UserAddr{
		Ref:         u.Ref,
		Pid:         u.Pid,
		Lat:         u.Lat,
		Lng:         u.Lng,
		Name:        u.Name,
		Cmna:        u.Cmna,
		Route:       u.Route,
		Street:      u.Street,
		Neighb:      u.Neighb,
		Locality:    u.Locality,
		Sublocal:    u.Sublocal,
		Address1:    u.Address1,
		Address2:    u.Address2,
		IsDefault:   u.IsDefault,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}
}

// CATLG__

// GENRE__

type Genre struct {
	Ref         string
	Name        string
	Descr       *string
	Imurl       *string
	Display     int32
	IsPublic    bool
	DateCreated time.Time
}

func (uc *Genre) ToProto() *v1.Genre {
	var dateCreated *timestamp.Timestamp

	if !uc.DateCreated.IsZero() {
		dateCreated = timestamppb.New(uc.DateCreated)
	}

	return &v1.Genre{
		Ref:         uc.Ref,
		Name:        uc.Name,
		Descr:       uc.Descr,
		Imurl:       uc.Imurl,
		Display:     uc.Display,
		IsPublic:    uc.IsPublic,
		DateCreated: dateCreated,
	}
}

// PRODUCT__

type Product struct {
	Ref     string
	Upc     *string
	Code    int32
	Name    string
	Genre   *Genre
	Descr   *string
	Imurl   *string
	Display int32

	Weight   int32
	Unitype  string
	Quantity int32

	IsActive    bool
	IsPublic    bool
	CostPrice   int32
	BasePrice   int32
	NumInStock  int32
	NumInAlloc  int32
	NumInAvail  int32
	DateCreated time.Time
	DateUpdated *time.Time
}

func (p *Product) ToProto() *v1.Product {
	var genrepb *v1.Genre
	if p.Genre != nil {
		// genre puede ser nil
		genrepb = p.Genre.ToProto()
	}

	var dateCreated *timestamp.Timestamp
	if !p.DateCreated.IsZero() {
		dateCreated = timestamppb.New(p.DateCreated)
	}

	var dateUpdated *timestamp.Timestamp
	if p.DateUpdated != nil && !p.DateUpdated.IsZero() {
		dateCreated = timestamppb.New(*p.DateUpdated)
	}

	return &v1.Product{
		Ref:         p.Ref,
		Upc:         p.Upc,
		Code:        p.Code,
		Name:        p.Name,
		Genre:       genrepb,
		Descr:       p.Descr,
		Imurl:       p.Imurl,
		Display:     p.Display,
		Weight:      p.Weight,
		Unitype:     p.Unitype,
		Quantity:    p.Quantity,
		IsPublic:    p.IsPublic,
		IsActive:    p.IsActive,
		CostPrice:   p.CostPrice,
		BasePrice:   p.BasePrice,
		NumInStock:  p.NumInStock,
		NumInAlloc:  p.NumInAlloc,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}
}

// SALES__

// ORDER__

type Order struct {
	Ref           string `db:"id"`
	User          *User
	Addr          *UserAddr
	Number        int32      `db:"number"`
	Status        string     `db:"status"`
	BasePrice     int32      `db:"base_price"`
	DiscPrice     int32      `db:"disc_price"`
	DateCreated   time.Time  `db:"date_created"`
	DateUpdated   *time.Time `db:"date_updated"`
	DeliveryDate  time.Time  `db:"delivery_date"`
	PaymentStatus string     `db:"payment_status"`
	PaymentMethod string     `db:"payment_method"`
}

func (p *Order) ToProto() *v1.Order {
	user := p.User.ToProto()
	addr := p.Addr.ToProto()

	var dateCreated *timestamp.Timestamp
	if !p.DateCreated.IsZero() {
		dateCreated = timestamppb.New(p.DateCreated)
	}

	var dateUpdated *timestamp.Timestamp
	if p.DateUpdated != nil && !p.DateUpdated.IsZero() {
		dateCreated = timestamppb.New(*p.DateUpdated)
	}

	deliveryDate := &date.Date{
		Day:   int32(p.DeliveryDate.Day()),
		Year:  int32(p.DeliveryDate.Year()),
		Month: int32(p.DeliveryDate.Month()),
	}

	return &v1.Order{
		Ref:           p.Ref,
		User:          user,
		Addr:          addr,
		Number:        p.Number,
		Status:        p.Status,
		DateCreated:   dateCreated,
		DateUpdated:   dateUpdated,
		DeliveryDate:  deliveryDate,
		PaymentStatus: p.PaymentStatus,
		PaymentMethod: p.PaymentMethod,
	}
}

// ORDER_LINE__

// ORDER_LINE__

type OrderLine struct {
	Ref        string
	Item       *Product
	Status     string
	Quantity   int32
	BasePrice  int32
	TotalPrice int32
}

func (p *OrderLine) ToProto() *v1.OrderLine {
	product := p.Item.ToProto()

	return &v1.OrderLine{
		Ref:        p.Ref,
		Item:       product,
		Status:     p.Status,
		Quantity:   p.Quantity,
		BasePrice:  p.BasePrice,
		TotalPrice: p.TotalPrice,
	}
}

// DELIVERY_DAY__

type DeliveryDay struct {
	Ref           string     `db:"id"`
	Kind          string     `db:"kind"`
	Note          *string    `db:"note"`
	Wday          time.Time  `db:"work_date"`
	IsOpen        bool       `db:"is_open"`
	Capacity      int32      `db:"capacity"`
	Reserved      int32      `db:"reserved"`
	CutoffMin     int32      `db:"cutoff_min"`
	DateCreated   time.Time  `db:"date_created"`
	DateUpdated   *time.Time `db:"date_updated"`
	DeliveryStart int32      `db:"delivery_start"`
	DeliveryUntil int32      `db:"delivery_until"`
}

func (p *DeliveryDay) ToProto() *v1.DeliveryDay {
	workDate := &date.Date{
		Day:   int32(p.Wday.Day()),
		Year:  int32(p.Wday.Year()),
		Month: int32(p.Wday.Month()),
	}

	var dateCreated *timestamp.Timestamp
	if !p.DateCreated.IsZero() {
		dateCreated = timestamppb.New(p.DateCreated)
	}

	var dateUpdated *timestamp.Timestamp
	if p.DateUpdated != nil && !p.DateUpdated.IsZero() {
		dateUpdated = timestamppb.New(*p.DateUpdated)
	}

	return &v1.DeliveryDay{
		Ref:           p.Ref,
		Kind:          p.Kind,
		Note:          p.Note,
		WorkDate:      workDate,
		IsOpen:        p.IsOpen,
		Capacity:      p.Capacity,
		Reserved:      p.Reserved,
		CutoffMin:     int64(p.CutoffMin),
		DateCreated:   dateCreated,
		DateUpdated:   dateUpdated,
		DeliveryStart: int64(p.DeliveryStart),
		DeliveryUntil: int64(p.DeliveryUntil),
	}
}
