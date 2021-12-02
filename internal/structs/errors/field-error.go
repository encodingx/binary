package errors

type fieldError struct {
	wordError
	fieldName string
}
