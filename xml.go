package binding

import (
	"encoding/xml"
	"io"
	"net/http"
	"reflect"
)

type xmlBinding struct{}

func (_ xmlBinding) Name() string {
	return "xml"
}

func (_ xmlBinding) Bind(dst interface{}, req *http.Request) Errors {
	var bindErrors Errors

	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return append(bindErrors, ErrorInputNotByReference)
	}

	if req.Body != nil {
		defer req.Body.Close()
		err := xml.NewDecoder(req.Body).Decode(dst)
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
