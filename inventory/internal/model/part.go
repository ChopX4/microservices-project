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
	UUID          string         `bson:"uuid"`
	Name          string         `bson:"name"`
	Description   string         `bson:"description"`
	Price         float64        `bson:"price"`
	StockQuantity int64          `bson:"stock_quantity"`
	Category      Category       `bson:"category"`
	Dimensions    Dimensions     `bson:"dimensions"`
	Manufacturer  Manufacturer   `bson:"manufacturer"`
	Tags          []string       `bson:"tags"`
	Metadata      map[string]any `bson:"metadata"`
	CreatedAt     *time.Time     `bson:"created_at"`
	UpdatedAt     *time.Time     `bson:"updated_at"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	WebSite string `bson:"website"`
}

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type PartsFilter struct {
	UUIDS                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

func (c Category) IsValid() bool {
	switch c {
	case CategoryEngine, CategoryFuel, CategoryPorthole, CategoryWing:
		return true
	default:
		return false
	}
}
