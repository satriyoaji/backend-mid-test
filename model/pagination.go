package model

type PageRequest struct {
	PageNum  int `query:"page_num"`
	PageSize int `query:"page_size"`
}

type Pagination struct {
	PageNum   *int `json:"page_num,omitempty"`
	PageSize  *int `json:"page_size,omitempty"`
	TotalData *int `json:"total_data,omitempty"`
}
