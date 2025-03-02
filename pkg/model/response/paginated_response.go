package response

type PaginatedResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalCount int         `json:"totalCount"`
	Messages   interface{} `json:"messages"`
}
