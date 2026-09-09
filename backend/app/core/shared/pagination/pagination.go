package pagination

// Query is the normalized pagination input shared by list endpoints.
type Query struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// Meta is returned with every paginated response.
type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewMeta(query Query, total int) Meta {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 {
		query.PerPage = 20
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + query.PerPage - 1) / query.PerPage
	}
	return Meta{Page: query.Page, PerPage: query.PerPage, Total: total, TotalPages: totalPages}
}
