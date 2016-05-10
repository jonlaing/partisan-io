package modelsv2

import "encoding/json"

// ValidationErrors is a slice of errors for when model
// initialization fails validation tests. It implements the
// Error interface.
type ValidationErrors map[string]error

func (errs ValidationErrors) Error() string {
	m := make(map[string]string)
	for k, e := range errs {
		m[k] = e.Error()
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "An Unknown error has occurred"
	}

	return string(b)
}
