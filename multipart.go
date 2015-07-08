package binding

import (
	"net/http"
	"reflect"
)

type multipartBinding struct{}

func (_ multipartBinding) Name() string {
	return "multipart"
}

// MultipartForm works much like Form, except it can parse multipart forms
// and handle file uploads. Like the other deserialization middleware handlers,
// you can pass in an interface to make the interface available for injection
// into other handlers later.
func (_ multipartBinding) Bind(dst interface{}, req *http.Request) Errors {
	var bindErrors Errors

	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return append(bindErrors, ErrorInputNotByReference)
	}

	//reset element to zero variant
	v = v.Elem()
	if v.Kind() == reflect.Ptr && v.CanSet() && v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct || !v.CanSet() {
		return append(bindErrors, ErrorInputIsNotStructure)
	}

	// This if check is necessary due to https://github.com/martini-contrib/csrf/issues/6
	if req.MultipartForm == nil {
		// Workaround for multipart forms returning nil instead of an error
		// when content is not multipart; see https://code.google.com/p/go/issues/detail?id=6334
		if multipartReader, err := req.MultipartReader(); err != nil {
			// TODO: Cover this and the next error check with tests
			bindErrors.Add([]string{}, DeserializationError, err.Error())
		} else {
			form, parseErr := multipartReader.ReadForm(MaxMemory)
			if parseErr != nil {
				bindErrors.Add([]string{}, DeserializationError, parseErr.Error())
			}
			req.MultipartForm = form
		}
	}

	mapForm("", v, req.MultipartForm.Value, req.MultipartForm.File, bindErrors)
	validateErrs := validate(v.Interface())
	if validateErrs != nil {
		return append(bindErrors, validateErrs...)
	}
	return bindErrors
}
