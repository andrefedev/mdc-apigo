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
	SlotRef           string     `db:"slot_ref"`
	SlotKind          string     `db:"slot_kind"`
	SlotNote          *string    `db:"slot_note"`
	SlotWday          time.Time  `db:"slot_wday"`
	SlotIsOpen        bool       `db:"slot_is_open"`
	SlotCapacity      int32      `db:"slot_capacity"`
	SlotReserved      int32      `db:"slot_reserved"`
	SlotCutoffMin     int32      `db:"slot_cutoff_min"`
	SlotDateCreated   time.Time  `db:"slot_date_created"`
	SlotDateUpdated   *time.Time `db:"slot_date_updated"`
	SlotDeliveryStart int32      `db:"slot_delivery_start"`
	SlotDeliveryUntil int32      `db:"slot_delivery_until"`
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
		PaymentMethod: o.PaymentMethod,

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

		// SLOT__
		Slot: &DeliverySlot{
			Ref:           o.SlotRef,
			Kind:          o.SlotKind,
			Note:          o.SlotNote,
			Wday:          o.SlotWday,
			IsOpen:        o.SlotIsOpen,
			Capacity:      o.SlotCapacity,
			Reserved:      o.SlotReserved,
			CutoffMin:     o.SlotCutoffMin,
			DeliveryStart: o.SlotDeliveryStart,
			DeliveryUntil: o.SlotDeliveryUntil,
		},
	}
}

type OrderLineRaw struct {
	Ref        string `db:"id"`
	Status     string `db:"status"`
	Quantity   int32  `db:"quantity"`
	BasePrice  int32  `db:"base_price"`
	TotalPrice int32  `db:"total_price"`

	// PRODUCT__
	ProductRef     string  `db:"product_ref"`
	ProductUpc     *string `db:"product_upc"`
	ProductCode    int32   `db:"product_code"`
	ProductName    string  `db:"product_name"`
	ProductDescr   *string `db:"product_descr"`
	ProductImurl   *string `db:"product_imurl"`
	ProductDisplay int32   `db:"product_display"`

	ProductWeight   int32  `db:"product_weight"`
	ProductUnitype  string `db:"product_unitype"`
	ProductQuantity int32  `db:"product_quantity"`

	ProductIsActive    bool       `db:"product_is_active"`
	ProductIsPublic    bool       `db:"product_is_public"`
	ProductCostPrice   int32      `db:"product_cost_price"`
	ProductBasePrice   int32      `db:"product_base_price"`
	ProductNumInStock  int32      `db:"product_num_in_stock"`
	ProductNumInAlloc  int32      `db:"product_num_in_alloc"`
	ProductNumInAvail  int32      `db:"product_num_in_avail"`
	ProductDateCreated time.Time  `db:"product_date_created"`
	ProductDateUpdated *time.Time `db:"product_date_updated"`
}

func (o *OrderLineRaw) ToOrderLine() *OrderLine {
	return &OrderLine{
		Ref:        o.Ref,
		Status:     o.Status,
		Quantity:   o.Quantity,
		BasePrice:  o.BasePrice,
		TotalPrice: o.TotalPrice,

		Item: &Product{
			Ref:         o.ProductRef,
			Upc:         o.ProductUpc,
			Code:        o.ProductCode,
			Name:        o.ProductName,
			Descr:       o.ProductDescr,
			Imurl:       o.ProductImurl,
			Display:     o.ProductDisplay,
			Weight:      o.ProductWeight,
			Unitype:     o.ProductUnitype,
			Quantity:    o.ProductQuantity,
			IsActive:    o.ProductIsActive,
			IsPublic:    o.ProductIsPublic,
			CostPrice:   o.ProductCostPrice,
			BasePrice:   o.ProductBasePrice,
			NumInStock:  o.ProductNumInStock,
			NumInAlloc:  o.ProductNumInAlloc,
			NumInAvail:  o.ProductNumInAvail,
			DateCreated: o.ProductDateCreated,
			DateUpdated: o.ProductDateUpdated,
		},
	}
}

// PRODUCT__

type ProductRaw struct {
	Ref     string  `db:"id"`
	Upc     *string `db:"upc"`
	Code    int32   `db:"code"`
	Name    string  `db:"name"`
	Descr   *string `db:"descr"`
	Imurl   *string `db:"imurl"`
	Display int32   `db:"display"`

	Weight   int32  `db:"weight"`
	Unitype  string `db:"unitype"`
	Quantity int32  `db:"quantity"`

	IsActive    bool       `db:"is_active"`
	IsPublic    bool       `db:"is_public"`
	CostPrice   int32      `db:"cost_price"`
	BasePrice   int32      `db:"base_price"`
	NumInStock  int32      `db:"num_in_stock"`
	NumInAlloc  int32      `db:"num_in_alloc"`
	NumInAvail  int32      `db:"num_in_avail"`
	DateCreated time.Time  `db:"date_created"`
	DateUpdated *time.Time `db:"date_updated"`

	// GENRE__
	GenreRef         string    `db:"genre_ref"`
	GenreName        string    `db:"genre_name"`
	GenreDescr       *string   `db:"genre_descr"`
	GenreImurl       *string   `db:"genre_imurl"`
	GenreDisplay     int32     `db:"genre_display"`
	GenreIsPublic    bool      `db:"genre_is_public"`
	GenreDateCreated time.Time `db:"genre_date_created"`
}

func (p *ProductRaw) ToProduct() *Product {
	return &Product{
		Ref:     p.Ref,
		Upc:     p.Upc,
		Code:    p.Code,
		Name:    p.Name,
		Descr:   p.Descr,
		Imurl:   p.Imurl,
		Display: p.Display,

		Weight:   p.Weight,
		Unitype:  p.Unitype,
		Quantity: p.Quantity,

		IsActive:    p.IsActive,
		IsPublic:    p.IsPublic,
		CostPrice:   p.CostPrice,
		BasePrice:   p.BasePrice,
		NumInStock:  p.NumInStock,
		NumInAlloc:  p.NumInAlloc,
		DateCreated: p.DateCreated,
		DateUpdated: p.DateUpdated,

		Genre: &Genre{
			Ref:         p.GenreRef,
			Name:        p.GenreName,
			Descr:       p.GenreDescr,
			Imurl:       p.GenreImurl,
			Display:     p.GenreDisplay,
			IsPublic:    p.GenreIsPublic,
			DateCreated: p.GenreDateCreated,
		},
	}
}
