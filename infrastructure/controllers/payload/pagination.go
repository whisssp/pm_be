package payload

type Pagination struct {
	Limit      int    `form:"limit"`
	Page       int    `form:"page"`
	Sort       string `form:"sort"`
	TotalRows  int64  `form:"-"`
	TotalPages int    `form:"-"`
}

func InitPaginate() *Pagination {
	return &Pagination{
		Page:  0,
		Limit: 10,
	}
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id asc"
	}
	return p.Sort
}