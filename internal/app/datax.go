package app

import (
	"apigo/internal/platforms/validatex/normalizex"
	"fmt"
	"strings"

	v1 "apigo/protobuf/gen/v1"

	"apigo/internal/platforms/validatex/validationx"

	"github.com/google/uuid"
)

// CODE__

type CodeInput struct {
	Phone string
}

func NewCodeInput(req *v1.CodeReq) *CodeInput {
	return &CodeInput{
		Phone: req.GetPhone(),
	}
}

func (r *CodeInput) Validate() error {
	const oper = "App.CodeInput.Validate"

	// Normalize
	r.Phone = validationx.ClearString(r.Phone)

	// Validation
	if !validationx.IsPhoneNumber(r.Phone) {
		return fmt.Errorf("%s: %w", oper, ErrInvalidPhone)
	}

	return nil
}

type CodeDetailInput struct {
	Ref string
}

func NewCodeDetailInput(req *v1.CodeDetailReq) *CodeDetailInput {
	return &CodeDetailInput{
		Ref: req.GetRef(),
	}
}

func (r *CodeDetailInput) Validate() error {
	const oper = "App.CodeDetailInput.Validate"

	// Normalize
	r.Ref = validationx.ClearString(r.Ref)

	// Validation
	if err := uuid.Validate(r.Ref); err != nil {
		return fmt.Errorf("%s: %w", oper, err)
	}

	return nil
}

// ····

type CodeVerifyInput struct {
	Ref  string
	Code string
}

func NewCodeVerifyInput(req *v1.CodeVerifyReq) *CodeVerifyInput {
	return &CodeVerifyInput{
		Ref:  req.GetRef(),
		Code: req.GetCode(),
	}
}

func (r *CodeVerifyInput) Validate() error {
	const oper = "App.VerifyCodeInput.Validate"
	// Normalize
	// Validation
	if !validationx.IsOneTimeCode(r.Code) {
		return fmt.Errorf("%s: %w", oper, ErrInvalidCode)
	}

	return nil
}

// # USER__ #

type UserPagingInput struct {
	Limit  int32
	Offset int32
}

