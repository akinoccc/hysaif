package types

// 通用分页类型
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// 通用响应类型
type ListResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Error   string              `json:"error"`
	Details map[string][]string `json:"details,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
