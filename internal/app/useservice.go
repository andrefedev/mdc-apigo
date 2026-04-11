package app

import (
	"apigo/internal/modules/gmaps"
	"apigo/internal/modules/whatsapp/messages"
	"apigo/internal/platforms/cryptox"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UseService struct {
	deps UseServiceDeps
}

type UseServiceDeps struct {
	Repository     *Repository
	GoogleMapx     *gmaps.Client
	MessageService *messages.Service
}

func NewUseService(deps UseServiceDeps) *UseService {
	return &UseService{deps: deps}
}

func slicesContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func (s *UseService) Code(ctx context.Context, input *CodeInput) (string, string, error) {
	oper := "App.UseService.Code"

	code, err := cryptox.GenerateRandomNumberString(6)
	if err != nil {
		return "", "", fmt.Errorf("%s: generate otp: %w", oper, err)
	}

	data := &CodeInsertData{
		Code:  code,
		Phone: input.Phone,
	}
	if err := data.Validate(); err != nil {
		return "", "", fmt.Errorf("%s: %w", oper, err)
	}

	ref, err := s.deps.Repository.CodeInsert(ctx, data)
	if err != nil {
		return "", "", fmt.Errorf("%s: insert code: %w", oper, err)
	}

	templ := &messages.TemplateMessageRequest{
		To:   data.Phone,
		Type: messages.TypeTemplate,
		Template: &messages.TemplContent{
			Name: "verify_code",
			Language: messages.TemplLang{
				Code: "es_CO",
			},
			Components: []messages.TemplComp{
				{
					Type: "body",
					Parameters: []messages.TemplParam{
						{
							Type: "text",
							Text: new(code),
						},
					},
				},
				{
					Type:    "button",
					SubType: new("url"),
					Index:   new(0),
					Parameters: []messages.TemplParam{
						{
							Type: "text",
							Text: new(code),
						},
					},
				},
			},
		},
	}
	if err := s.deps.MessageService.SendTemplate(ctx, templ); err != nil {
		return "", "", fmt.Errorf("%s: send template: %w", oper, err)
	}

	return ref, code, nil
}

func (s *UseService) CodeVerify(ctx context.Context, input *CodeVerifyInput) (string, string, error) {
	const op = "App.UseService.CodeVerify"

	var uid string
	var idk string
	// REQUIERE SESSION >>>
	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		code, err := s.deps.Repository.CodeSelect(ctx, input.Ref)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if code.Code != input.Code {
			return fmt.Errorf("%s: %w", op, WrapInvalidCode(nil))
		}
		if time.Now().After(code.DateExpired) {
			return fmt.Errorf("%s: %w", op, WrapCodeExpired(nil))
		}

		// Es probable que el usuario no exista
		// por lo que se debe crearlo.
		uid, err = s.deps.Repository.UserRefByPhone(ctx, code.Phone)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		idk, err = cryptox.GenerateRandomString(32)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		// session insert
		if _, err := s.deps.Repository.SessionInsert(
			ctx,
			&SessionInsertData{
				UserRef:   uid,
				TokenHash: cryptox.HashIdToken(idk),
			},
		); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		// SUCCESS...
		if _, err := s.deps.Repository.CodeDelete(ctx, input.Ref); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}); err != nil {
		return "", "", err
	}

	return uid, idk, nil
}

