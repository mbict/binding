package binding

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

type jsonBinding struct{}

func (_ jsonBinding) Name() string {
	return "json"
}

// Json is middleware to deserialize a JSON payload from the request
// into the struct that is passed in. The resulting struct is then
// validated, but no error handling is actually performed here.
// An interface pointer can be added as a second argument in order
// to map the struct to a specific interface.
func (_ jsonBinding) Bind(dst interface{}, req *http.Request) Errors {

	var bindErrors Errors

	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return append(bindErrors, ErrorInputNotByReference)
	}

	if req.Body != nil {
		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(dst)
		if err != nil && err != io.EOF {
			bindErrors.Add([]string{}, DeserializationError, err.Error())
		}
	}

	validateErrs := validate(dst)
	if validateErrs != nil {
		return append(bindErrors, validateErrs...)
	}
	return bindErrors
}
