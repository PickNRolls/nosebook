package errors

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func New(errorType string, message string) *Error {
	return &Error{
		Type:    errorType,
		Message: message,
	}
}

func From(err error) *Error {
	e, ok := err.(*Error)
	if ok {
		return e
	}

	return New("Unknown", err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	return e.Message
}