func (s *UseService) CodeDetail(ctx context.Context, input *CodeDetailInput) (*Code, error) {
	const op = "App.UseService.CodeDetail"

	res, err := s.deps.Repository.CodeSelect(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if time.Now().After(res.DateExpired) {
		return nil, WrapCodeExpired(err)
	}

	return res, nil
}

func (s *UseService) SessionByIdToken(ctx context.Context, idk string) (*Session, error) {
	const op = "App.UseService.SessionByIdToken"

	if idk == "" {
		return nil, fmt.Errorf("%s: %w", op, WrapSessionRequired(nil))
	}

	idk = cryptox.HashIdToken(idk)
	session, err := s.deps.Repository.SessionSelectByToken(ctx, idk)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			err = WrapSessionRequired(err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// session expirada
	if time.Now().After(session.DateExpired) {
		return nil, fmt.Errorf("%s: %w", op, WrapSessionExpired(nil))
	}

	// Session revocada
	if session.DateRevoked != nil {
		return nil, fmt.Errorf("%s: %w", op, WrapSessionRevoked(nil))
	}

	return session, nil
}

// USER__

func (s *UseService) UserCreate(ctx context.Context, input *UserInsertInput) (*User, error) {
	const op = "App.UseService.UserCreate"

	// manejar: usuario existente
	user, err := s.deps.Repository.UserSelectByPhone(ctx, input.Phone)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if user != nil && user.Phone == input.Phone {
		return nil, fmt.Errorf("%s: %w", op, WrapUserExists(nil))
	}

	data := NewUserInsertData(input)
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		ref, err := s.deps.Repository.UserInsert(ctx, data)
		if err != nil {
			return err
		}

		user, err = s.deps.Repository.UserSelect(ctx, ref)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UseService) UserUpdate(ctx context.Context, ref string, paths []string, input *UserUpdateInput) (*User, error) {
	const op = "App.UseService.UserUpdate"

	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data := NewUserUpdateData(input)
	if err := data.Validate(paths); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// manejar: usuario
	user, err := s.deps.Repository.UserSelect(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		_, err := s.deps.Repository.UserUpdate(ctx, ref, paths, data)
		if err != nil {
			return err
		}

		user, err = s.deps.Repository.UserSelect(ctx, ref)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UseService) UserDetail(ctx context.Context, ref string) (*User, error) {
	const op = "App.UseService.UserDetail"

	user, err := s.deps.Repository.UserSelect(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UseService) UserListAll(ctx context.Context, filter *UserFilterInput, paging *UserPagingInput) ([]*User, error) {
	const op = "App.UseService.UserListAll"

	// aqui se convierten
	f := NewUserFilterData(filter)
	if err := f.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	p := NewUserPagingData(paging)
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	users, err := s.deps.Repository.UserSelectAll(ctx, f, p)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

// USER_ADDR__

func (s *UseService) UserAddrCreate(ctx context.Context, uid string, input *UserAddrCreateInput) (*UserAddr, error) {
	const op = "App.UseService.UserAddrCreate"

	// manejar: usuario existente
	_, err := s.deps.Repository.UserSelect(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data := NewUserAddrInsertData(input)
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var result *UserAddr
	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		ref, err := s.deps.Repository.UserAddrInsert(ctx, uid, data)
		if err != nil {
			return err
		}

		result, err = s.deps.Repository.UserAddrSelect(ctx, ref)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) UserAddrUpdate(ctx context.Context, ref string, paths []string, input *UserAddrUpdateInput) (*UserAddr, error) {
	const op = "App.UseService.UserAddrUpdate"

	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	updata := NewUserAddrUpdateData(input)
	if err := updata.Validate(paths); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userAddr, err := s.deps.Repository.UserAddrSelect(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		_, err := s.deps.Repository.UserAddrUpdate(ctx, ref, paths, updata)
		if err != nil {
			return err
		}

		userAddr, err = s.deps.Repository.UserAddrSelect(ctx, ref)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userAddr, nil
}

func (s *UseService) UserAddrDetail(ctx context.Context, ref string) (*UserAddr, error) {
	const op = "App.UseService.UserAddrDetail"

	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.UserAddrSelect(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) UserAddrListAll(ctx context.Context, uid string) ([]*UserAddr, error) {
	const op = "App.UseService.UserAddrListAll"

	if err := uuid.Validate(uid); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.UserAddrSelectAll(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// CATLG__

func (s *UseService) ProductDetail(ctx context.Context, ref string) (*Product, error) {
	const op = "App.UseService.ProductDetail"

	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.ProductSelect(ctx, ref, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) ProductListAll(ctx context.Context, filter *ProductFilterInput) ([]*Product, error) {
	const op = "App.UseService.ProductListAll"

	f := NewProductFilterData(filter)
	if err := f.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results, err := s.deps.Repository.ProductSelectAll(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return results, nil
}

// ORDER__

func (s *UseService) OrderCreate(ctx context.Context, input *OrderInsertInput) (*Order, error) {
	const op = "App.UseService.OrderCreate"

	data := NewOrderInsertData(input)
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var order *Order
	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		ref, err := s.deps.Repository.OrderInsert(ctx, data)
		if err != nil {
			return err
		}

		order, err = s.deps.Repository.OrderSelect(ctx, ref, false)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

func (s *UseService) OrderUpdate(ctx context.Context, ref string, paths []string, input *OrderUpdateInput) (*Order, error) {
	const op = "App.UseService.OrderUpdate"

	log.Println("paths", paths)
	log.Println("input", input)
	
	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	updata := NewOrderUpdateData(input)
	if err := updata.Validate(paths); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.OrderSelect(ctx, ref, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		result, err = s.deps.Repository.OrderSelect(ctx, ref, true)
		if err != nil {
			return err
		}

		_, err = s.deps.Repository.OrderUpdate(ctx, ref, paths, updata)
		if err != nil {
			return err
		}

		result, err = s.deps.Repository.OrderSelect(ctx, ref, false)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) OrderDelete(ctx context.Context, ref string) (*Order, error) {
	const op = "App.UseService.OrderDelete"

	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.OrderSelect(ctx, ref, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if result.Status != "pending" {
		return nil, fmt.Errorf("%s: %w", op, WrapOrderDeleteNotAllowed(nil))
	}

	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		current, err := s.deps.Repository.OrderSelect(ctx, ref, true)
		if err != nil {
			return err
		}
		if current.Status != orderStatusPending {
			return WrapOrderDeleteNotAllowed(nil)
		}

		affected, err := s.deps.Repository.OrderDelete(ctx, ref)
		if err != nil {
			return err
		}
		if affected == 0 {
			return WrapOrderDeleteNotAllowed(nil)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) OrderDetail(ctx context.Context, ref string) (*Order, error) {
	const op = "App.UseService.OrderDetail"

	order, err := s.deps.Repository.OrderSelect(ctx, ref, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

func (s *UseService) OrderListAll(ctx context.Context, filter *OrderFilterInput, paging *OrderPagingInput) ([]*Order, error) {
	const op = "App.UseService.OrderListAll"

	// aqui se convierten
	f := NewOrderFilterData(filter)
	if err := f.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	p := NewOrderPagingData(paging)
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	orders, err := s.deps.Repository.OrderSelectAll(ctx, f, p)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (s *UseService) OrderChangeStatus(ctx context.Context, ref string, input *OrderChangeStatusInput) (*Order, error) {
	const op = "App.UseService.OrderChangeStatus"

	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data := NewOrderChangeStatusData(input)
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var order *Order
	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		var err error

		order, err = s.deps.Repository.OrderSelect(ctx, ref, true)
		if err != nil {
			return err
		}

		currentStatus := strings.TrimSpace(order.Status)
		nextStatus := data.Status

		if currentStatus == nextStatus {
			return nil
		}

		if !canTransitionOrderStatus(currentStatus, nextStatus) {
			return WrapOrderInvalidTransition(nil)
		}

		switch nextStatus {
		case orderStatusAcepted:
			lines, err := s.deps.Repository.OrderLineSelectAll(ctx, order.Ref)
			if err != nil {
				return err
			}
			if len(lines) == 0 {
				return WrapOrderLineEmpty(nil)
			}
		}

		status := &OrderUpdateData{Status: nextStatus}
		if _, err := s.deps.Repository.OrderUpdate(ctx, ref, []string{"status"}, status); err != nil {
			return err
		}

		order, err = s.deps.Repository.OrderSelect(ctx, ref, false)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

// ORDER_LINE__

func (s *UseService) OrderLineCreate(ctx context.Context, oid string, input *OrderLineCreateInput) (*OrderLine, error) {
	const op = "App.UseService.OrderLineCreate"

	if err := uuid.Validate(oid); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// # CHECK ORDER EXISTS #
	if _, err := s.deps.Repository.OrderSelect(ctx, oid, false); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data := NewOrderLineInsertData(input)
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var result *OrderLine
	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		ref, err := s.deps.Repository.OrderLineInsert(ctx, oid, data)
		if err != nil {
			return err
		}

		result, err = s.deps.Repository.OrderLineSelect(ctx, ref, false)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) OrderLineUpdate(ctx context.Context, ref string, paths []string, input *OrderLineUpdateInput) (*OrderLine, error) {
	const op = "App.UseService.OrderLineUpdate"

	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data := NewOrderLineUpdateData(input)
	if err := data.Validate(paths); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.OrderLineSelect(ctx, ref, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		_, err := s.deps.Repository.OrderLineUpdate(ctx, ref, paths, data)
		if err != nil {
			return err
		}

		result, err = s.deps.Repository.OrderLineSelect(ctx, ref, false)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) OrderLineDelete(ctx context.Context, ref string) (*OrderLine, error) {
	const op = "App.UseService.OrderLineDelete"

	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.OrderLineSelect(ctx, ref, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.deps.Repository.db.WithTx(ctx, func(ctx context.Context) error {
		_, err := s.deps.Repository.OrderLineSelect(ctx, ref, true)
		if err != nil {
			return err
		}

		if _, err := s.deps.Repository.OrderLineDelete(ctx, ref); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) OrderLineDetail(ctx context.Context, ref string) (*OrderLine, error) {
	const op = "App.UseService.OrderLineDetail"

	if err := uuid.Validate(ref); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.OrderLineSelect(ctx, ref, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) OrderLineListAll(ctx context.Context, oid string) ([]*OrderLine, error) {
	const op = "App.UseService.OrderLineListAll"

	if err := uuid.Validate(oid); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.OrderLineSelectAll(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// DELIVERY_DAY__

func (s *UseService) DeliveryDayListAll(ctx context.Context, filter *DeliveryDayFilterInput, paging *DeliveryDayPagingInput) ([]*DeliverySlot, error) {
	const op = "App.UseService.DeliveryDayListAll"

	f := NewDeliveryDayFilterData(filter)
	if err := f.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	p := NewDeliveryDayPagingData(paging)
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.Repository.DeliveryDaySelectAll(ctx, f, p)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// GOOGLE_MAPS__

func (s *UseService) PlaceDetail(ctx context.Context, input *gmaps.PlaceDetailInput) (*gmaps.Place, error) {
	const op = "App.UseService.PlaceDetail"

	data := gmaps.NewPlaceDetailData(input)
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := gmaps.ParseSessionToken(data.Token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.GoogleMapx.PlaceDetails(ctx, data.Ref, token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) ReverseGeocode(ctx context.Context, input *gmaps.ReverseGeocodeInput) (*gmaps.Place, error) {
	const op = "App.UseService.ReverseGeocode"

	data := gmaps.NewReverseGeocodeData(input)
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.deps.GoogleMapx.ReverseGeocode(ctx, data.Lat, data.Lng)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *UseService) PlaceAutocomplete(ctx context.Context, input *gmaps.PlaceAutocompleteInput) ([]*gmaps.Prediction, string, error) {
	const op = "App.UseService.PlaceAutocomplete"

	data := gmaps.NewPlaceAutocompleteData(input)
	if err := data.Validate(); err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	token := gmaps.NewSessionToken()
	results, err := s.deps.GoogleMapx.Autocomplete(ctx, data.Query, token)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return results, uuid.UUID(token).String(), nil
}
