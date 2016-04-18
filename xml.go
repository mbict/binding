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

func (_ xmlBinding) Bind(dst interface{}, req *http.Request) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return ErrorInputNotByReference
	}

	if req.Body != nil {
		defer req.Body.Close()
		err := xml.NewDecoder(req.Body).Decode(dst)
		if err != nil && err != io.EOF {
			return ErrorDeserialization
		}
	}
	return nil
}
