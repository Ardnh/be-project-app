// internal/interfaces/http/dto/response/base_response.go
package dto

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

type PaginationResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    MetaData    `json:"meta"`
}

type MetaData struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Helper functions
func Success(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func Error(message string, errors ...string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}

func Paginated(message string, data interface{}, meta MetaData) PaginationResponse {
	return PaginationResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}
