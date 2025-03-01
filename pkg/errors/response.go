package errors

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Message string `json:"message"`
}
