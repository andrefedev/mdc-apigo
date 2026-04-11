package app

import (
	"apigo/internal/platforms/validatex/normalizex"
	"fmt"
	"strings"
	"time"

	v1 "apigo/protobuf/gen/v1"

	"apigo/internal/platforms/validatex/validationx"

	"github.com/google/uuid"
	datepb "google.golang.org/genproto/googleapis/type/date"
)

func protoDateToTime(day *datepb.Date) (time.Time, error) {
	if day == nil {
		return time.Time{}, WrapInvalidDeliveryDayDate(nil)
	}

	result := time.Date(
		int(day.GetYear()),
		time.Month(day.GetMonth()),
		int(day.GetDay()),
		0, 0, 0, 0,
		time.UTC,
	)

	if result.Year() != int(day.GetYear()) ||
		int(result.Month()) != int(day.GetMonth()) ||
		result.Day() != int(day.GetDay()) {
		return time.Time{}, WrapInvalidDeliveryDayDate(nil)
	}

	return result, nil
}

func protoDateToTimePtr(day *datepb.Date) (*time.Time, error) {
	if day == nil {
		return nil, nil
	}

	result, err := protoDateToTime(day)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

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
	DeliveryDay   time.Time
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
		PaymentMethod: payload.GetPaymentMethod(),
	}
}

func (r *OrderInsertInput) Validate() error {
	const op = "App.OrderInsertInput.Validate"

	r.User = strings.TrimSpace(r.User)
	r.Addr = strings.TrimSpace(r.Addr)
	r.Slot = strings.TrimSpace(r.Slot)

	if err := uuid.Validate(r.User); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := uuid.Validate(r.Addr); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := uuid.Validate(r.Slot); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	status, err := normalizeOrderStatus(r.Status)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	paymentStatus, err := normalizeOrderPaymentStatus(r.PaymentStatus)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	paymentMethod, err := normalizeOrderPaymentMethod(r.PaymentMethod)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	r.Status = status
	r.PaymentStatus = paymentStatus
	r.PaymentMethod = paymentMethod
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
		PaymentMethod: payload.GetPaymentMethod(),
	}
}

func (r *OrderUpdateInput) Validation(paths []string) error {
	const op = "App.OrderUpdateInput.Validation"

	for _, path := range paths {
		switch strings.TrimSpace(path) {
		case "addr":
			r.Addr = strings.TrimSpace(r.Addr)
			if err := uuid.Validate(r.Addr); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case "slot":
			r.Slot = strings.TrimSpace(r.Slot)
			if err := uuid.Validate(r.Slot); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case "status":
			status, err := normalizeOrderStatus(r.Status)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			r.Status = status
		case "payment_status":
			status, err := normalizeOrderPaymentStatus(r.PaymentStatus)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			r.PaymentStatus = status
		case "payment_method":
			method, err := normalizeOrderPaymentMethod(r.PaymentMethod)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			r.PaymentMethod = method
		}
	}

	return nil
}

// ORDER_CHANGE_STATUS_INPUT__

type OrderChangeStatusInput struct {
	Status string
}

func NewOrderChangeStatusInput(req *v1.OrderChangeStatusReq) *OrderChangeStatusInput {
	if req == nil {
		return &OrderChangeStatusInput{}
	}
	return &OrderChangeStatusInput{
		Status: req.GetStatus(),
	}
}

func (r *OrderChangeStatusInput) Validate() error {
	const op = "App.OrderChangeStatusInput.Validate"

	status, err := normalizeOrderStatus(r.Status)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	r.Status = status
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

// DELIVERY_DAY__

type DeliveryDayDateInput struct {
	WorkDate time.Time
}

func NewDeliveryDayDateInput(day *datepb.Date) *DeliveryDayDateInput {
	if day == nil {
		return &DeliveryDayDateInput{}
	}

	workDate, _ := protoDateToTime(day)
	return &DeliveryDayDateInput{
		WorkDate: workDate,
	}
}

func (r *DeliveryDayDateInput) Validate() error {
	const op = "App.DeliveryDayDateInput.Validate"

	if r.WorkDate.IsZero() {
		return fmt.Errorf("%s: %w", op, ErrInvalidDeliveryDayDate)
	}

	return nil
}

type DeliveryDayFilterInput struct {
	FromDate  *time.Time
	UntilDate *time.Time
	IsOpen    *bool
	Kind      *string
}

func NewDeliveryDayFilterInput(req *v1.DeliverySlotListAllReq_Filter) *DeliveryDayFilterInput {
	if req == nil {
		return &DeliveryDayFilterInput{}
	}

	fromDate, _ := protoDateToTimePtr(req.GetFromDate())
	untilDate, _ := protoDateToTimePtr(req.GetUntilDate())

	return &DeliveryDayFilterInput{
		FromDate:  fromDate,
		UntilDate: untilDate,
		IsOpen:    req.IsOpen,
		Kind:      req.Kind,
	}
}

func (r *DeliveryDayFilterInput) Validate() error {
	const op = "App.DeliveryDayFilterInput.Validate"

	if r.Kind != nil {
		value := strings.TrimSpace(*r.Kind)
		if value == "" {
			r.Kind = nil
		} else {
			r.Kind = &value
		}
	}

	if r.FromDate != nil && r.UntilDate != nil && r.UntilDate.Before(*r.FromDate) {
		return fmt.Errorf("%s: %w", op, ErrInvalidDeliveryDayRange)
	}

	return nil
}

// DELIVERY_DAY_PAGING_INPUT__

type DeliveryDayPagingInput struct {
	Limit  int32
	Offset int32
}

func NewDeliveryDayPagingInput(req *v1.DeliverySlotListAllReq_Paging) *DeliveryDayPagingInput {
	if req == nil {
		return &DeliveryDayPagingInput{}
	}

	return &DeliveryDayPagingInput{
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	}
}

func (r *DeliveryDayPagingInput) Validate() error {
	const limit int32 = 90
	if r.Limit <= 0 {
		r.Limit = 30
	}
	if r.Limit > limit {
		r.Limit = limit
	}
	if r.Offset < 0 {
		r.Offset = 0
	}

	return nil
}

// CATLG__

type ProductFilterInput struct {
	Query    *string
	Genre    *string
	IsActive *bool
	IsPublic *bool
}

func NewProductFilterInput(req *v1.ProductListAllReq_Filter) *ProductFilterInput {
	if req == nil {
		return &ProductFilterInput{}
	}

	return &ProductFilterInput{
		Query:    req.Query,
		Genre:    req.Genre,
		IsActive: req.IsActive,
		IsPublic: req.IsPublic,
	}
}

func (r *ProductFilterInput) Validate() error {
	const op = "App.ProductFilterInput.Validate"

	if r.Query != nil {
		value := strings.TrimSpace(*r.Query)
		if value == "" {
			r.Query = nil
		} else {
			value = normalizex.NormalizeName(value)
			r.Query = &value
		}
	}

	if r.Genre != nil {
		value := strings.TrimSpace(*r.Genre)
		if value == "" {
			r.Genre = nil
		} else {
			if err := uuid.Validate(value); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			r.Genre = &value
		}
	}

	return nil
}
