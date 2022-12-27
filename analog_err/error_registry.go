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
			Code:    999,
			Message: "Unknown error occurred",
			Detail:  "Sorry for the inconvenient, an unexpected error occurred, we will check on this event.",
		}
	}

	return analogErr
}

type AnalogError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
	IsError bool   `json:"error"`
}

func (aErr *AnalogError) Error() string {
	return aErr.Detail
}
