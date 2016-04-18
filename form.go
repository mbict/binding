package binding

import (
	"net/http"
	"reflect"
)

type formBinding struct{}

func (_ formBinding) Name() string {
	return "form"
}

// Form is middleware to deserialize form-urlencoded data from the request.
// It gets data from the form-urlencoded body, if present, or from the
// query string. It uses the http.Request.ParseForm() method
// to perform deserialization, then reflection is used to map each field
// into the struct with the proper type. Structs with primitive slice types
// (bool, float, int, string) can support deserialization of repeated form
// keys, for example: key=val1&key=val2&key=val3
// An interface pointer can be added as a second argument in order
// to map the struct to a specific interface.
func (_ formBinding) Bind(dst interface{}, req *http.Request) error {

	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return ErrorInputNotByReference
	}

	//reset element to zero variant
	v = v.Elem()
	if v.Kind() == reflect.Ptr && v.CanSet() && v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct || !v.CanSet() {
		return ErrorInputIsNotStructure
	}

	// Format validation of the request body or the URL would add considerable overhead,
	// and ParseForm does not complain when URL encoding is off.
	// Because an empty request body or url can also mean absence of all needed values,
	// it is not in all cases a bad request, so let's return 422.
	parseErr := req.ParseForm()
	if parseErr != nil {
		return ErrorDeserialization
	}
	return mapForm("", v, req.Form, nil)
}
