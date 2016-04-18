package binding

import (
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	MIMEJSON      = "application/json"
	MIMEHTML      = "text/html"
	MIMEXML       = "application/xml"
	MIMEXML2      = "text/xml"
	MIMEPlain     = "text/plain"
	MIMEPOSTForm  = "application/x-www-form-urlencoded"
	MIMEMultipart = "multipart/form-data"
)

type Binding interface {
	Name() string
	Bind(interface{}, *http.Request) error
}

var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	MultipartForm = multipartBinding{}
)

func Default(method, contentType string) Binding {
	if method == "POST" || method == "PUT" || method == "PATCH" || contentType != "" {
		switch contentType {
		case MIMEMultipart:
			return MultipartForm
		case MIMEPOSTForm:
			return Form
		case MIMEJSON:
			return JSON
		case MIMEXML, MIMEXML2:
			return XML
		default:
			/*if contentType == "" {
				return Errors{ErrorEmptyContentType}
			} else {
				return Errors{ErrorUnsupportedContentType}
			}*/
			return Form
		}
	}
	return Form
}

var (
	ErrorEmptyContentType       = NewError([]string{}, ContentTypeError, "Empty Content-Type")
	ErrorUnsupportedContentType = NewError([]string{}, ContentTypeError, "Unsupported Content-Type")
	ErrorInputNotByReference    = NewError([]string{}, DeserializationError, "input binding model is not by reference")
	ErrorInputIsNotStructure    = NewError([]string{}, DeserializationError, "binding model is required to be structure")
)

func Bind(obj interface{}, req *http.Request) error {
	contentType := req.Header.Get("Content-Type")
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" || contentType != "" {
		if strings.Contains(contentType, "form-urlencoded") {
			return Form.Bind(obj, req)
		} else if strings.Contains(contentType, "multipart/form-data") {
			return MultipartForm.Bind(obj, req)
		} else if strings.Contains(contentType, "json") {
			return JSON.Bind(obj, req)
		} else {
			if contentType == "" {
				return ErrorEmptyContentType
			} else {
				return ErrorUnsupportedContentType
			}
		}
	} else {
		return Form.Bind(obj, req)
	}
}

