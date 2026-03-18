package wabas

import "time"

type WabaMessage struct {
	Ref         string    `db:"id"`
	Uid         string    `db:"uid"`
	Wid         string    `db:"wid"`
	Status      string    `db:"status"`
	Category    string    `db:"category"`
	Template    string    `db:"template"`
	DateCreated time.Time `db:"date_created"`
	DateExpired time.Time `db:"date_expired"`
}

type WabaMessageEvent struct {
	Ref         string    `db:"id"`
	Event       string    `db:"event"`
	Status      string    `db:"status"`
	DateCreated time.Time `db:"date_created"`
}