func NewUserPagingInput(req *v1.UserListAllReq_Paging) *UserPagingInput {
	if req == nil {
		return &UserPagingInput{}
	}

	return &UserPagingInput{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
}

func (r *UserPagingInput) Validate() error {
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

// USER_FILTER_INPUT

type UserFilterInput struct {
	IsSuper   *bool
	IsStaff   *bool
	IsActive  *bool
	FlatQuery *string
}

func NewUserFilterInput(req *v1.UserListAllReq_Filter) *UserFilterInput {
	if req == nil {
		return &UserFilterInput{}
	}

	return &UserFilterInput{
		IsSuper:   req.IsSuper,
		IsStaff:   req.IsStaff,
		IsActive:  req.IsActive,
		FlatQuery: req.FlatQuery,
	}
}

func (r *UserFilterInput) Validate() error {
	const op = "App.UserFilterInput.Validate"

	// Normalize
	if r.FlatQuery != nil {
		r.FlatQuery = new(normalizex.NormalizeName(*r.FlatQuery))
	}

	// Validation
	if r.FlatQuery != nil {
		if len(*r.FlatQuery) < 2 {
			// return fmt.Errorf("%s: %w", oper, ErrInvalidPhone)
		}
	}

	return nil
}

// USER_INSERT_INPUT

type UserInsertInput struct {
	Name     string
	Phone    string
	IsSuper  bool
	IsStaff  bool
	IsActive bool
}

func NewUserInsertInput(payload *v1.UserCreateReq_Payload) *UserInsertInput {
	if payload == nil {
		return &UserInsertInput{}
	}

	return &UserInsertInput{
		Name:     payload.GetName(),
		Phone:    payload.GetPhone(),
		IsSuper:  payload.GetIsSuper(),
		IsStaff:  payload.GetIsStaff(),
		IsActive: payload.GetIsActive(),
	}
}

func (r *UserInsertInput) Validate() error {
	const op = "App.UserInsertInput.Validate"

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

// USER_UPDATE_INPUT

type UserUpdateInput struct {
	Name     string
	Phone    string
	IsSuper  bool
	IsStaff  bool
	IsActive bool
}

func NewUserUpdateInput(payload *v1.UserUpdateReq_Payload) *UserUpdateInput {
	if payload == nil {
		return &UserUpdateInput{}
	}

	return &UserUpdateInput{
		Name:     payload.GetName(),
		Phone:    payload.GetPhone(),
		IsSuper:  payload.GetIsSuper(),
		IsStaff:  payload.GetIsStaff(),
		IsActive: payload.GetIsActive(),
	}
}

func (r *UserUpdateInput) Validate(paths []string) error {
	const op = "App.UserUpdateInput.Validate"

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
		}
	}

	return nil
}

// USER_ADDR__

// USER_ADDR_INSERT_DATA__

type UserAddrCreateInput struct {
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

func NewUserAddrCreateInput(payload *v1.UserAddrCreateReq_Payload) *UserAddrCreateInput {
	if payload == nil {
		return &UserAddrCreateInput{}
	}
	return &UserAddrCreateInput{
		Pid:      payload.GetPid(),
		Lat:      payload.GetLat(),
		Lng:      payload.GetLng(),
		Name:     payload.GetName(),
		Cmna:     payload.GetCmna(),
		Route:    payload.GetRoute(),
		Street:   payload.GetStreet(),
		Neighb:   payload.GetNeighb(),
		Locality: payload.GetLocality(),
		Sublocal: payload.GetSublocal(),
	}
}

func (r *UserAddrCreateInput) Validate() error {
	return nil
}

// USER_ADDR_UPDATE_DATA__

type UserAddrUpdateInput struct {
	Lat      float64
	Lng      float64
	Route    string
	Street   string
	Address1 string // casa / apto complemento
	Address2 string // instrucciones de entrega
}

func NewUserAddrUpdateInput(payload *v1.UserAddrUpdateReq_Payload) *UserAddrUpdateInput {
	if payload == nil {
		return &UserAddrUpdateInput{}
	}
	return &UserAddrUpdateInput{
		Lat:      payload.GetLat(),
		Lng:      payload.GetLng(),
		Route:    payload.GetRoute(),
		Street:   payload.GetStreet(),
		Address1: payload.GetAddress1(),
		Address2: payload.GetAddress2(),
	}
}

func (r *UserAddrUpdateInput) Validate(paths []string) error {
	return nil
}

// SALES__

// ORDER_INSERT_INPUT__

type OrderInsertInput struct {
	User          string
	Addr          string
	Slot          string
	Status        string
	PaymentStatus string
	PaymentMethod string
}

func NewOrderInsertInput(payload *v1.OrderCreateReq_Payload) *OrderInsertInput {
	if payload == nil {
		return &OrderInsertInput{}
	}
	return &OrderInsertInput{
		User:          payload.GetUser(),
		Addr:          payload.GetAddr(),
		Slot:          payload.GetSlot(),
		Status:        payload.GetStatus(),
		PaymentStatus: payload.GetPaymentStatus(),
	}
}

func (r *OrderInsertInput) Validation(paths []string) error {
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

type OrderUpdateInput struct {
	Addr          string
	Slot          string
	Status        string
	PaymentStatus string
	PaymentMethod string
}

func NewOrderUpdateInput(payload *v1.OrderUpdateReq_Payload) *OrderUpdateInput {
	if payload == nil {
		return &OrderUpdateInput{}
	}
	return &OrderUpdateInput{
		Addr:          payload.GetAddr(),
		Slot:          payload.GetSlot(),
		Status:        payload.GetStatus(),
		PaymentStatus: payload.GetPaymentStatus(),
	}
}

func (r *OrderUpdateInput) Validation(paths []string) error {
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

// ORDER_FILTER_INPUT__

type OrderFilterInput struct {
	Query         *string
	Status        *string
	Delivery      *string
	PaymentStatus *string
}

func NewOrderFilterInput(req *v1.OrderListAllReq_Filter) *OrderFilterInput {
	if req == nil {
		return &OrderFilterInput{}
	}
	return &OrderFilterInput{
		Query:         req.Query,
		Status:        req.Status,
		Delivery:      req.Delivery,
		PaymentStatus: req.PaymentStatus,
	}
}

func (r *OrderFilterInput) Validate() error {
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

// ORDER_PAGING_INPUT__

type OrderPagingInput struct {
	Limit  int32
	Offset int32
}

func NewOrderPagingInput(req *v1.OrderListAllReq_Paging) *OrderPagingInput {
	if req == nil {
		return &OrderPagingInput{}
	}
	return &OrderPagingInput{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
}

func (r *OrderPagingInput) Validate() error {
	const op = "App.OrderPagingInput.Validate"

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

// ORDER_LINE_CREATE_INPUT__

type OrderLineCreateInput struct {
	Pid        string
	Status     string
	Quantity   int32
	BasePrice  int32
	OfferPrice int32
}

func NewOrderLineCreateInput(payload *v1.OrderLineCreateReq_Payload) *OrderLineCreateInput {
	if payload == nil {
		return &OrderLineCreateInput{}
	}

	return &OrderLineCreateInput{
		Pid:       payload.GetPid(),
		Status:    payload.GetStatus(),
		Quantity:  payload.GetQuantity(),
		BasePrice: payload.GetBasePrice(),
		// OfferPrice: payload.Ge
	}
}

func (r *OrderLineCreateInput) Validate() error {
	const op = "App.OrderLineCreateInput.Validate"

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

// ORDER_LINE_UPDATE_INPUT__

type OrderLineUpdateInput struct {
	Status     string
	Quantity   int32
	BasePrice  int32
	OfferPrice int32
}

func NewOrderLineUpdateInput(payload *v1.OrderLineUpdateReq_Payload) *OrderLineUpdateInput {
	if payload == nil {
		return &OrderLineUpdateInput{}
	}

	return &OrderLineUpdateInput{
		Status:    payload.GetStatus(),
		Quantity:  payload.GetQuantity(),
		BasePrice: payload.GetBasePrice(),
		// OfferPrice: payload.Ge
	}
}

func (r *OrderLineUpdateInput) Validate(paths []string) error {
	const op = "App.OrderLineUpdateInput.Validate"

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