/*
var (
	alphaDashPattern    = regexp.MustCompile("[^\\d\\w-_]")
	alphaDashDotPattern = regexp.MustCompile("[^\\d\\w-_\\.]")
	emailPattern        = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")
	urlPattern          = regexp.MustCompile(`(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
)
*/
/*
// Performs required field checking on a struct
func validateStruct(errors Errors, obj interface{}, path string) Errors {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Allow ignored fields in the struct
		if field.Tag.Get("form") == "-" || !val.Field(i).CanInterface() {
			continue
		}

		fieldVal := val.Field(i)
		fieldValue := fieldVal.Interface()
		zero := reflect.Zero(field.Type).Interface()

		// Validate nested and embedded structs (if pointer, only do so if not nil)
		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue) &&
				field.Type.Elem().Kind() == reflect.Struct) {
			fieldPath := path
			if field.Anonymous == false {
				fieldPath = path + field.Name + "."
			}
			errors = validateStruct(errors, fieldValue, fieldPath)
			// Validate structure slices
		} else if field.Type.Kind() == reflect.Slice &&
			(field.Type.Elem().Kind() == reflect.Struct ||
				(field.Type.Elem().Kind() == reflect.Ptr && field.Type.Elem().Elem().Kind() == reflect.Struct)) {
			for i := 0; i < fieldVal.Len(); i++ {
				fieldPath := path + field.Name + "." + strconv.Itoa(i) + "."
				errors = validateStruct(errors, fieldVal.Index(i).Interface(), fieldPath)
			}
		}

		// Match rules.
	VALIDATE_RULES:
		for _, rule := range strings.Split(field.Tag.Get("validate"), ";") {
			if len(rule) == 0 {
				continue
			}

			switch {
			case rule == "Required":
				if reflect.DeepEqual(zero, fieldValue) {
					errors.Add([]string{path + field.Name}, RequiredError, "Required")
					break
				}
			case rule == "AlphaDash":
				if alphaDashPattern.MatchString(fmt.Sprintf("%v", fieldValue)) {
					errors.Add([]string{path + field.Name}, AlphaDashError, "AlphaDash")
					break VALIDATE_RULES
				}
			case rule == "AlphaDashDot":
				if alphaDashDotPattern.MatchString(fmt.Sprintf("%v", fieldValue)) {
					errors.Add([]string{path + field.Name}, AlphaDashDotError, "AlphaDashDot")
					break VALIDATE_RULES
				}
			case strings.HasPrefix(rule, "MinSize("):
				min, _ := strconv.Atoi(rule[8 : len(rule)-1])
				if str, ok := fieldValue.(string); ok && utf8.RuneCountInString(str) < min {
					errors.Add([]string{path + field.Name}, MinSizeError, "MinSize")
					break VALIDATE_RULES
				}
				v := reflect.ValueOf(fieldValue)
				if v.Kind() == reflect.Slice && v.Len() < min {
					errors.Add([]string{path + field.Name}, MinSizeError, "MinSize")
					break VALIDATE_RULES
				}
			case strings.HasPrefix(rule, "MaxSize("):
				max, _ := strconv.Atoi(rule[8 : len(rule)-1])
				if str, ok := fieldValue.(string); ok && utf8.RuneCountInString(str) > max {
					errors.Add([]string{path + field.Name}, MaxSizeError, "MaxSize")
					break VALIDATE_RULES
				}
				v := reflect.ValueOf(fieldValue)
				if v.Kind() == reflect.Slice && v.Len() > max {
					errors.Add([]string{path + field.Name}, MaxSizeError, "MaxSize")
					break VALIDATE_RULES
				}
			case rule == "Email":
				if !emailPattern.MatchString(fmt.Sprintf("%v", fieldValue)) {
					errors.Add([]string{path + field.Name}, EmailError, "Email")
					break VALIDATE_RULES
				}
			case rule == "Url":
				str := fmt.Sprintf("%v", fieldValue)
				if len(str) == 0 {
					continue
				} else if !urlPattern.MatchString(str) {
					errors.Add([]string{path + field.Name}, UrlError, "Url")
					break VALIDATE_RULES
				}

			// TODO write test for these validation rules
			case strings.HasPrefix(rule, "Range("):
				nums := strings.Split(rule[6:len(rule)-1], ",")
				if len(nums) != 2 {
					break
				}
				val, _ := strconv.ParseInt(fmt.Sprintf("%v", fieldValue), 10, 32)
				a, _ := strconv.ParseInt(nums[0], 10, 32)
				b, _ := strconv.ParseInt(nums[1], 10, 32)
				if val < a || val > b {
					errors.Add([]string{path + field.Name}, RangeError, "Range")
					break VALIDATE_RULES
				}
			case strings.HasPrefix(rule, "In("):
				if !in(fieldValue, rule[3:len(rule)-1]) {
					errors.Add([]string{path + field.Name}, InError, "In")
					break VALIDATE_RULES
				}
			case strings.HasPrefix(rule, "NotIn("):
				if in(fieldValue, rule[6:len(rule)-1]) {
					errors.Add([]string{path + field.Name}, NotInError, "NotIn")
					break VALIDATE_RULES
				}
			case strings.HasPrefix(rule, "Include("):
				if !strings.Contains(fmt.Sprintf("%v", fieldValue), rule[8:len(rule)-1]) {
					errors.Add([]string{path + field.Name}, IncludeError, "Include")
					break VALIDATE_RULES
				}
			case strings.HasPrefix(rule, "Exclude("):
				if strings.Contains(fmt.Sprintf("%v", fieldValue), rule[8:len(rule)-1]) {
					errors.Add([]string{path + field.Name}, ExcludeError, "Exclude")
					break
				}
			case strings.HasPrefix(rule, "Default("):
				if reflect.DeepEqual(zero, fieldValue) {
					if fieldVal.CanAddr() {
						setWithProperType(field.Type.Kind(), rule[8:len(rule)-1], fieldVal, field.Tag.Get("form"), errors)
					} else {
						errors.Add([]string{path + field.Name}, DefaultError, "Default")
						break VALIDATE_RULES
					}
				}

			}
		}
	}
	return errors
}
*/
//validation in function
/*
func in(fieldValue interface{}, arr string) bool {
	val := fmt.Sprintf("%v", fieldValue)
	vals := strings.Split(arr, ",")
	isIn := false
	for _, v := range vals {
		if v == val {
			isIn = true
			break
		}
	}
	return isIn
}
*/
/*
func mapFormValues(field string, form map[string][]string) (result []map[string][]string) {
	for key, values := range form {
		parts := strings.SplitN(key, ".", 3)
		if len(parts) == 3 && parts[0] == field {
			index, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			if len(result) < index+1 {
				tmp := make([]map[string][]string, index+1)
				copy(tmp, result)
				result = tmp
			}

			if result[index] == nil {
				result[index] = make(map[string][]string)
			}
			result[index][parts[2]] = values
		}
	}
	return result
}*/

