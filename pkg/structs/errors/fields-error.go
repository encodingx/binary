package errors

type fieldsError struct {
	wordError
	leftFieldName  string
	rightFieldName string
}
