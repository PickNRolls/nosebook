package dto

type ResponseDTO struct {
	Errors []error `json:"errors"`
	Data   any     `json:"data"`
}