func pathSliceSize(field string, form map[string][]string) int {
	size := 0
	for key, _ := range form {
		parts := strings.SplitN(key, ".", 3)
		if len(parts) == 3 && parts[0] == field {
			index, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			if size < index+1 {
				size = index + 1
			}
		}
	}
	return size
}

var fhType = reflect.TypeOf((*multipart.FileHeader)(nil))

// Takes values from the form data and puts them into a struct
func mapForm(path string, formStruct reflect.Value, form map[string][]string, formfile map[string][]*multipart.FileHeader) error {
	formStruct = reflect.Indirect(formStruct)
	typ := formStruct.Type()

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := formStruct.Field(i)

		inputFieldName := typeField.Tag.Get("form")
		if inputFieldName == "" {
			inputFieldName = strings.ToLower(typeField.Name)
		}

		if typeField.Anonymous {
			if typeField.Type.Kind() == reflect.Ptr {
				structField.Set(reflect.New(typeField.Type.Elem()))
				if err := mapForm(path, structField.Elem(), form, formfile); err != nil {
					return err
				}
				if reflect.DeepEqual(structField.Elem().Interface(), reflect.Zero(structField.Elem().Type()).Interface()) {
					structField.Set(reflect.Zero(structField.Type()))
				}
			} else {
				if err := mapForm(path, structField, form, formfile); err != nil {
					return err
				}
			}
		} else if structField.Kind() == reflect.Slice && structField.Type().Elem() == fhType {
			//slice of file uploads
			inputFile, exists := formfile[path+inputFieldName]
			if exists {
				numFiles := len(inputFile)
				if numFiles > 0 {
					slice := reflect.MakeSlice(structField.Type(), numFiles, numFiles)
					for i := 0; i < numFiles; i++ {
						slice.Index(i).Set(reflect.ValueOf(inputFile[i]))
					}
					structField.Set(slice)
				}
			}
		} else if structField.Type() == fhType {
			//single file
			inputFile, exists := formfile[path+inputFieldName]
			if exists && len(inputFile) >= 1 {
				structField.Set(reflect.ValueOf(inputFile[0]))
			}
		} else if typeField.Type.Kind() == reflect.Ptr && typeField.Type.Elem().Kind() == reflect.Struct {
			//find if we have posted this field and or need to init the pointer
			for key, _ := range form {
				if strings.HasPrefix(key, path+inputFieldName+".") {
					if structField.IsNil() {
						structField.Set(reflect.New(typeField.Type.Elem()))
					}
					if err := mapForm(path+inputFieldName+".", structField.Elem(), form, formfile); err != nil {
						return err
					}
					break
				}
			}
		} else if typeField.Type.Kind() == reflect.Struct {
			if err := mapForm(path+inputFieldName+".", structField, form, formfile); err != nil {
				return err
			}
		} else if typeField.Type.Kind() == reflect.Slice &&
			(typeField.Type.Elem().Kind() == reflect.Struct ||
				(typeField.Type.Elem().Kind() == reflect.Ptr && typeField.Type.Elem().Elem().Kind() == reflect.Struct)) {

			//size slice (if necessary)
			size := pathSliceSize(path+inputFieldName, form)
			if structField.Len() < size {
				value := reflect.MakeSlice(structField.Type(), size, size)
				if structField.Len() > 0 {
					reflect.Copy(value, structField)
				}
				structField.Set(value)
			}

			//assign structs
			for i := 0; i < size; i++ {
				sliceValue := structField.Index(i)
				if sliceValue.Kind() == reflect.Ptr && sliceValue.IsNil() {
					sliceValue.Set(reflect.New(sliceValue.Type().Elem()))
				}
				if err := mapForm(path+inputFieldName+"."+strconv.Itoa(i)+".", sliceValue, form, formfile); err != nil {
					return err
				}
			}

		} else if inputFieldName := typeField.Tag.Get("form"); inputFieldName != "" {
			if !structField.CanSet() {
				continue
			}

			inputValue, exists := form[path+inputFieldName]
			if exists {
				numElems := len(inputValue)
				if structField.Kind() == reflect.Slice && numElems > 0 {
					sliceOf := structField.Type().Elem().Kind()
					slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
					for i := 0; i < numElems; i++ {
						setWithProperType(sliceOf, inputValue[i], slice.Index(i), inputFieldName, errors)
					}
					formStruct.Field(i).Set(slice)
				} else {
					setWithProperType(typeField.Type.Kind(), inputValue[0], structField, inputFieldName, errors)
				}
			}
		}
	}
}

