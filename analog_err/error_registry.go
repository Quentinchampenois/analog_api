package analog_err

type ErrorRegistry struct {
	Registry []*AnalogError
}

func (er *ErrorRegistry) Register(e *AnalogError) {
	er.Registry = append(er.Registry, e)
}

func (er *ErrorRegistry) Find(code int) *AnalogError {
	for _, analogErr := range er.Registry {
		if analogErr.Code == code {
			return analogErr
		}
	}
	return nil
}

func (er *ErrorRegistry) FindOrUnknown(code int) *AnalogError {
	var analogErr *AnalogError
	if analogErr = er.Find(code); analogErr == nil {
		analogErr = &AnalogError{
			Code:   999,
			Detail: "An unknown error occurred",
		}
	}

	return analogErr
}

type AnalogError struct {
	Code   int
	Detail string
}

func (aErr *AnalogError) Error() string {
	return aErr.Detail
}
