package models

type ResponseOK struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Error string `json:"error"`
}