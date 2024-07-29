package entity

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

func (f *ProductFilter) IsNil() bool {
	return f.Keyword == "" &&
		f.ID == 0 &&
		f.Name == "" &&
		f.PriceFrom == 0 &&
		f.PriceTo == 0 &&
		f.Description == "" &&
		f.CategoryID == 0 &&
		f.CreatedAtTo == nil && f.CreatedAtFrom == nil &&
		f.UpdatedAtTo == nil && f.UpdatedAtFrom == nil &&
		f.Deleted == false
}

type CategoryFilter struct {
	Keyword       string     `form:"keyword"`
	ID            int64      `form:"id"`
	Name          string     `form:"categoryName"`
	CreatedAtFrom *time.Time `form:"createdAtFrom"`
	CreatedAtTo   *time.Time `form:"createdAtTo"`
	UpdatedAtFrom *time.Time `form:"updatedAtFrom"`
	UpdatedAtTo   *time.Time `form:"updatedAtTo"`
	Deleted       bool       `form:"deleted"`
}

type OrderFilter struct {
}