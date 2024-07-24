package payload

import "time"

type ProductFilter struct {
	Keyword       string     `form:"keyword"`
	ID            int64      `form:"id"`
	Name          string     `form:"name"`
	PriceFrom     float64    `form:"priceFrom"`
	PriceTo       float64    `form:"priceTo"`
	Description   string     `form:"description"`
	CategoryID    int64      `form:"categoryId"`
	CreatedAtFrom *time.Time `form:"createdAtFrom"`
	CreatedAtTo   *time.Time `form:"createdAtTo"`
	UpdatedAtFrom *time.Time `form:"updatedAtFrom"`
	UpdatedAtTo   *time.Time `form:"updatedAtTo"`
	Deleted       bool       `form:"deleted"`
}