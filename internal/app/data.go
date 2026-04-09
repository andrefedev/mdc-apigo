package app

import (
	"apigo/internal/platforms/validatex/normalizex"
	"fmt"
	"strings"
	"time"

	"apigo/internal/platforms/validatex/validationx"

	"github.com/google/uuid"
)

// # CODE__ #

type CodeInsertData struct {
	Code  string `db:"code"`
	Phone string `db:"phone"`
}

func (d *CodeInsertData) Validate() error {
	const op = "App.CodeInsertData.Validate"

	// Normalize
	d.Phone = validationx.ClearString(d.Phone)

	// Validation
	if !validationx.IsPhoneNumber(d.Phone) {
		return fmt.Errorf("%s: %w", op, ErrInvalidPhone)
	}

	return nil
}

// # SESSION__ #

type SessionInsertData struct {
	UserRef     string    `db:"uid"`
	TokenHash   string    `db:"token_hash"`
	DateExpired time.Time `db:"date_expired"`
}

// # USER__ #

type UserFilterData struct {
	IsSuper   *bool   `db:"is_super"`
	IsStaff   *bool   `db:"is_staff"`
	IsActive  *bool   `db:"is_active"`
	FlatQuery *string `db:"flat_query"`
}

func NewUserFilterData(input *UserFilterInput) *UserFilterData {
	if input == nil {
		return &UserFilterData{}
	}

	return &UserFilterData{
		IsSuper:   input.IsSuper,
		IsStaff:   input.IsStaff,
		IsActive:  input.IsActive,
		FlatQuery: input.FlatQuery,
	}
}

func (r *UserFilterData) Validate() error {
	const op = "App.UserFilterData.Validate"

	// Normalize
	if r.FlatQuery != nil {
		r.FlatQuery = new(normalizex.NormalizeName(*r.FlatQuery))
	}

	// Validation

	return nil
}

type UserPagingData struct {
	Limit  int32
	Offset int32
}

