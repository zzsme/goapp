// Placeholder for paginator.go
package utils

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// Offset 获取分页偏移量
func (p *Pagination) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// Limit 获取分页大小
func (p *Pagination) Limit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

// SetTotal 设置总数
func (p *Pagination) SetTotal(total int) {
	p.Total = total
}
