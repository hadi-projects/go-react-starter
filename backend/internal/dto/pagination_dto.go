package dto

type PaginationRequest struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
}

func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PaginationRequest) GetLimit() int {
	if p.Limit == 0 {
		return 10
	}
	return p.Limit
}

func (p *PaginationRequest) GetPage() int {
	if p.Page == 0 {
		return 1
	}
	return p.Page
}
