package dto

type Params struct {
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string
}

type GetProjectsParams struct {
	UserID    string `json:"user_id"`
	Search    string `json:"search,omitempty"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}
