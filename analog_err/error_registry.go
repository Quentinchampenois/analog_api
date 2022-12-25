package analog_err

import "strconv"

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

type AnalogError struct {
	Code   int
	Detail string
}

func (aErr *AnalogError) Error() string {
	return aErr.Detail
}

func (aErr *AnalogError) Display() map[string]string {
	return map[string]string{
		"code":  strconv.Itoa(aErr.Code),
		"error": aErr.Error(),
	}
}