// This sets the value in a struct of an indeterminate type to the
// matching value from the request (via Form middleware) in the
// same type, so that not all deserialized values have to be strings.
// Supported types are string, int, float, and bool.
func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value, nameInTag string, errors Errors) {
	switch valueKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			val = "0"
		}
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			errors.Add([]string{nameInTag}, IntegerTypeError, "Value could not be parsed as integer")
		} else {
			structField.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val == "" {
			val = "0"
		}
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			errors.Add([]string{nameInTag}, IntegerTypeError, "Value could not be parsed as unsigned integer")
		} else {
			structField.SetUint(uintVal)
		}
	case reflect.Bool:
		if val == "on" {
			structField.SetBool(true)
			return
		}

		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			errors.Add([]string{nameInTag}, BooleanTypeError, "Value could not be parsed as boolean")
		} else {
			structField.SetBool(boolVal)
		}
	case reflect.Float32:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			errors.Add([]string{nameInTag}, FloatTypeError, "Value could not be parsed as 32-bit float")
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.Float64:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			errors.Add([]string{nameInTag}, FloatTypeError, "Value could not be parsed as 64-bit float")
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.String:
		structField.SetString(val)
	}
}

/*
// validate by the build in validation rules and tries to run the model ValidateBinder function if set
func validate(obj interface{}) Errors {
	if obj == nil {
		return nil
	}

	var bindErrors Errors
	v := reflect.ValueOf(obj)
	k := v.Kind()
	if k == reflect.Interface || k == reflect.Ptr {
		//skip nil pointers
		v = v.Elem()
		k = v.Kind()
	}

	if k == reflect.Slice || k == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i).Interface()
			bindErrors = validateStruct(bindErrors, e, strconv.Itoa(i)+".")
			if validator, ok := e.(Validator); ok {
				bindErrors = validator.Validate(bindErrors)
			}
		}
	} else {
		bindErrors = validateStruct(bindErrors, obj, "")
		if validator, ok := obj.(Validator); ok {
			bindErrors = validator.Validate(bindErrors)
		}
	}
	return bindErrors
}
*/

/*
type (
	// Implement the Validator interface to handle some rudimentary
	// request validation logic so your application doesn't have to.
	Validator interface {
		// ValidateBinder validates that the request is OK. It is recommended
		// that validation be limited to checking values for syntax and
		// semantics, enough to know that you can make sense of the request
		// in your application. For example, you might verify that a credit
		// card number matches a valid pattern, but you probably wouldn't
		// perform an actual credit card authorization here.
		Validate(Errors) Errors
	}
)
*/

var (
	// Maximum amount of memory to use when parsing a multipart form.
	// Set this to whatever value you prefer; default is 16 MB.
	MaxMemory = int64(1024 * 1024 * 16)
)

const (
	jsonContentType = "application/json; charset=utf-8"
)
