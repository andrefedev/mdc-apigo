package app

import "time"

type OrderRaw struct {
	Ref           string     `db:"id"`
	Number        int32      `db:"number"`
	Status        string     `db:"status"`
	BasePrice     int32      `db:"base_price"`
	DiscPrice     int32      `db:"disc_price"`
	DateCreated   time.Time  `db:"date_created"`
	DateUpdated   *time.Time `db:"date_updated"`
	DeliveryDate  time.Time  `db:"delivery_date"`
	PaymentStatus string     `db:"payment_status"`
	PaymentMethod string     `db:"payment_method"`

	// USER__
	UserRef   string `db:"user_ref"`
	UserName  string `db:"user_name"`
	UserPhone string `db:"user_phone"`

	// ADDR__
	AddrRef      string  `db:"addr_ref"`
	AddrLat      float64 `db:"addr_lat"`
	AddrLng      float64 `db:"addr_lng"`
	AddrName     string  `db:"addr_name"`
	AddrCmna     string  `db:"addr_cmna"`
	AddrRoute    string  `db:"addr_route"`
	AddrStreet   string  `db:"addr_street"`
	AddrNeighb   string  `db:"addr_neighb"`
	AddrLocality string  `db:"addr_locality"`
	AddrSublocal string  `db:"addr_sublocal"`
	AddrAddress1 string  `db:"addr_address1"`
	AddrAddress2 string  `db:"addr_address2"`

	// SLOT__
	SlotRef         string    `db:"slot_ref"`
	SlotCode        string    `db:"slot_code"`
	SlotWorkDate    time.Time `db:"slot_work_date"`
	SlotCapacity    int32     `db:"slot_capacity"`
	SlotReserved    int32     `db:"slot_reserved"`
	SlotRemaining   int32     `db:"slot_remaining"`
	SlotStartUnix   int64     `db:"slot_start_unix"` //  segundos epoch
	SlotUntilUnix   int64     `db:"slot_until_unix"` //  segundos epoch
	SlotIsAvailable bool      `db:"slot_is_available"`
}

func (o *OrderRaw) ToOrder() *Order {
	return &Order{
		Ref:           o.Ref,
		Number:        o.Number,
		Status:        o.Status,
		BasePrice:     o.BasePrice,
		DiscPrice:     o.DiscPrice,
		DateCreated:   o.DateCreated,
		DeliveryDate:  o.DeliveryDate,
		DateUpdated:   o.DateUpdated,
		PaymentStatus: o.PaymentStatus,

		// USER__
		User: &User{
			Ref:   o.UserRef,
			Name:  o.UserName,
			Phone: o.UserPhone,
		},

		// ADDR__
		Addr: &UserAddr{
			Ref:    o.AddrRef,
			Lat:    o.AddrLat,
			Lng:    o.AddrLng,
			Name:   o.AddrName,
			Cmna:   o.AddrCmna,
			Route:  o.AddrRoute,
			Street: o.AddrStreet,
			Neighb: o.AddrNeighb,
		},
	}
}
