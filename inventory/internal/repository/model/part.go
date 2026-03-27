package model

import "time"

type Category int32

const (
	CategoryUnknown  Category = 0 // Неизвестная категория
	CategoryEngine   Category = 1 // Двигатель
	CategoryFuel     Category = 2 // Топливо
	CategoryPorthole Category = 3 // Иллюминатор
	CategoryWing     Category = 4 // Крыло
)

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type Manufacturer struct {
	Name    string
	Country string
	WebSite string
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type PartsFilter struct {
	UUIDS                   []string
	Names                   []string
	Categories              []Category
	ManunufacturerCountries []string
	Tags                    []string
}

type Value interface {
	isKind()
}

type (
	StringValue  struct{ V string }
	Int64Value   struct{ V int64 }
	Float64Value struct{ V float64 }
	BoolValue    struct{ V bool }
)

func (StringValue) isKind()  {}
func (Int64Value) isKind()   {}
func (Float64Value) isKind() {}
func (BoolValue) isKind()    {}