func NewUserPagingData(input *UserPagingInput) *UserPagingData {
	if input == nil {
		return &UserPagingData{}
	}
	return &UserPagingData{
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}

func (r *UserPagingData) Validate() error {
	const op = "App.UserPagingData.Validate"

	// Normalize
	limit := int32(20)
	if r.Limit == 0 {
		r.Limit = limit
	}
	if r.Limit > limit {
		r.Limit = limit
	}

	// Validation

	return nil
}

// USER_INSERT_DATA__

type UserInsertData struct {
	Name      string     `db:"name"`
	Phone     string     `db:"phone"`
	IsSuper   bool       `db:"is_super"`
	IsStaff   bool       `db:"is_staff"`
	IsActive  bool       `db:"is_active"`
	LastLogin *time.Time `db:"last_login"`
}

func NewUserInsertData(input *UserInsertInput) *UserInsertData {
	if input == nil {
		return &UserInsertData{}
	}
	return &UserInsertData{
		Name:     input.Name,
		Phone:    input.Phone,
		IsSuper:  input.IsSuper,
		IsStaff:  input.IsStaff,
		IsActive: input.IsActive,
	}

}

func (r *UserInsertData) Validate() error {
	const op = "App.UserInsertData.Validate"

	// Normalize
	r.Name = normalizex.NormalizeName(r.Name)
	r.Phone = validationx.ClearString(r.Phone)

	// Validation

	if r.Name == "" {
		return fmt.Errorf("%s: %w", op, ErrInvalidName)
	}

	if !validationx.IsPhoneNumber(r.Phone) {
		return fmt.Errorf("%s: %w", op, ErrInvalidPhone)
	}

	return nil
}

// USER_UPDATE_DATA__

type UserUpdateData struct {
	Name      string     `db:"name"`
	Phone     string     `db:"phone"`
	IsSuper   bool       `db:"is_super"`
	IsStaff   bool       `db:"is_staff"`
	IsActive  bool       `db:"is_active"`
	LastLogin *time.Time `db:"last_login"`
}

func NewUserUpdateData(input *UserUpdateInput) *UserUpdateData {
	if input == nil {
		return &UserUpdateData{}
	}
	return &UserUpdateData{
		Name:     input.Name,
		Phone:    input.Phone,
		IsSuper:  input.IsSuper,
		IsStaff:  input.IsStaff,
		IsActive: input.IsActive,
	}

}

func (r UserUpdateData) Validate(paths []string) error {
	const op = "App.UserUpdateData.Validate"

	for _, path := range paths {
		switch strings.TrimSpace(path) {
		case "name":
			r.Name = normalizex.NormalizeName(r.Name)
			if r.Name == "" {
				return fmt.Errorf("%s: %w", op, ErrInvalidName)
			}
		case "phone":
			r.Phone = validationx.ClearString(r.Phone)
			if !validationx.IsPhoneNumber(r.Phone) {
				return fmt.Errorf("%s: %w", op, ErrInvalidPhone)
			}
		case "last_login":
			if r.LastLogin == nil {
				return fmt.Errorf("%s: %w", op, ErrInvalidLastLogin)
			}
		}
	}

	return nil
}

// USER_ADDR__

// USER_ADDR_INSERT_DATA__

type UserAddrInsertData struct {
	Pid       string
	Lat       float64
	Lng       float64
	Name      string
	Cmna      string
	Route     string
	Street    string
	Neighb    string
	Locality  string
	Sublocal  string
	Address1  string // casa / apto complemento
	Address2  string // instrucciones de entrega
	IsDefault bool
}

func NewUserAddrInsertData(input *UserAddrCreateInput) *UserAddrInsertData {
	if input == nil {
		return &UserAddrInsertData{}
	}

	return &UserAddrInsertData{
		Pid:      input.Pid,
		Lat:      input.Lat,
		Lng:      input.Lng,
		Name:     input.Name,
		Cmna:     input.Cmna,
		Route:    input.Route,
		Street:   input.Street,
		Neighb:   input.Neighb,
		Locality: input.Locality,
	}
}

func (r *UserAddrInsertData) Validate() error {
	return nil
}

// USER_ADDR_UPDATE_DATA__

type UserAddrUpdateData struct {
	Pid       string
	Lat       float64
	Lng       float64
	Name      string
	Cmna      string
	Route     string
	Street    string
	Neighb    string
	Locality  string
	Sublocal  string
	Address1  string // casa / apto complemento
	Address2  string // instrucciones de entrega
	IsDefault bool
}

func NewUserAddrUpdateData(input *UserAddrUpdateInput) *UserAddrUpdateData {
	if input == nil {
		return &UserAddrUpdateData{}
	}

	return &UserAddrUpdateData{
		Lat:      input.Lat,
		Lng:      input.Lng,
		Route:    input.Route,
		Street:   input.Street,
		Address1: input.Address1,
		Address2: input.Address2,
	}
}

func (r *UserAddrUpdateData) Validate(paths []string) error {
	return nil
}

// SALES__

// ORDER_INSERT_INPUT__

type OrderInsertData struct {
	User          string
	Addr          string
	Slot          string
	Status        string
	PaymentStatus string
	PaymentMethod string
}

func NewOrderInsertData(input *OrderInsertInput) *OrderInsertData {
	if input == nil {
		return &OrderInsertData{}
	}
	return &OrderInsertData{
		User:          input.User,
		Addr:          input.Addr,
		Slot:          input.Slot,
		Status:        input.Status,
		PaymentStatus: input.PaymentStatus,
		PaymentMethod: input.PaymentMethod,
	}
}

func (r *OrderInsertData) Validate() error {
	//for _, path := range paths {
	//	switch strings.TrimSpace(path) {
	//	case "user":
	//		if r.User == "" {
	//			return errors.New("la referencia del usuario es un campo obligatorio")
	//		}
	//	case "addr":
	//		if r.Addr == "" {
	//			return errors.New("la referencia de la dirección de envío es un campo obligatorio")
	//		}
	//	case "slot":
	//		if r.Slot == "" {
	//			return errors.New("la referencia del día y franja horaria es un campo obligatorio")
	//		}
	//	case "status":
	//		// validar opciones del status
	//		if r.Status == "" {
	//			return errors.New("el estado del pedido es un campo obligatorio")
	//		}
	//	case "payment_status":
	//		if r.PaymentStatus == "" {
	//			return errors.New("el estado del pago del pedido es un obligatorio")
	//		}
	//	case "payment_method":
	//		if r.PaymentMethod == "" {
	//			return errors.New("el método del pago del pedido es un obligatorio")
	//		}
	//	}
	//}

	return nil
}

// ORDER_UPDATE_INPUT__

type OrderUpdateData struct {
	Addr          string
	Slot          string
	Status        string
	PaymentStatus string
	PaymentMethod string
}

func NewOrderUpdateData(input *OrderUpdateInput) *OrderUpdateData {
	if input == nil {
		return &OrderUpdateData{}
	}
	return &OrderUpdateData{
		Addr:          input.Addr,
		Slot:          input.Slot,
		Status:        input.Status,
		PaymentStatus: input.PaymentStatus,
		PaymentMethod: input.PaymentMethod,
	}
}

func (r *OrderUpdateData) Validate(paths []string) error {
	return nil
}

// ORDER_CHANGE_STATUS_DATA__

type OrderChangeStatusData struct {
	Status string
}

func NewOrderChangeStatusData(input *OrderChangeStatusInput) *OrderChangeStatusData {
	if input == nil {
		return &OrderChangeStatusData{}
	}
	return &OrderChangeStatusData{
		Status: input.Status,
	}
}

func (r *OrderChangeStatusData) Validate() error {
	const op = "App.OrderChangeStatusData.Validate"

	status := strings.TrimSpace(r.Status)
	switch status {
	case "pending", "acepted", "canceled", "dispatched", "successfully":
		r.Status = status
		return nil
	default:
		return fmt.Errorf("%s: %w", op, ErrInvalidOrderStatus)
	}
}

// ORDER_FILTER_DATA__

type OrderFilterData struct {
	Query         *string
	Status        *string
	Delivery      *string
	PaymentStatus *string
}

func NewOrderFilterData(input *OrderFilterInput) *OrderFilterData {
	if input == nil {
		return &OrderFilterData{}
	}
	return &OrderFilterData{
		Query:         input.Query,
		Status:        input.Status,
		Delivery:      input.Delivery,
		PaymentStatus: input.PaymentStatus,
	}
}

func (r *OrderFilterData) Validate() error {
	const op = "App.UserFilterData.Validate"

	if r.Query != nil {
		if *r.Query == "" {
			r.Query = nil
		} else {
			r.Query = new(normalizex.NormalizeName(*r.Query))
		}
	}

	if r.Query != nil && *r.Query == "" {
		if *r.Query == "" {
			return fmt.Errorf("%s: %w", op, ErrInvalidFlatQuery)
		}
	}

	return nil
}

// ORDER_PAGING_DATA__

type OrderPagingData struct {
	Limit  int32
	Offset int32
}

func NewOrderPagingData(input *OrderPagingInput) *OrderPagingData {
	if input == nil {
		return &OrderPagingData{}
	}
	return &OrderPagingData{
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}

func (r *OrderPagingData) Validate() error {
	const op = "App.OrderPagingData.Validate"

	// Normalize
	limit := int32(40)
	if r.Limit == 0 {
		r.Limit = limit
	}
	if r.Limit > limit {
		r.Limit = limit
	}

	// Validation

	return nil
}

// ORDER_LINE__

type OrderLineSelectData struct {
	Ref       string
	ForUpdate bool
}

type OrderLineInsertData struct {
	Pid        string
	Status     string
	Quantity   int32
	BasePrice  int32
	OfferPrice int32
}

func NewOrderLineInsertData(input *OrderLineCreateInput) *OrderLineInsertData {
	if input == nil {
		return &OrderLineInsertData{}
	}
	return &OrderLineInsertData{
		Pid:        input.Pid,
		Status:     input.Status,
		Quantity:   input.Quantity,
		BasePrice:  input.BasePrice,
		OfferPrice: input.OfferPrice,
	}
}

func (r *OrderLineInsertData) Validate() error {
	const op = "App.OrderLineInsertData.Validate"

	if uuid.Validate(r.Pid) != nil {
		return fmt.Errorf("%s: %w", op, ErrInvalidOrderLinePid)
	}

	if r.Quantity == 0 {
		return fmt.Errorf("%s: %w", op, ErrInvalidOrderLineQuantity)
	}

	if r.BasePrice == 0 {
		return fmt.Errorf("%s, %w", op, ErrInvalidOrderLineBasePrice)
	}

	if r.OfferPrice == 0 {
		return fmt.Errorf("%s: %w", op, ErrInvalidOrderLineOfferPrice)
	}

	// BASE_PRICE < OFFER_PRICE
	if r.BasePrice < r.OfferPrice {
		return fmt.Errorf("%s: %w", op, ErrInvalidOrderLinePriceRange) // nombrar..
	}

	return nil
}

type OrderLineUpdateData struct {
	Status     string
	Quantity   int32
	BasePrice  int32
	OfferPrice int32
}

func NewOrderLineUpdateData(input *OrderLineUpdateInput) *OrderLineUpdateData {
	if input == nil {
		return &OrderLineUpdateData{}
	}
	return &OrderLineUpdateData{
		Status:     input.Status,
		Quantity:   input.Quantity,
		BasePrice:  input.BasePrice,
		OfferPrice: input.OfferPrice,
	}
}

func (r *OrderLineUpdateData) Validate(paths []string) error {
	const op = "App.OrderLineUpdateData.Validate"

	priceRange := 0
	for _, path := range paths {
		switch strings.TrimSpace(path) {
		case "base_price":
			priceRange += 1
			if r.BasePrice == 0 {
				return fmt.Errorf("%s, %w", op, ErrInvalidOrderLineBasePrice)
			}
		case "offer_price":
			priceRange += 1
		}
	}

	if priceRange == 2 {
		if r.BasePrice < r.OfferPrice {
			return fmt.Errorf("%s: %w", op, ErrInvalidOrderLinePriceRange)
		}
	}

	return nil
}
