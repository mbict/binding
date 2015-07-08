[![Build Status](https://drone.io/github.com/mbict/binding/status.png)](https://drone.io/github.com/mbict/binding/latest)
[![Build Status](https://travis-ci.org/mbict/binding.png?branch=master)](https://travis-ci.org/mbict/binding)
[![Coverage Status](https://coveralls.io/repos/mbict/binding/badge.png)](https://coveralls.io/r/mbict/binding)
[![GoDoc](https://godoc.org/github.com/mbict/binding?status.png)](http://godoc.org/github.com/mbict/binding)
[![GoCover](http://gocover.io/_badge/github.com/mbict/binding)](http://gocover.io/github.com/mbict/binding)

binding
=======


### Installation

	go get github.com/mbict/binding
	
## Features

 - Automatically converts data from a request into a struct
 - Supports form, JSON, and multipart form data (including file uploads)
 - Provides data validation facilities
 	- Enforces required fields
 	- Invoke your own data validator

## Usage

#### Getting form data from a request

Suppose you have a contact form on your site where at least name and message are required. We'll need a struct to receive the data:

```go
type ContactForm struct {
	Name    string `form:"name" bind:"required"`
	Email   string `form:"email"`
	Message string `form:"message" bind:"required"`
}
```

In your http handle add 

```go
func(w http.ResponseWriter, r *http.Request) {
	contactForm := ContactForm{}
	errs := binding.Bind(&contactForm, r)
	...
}
```

or in case you need to use the pointer variant. Yes `bind` support (nil) pointer structs

```go
func(w http.ResponseWriter, r *http.Request) {
	contactForm := (*ContactForm)(nil))
	errs := binding.Bind(&contactForm, r)
	...
}
```

That's it! The `binding.Bind` function takes care of validating required fields. If there are any errors (like a required field is empty), `bind` will return all validation errors.


#### Getting JSON data from a request

To get data from JSON payloads, simply use the `json:` struct tags instead of `form:`. Pro Tip: Use [JSON-to-Go](http://mholt.github.io/json-to-go/) to correctly convert JSON to a Go type definition. It's useful if you're new to this or the structure is large/complex.


#### Custom validation

If you want additional validation beyond just checking required fields, your struct can implement the `binding.Validator` interface like so:

```go
func (cf ContactForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	if strings.Contains(cf.Message, "Go needs generics") {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Go has generics. They're called interfaces.",
		})
	}
	return errors
}
```

Now, any contact form submissions with "Go needs generics" in the message will return an error explaining your folly.


### Bind

`binding.Bind` is a convenient wrapper over the other handlers in this package. It does the following boilerplate for you:

 1. Deserializes request data into a struct
 2. Performs validation with `binding.Validate`
 3. Returns any validation errors

Content-Type will be used to know how to deserialize the requests.

### Form

`binding.Form` deserializes form data from the request, whether in the query string or as a form-urlencoded payload. It only does these things:

 1. Deserializes request data into a struct
 2. Performs validation with `binding.Validate`
 3. Returns any validation errors


### MultipartForm and file uploads

Like `binding.Form`, `binding.MultipartForm` deserializes form data from a request into the struct you pass in. Additionally, this will deserialize a POST request that has a form of *enctype="multipart/form-data"*. If the bound struct contains a field of type [`*multipart.FileHeader`](http://golang.org/pkg/mime/multipart/#FileHeader) (or `[]*multipart.FileHeader`), you also can read any uploaded files that were part of the form.

This handler does the following:

 1. Deserializes request data into a struct
 2. Performs validation with `binding.Validate`
 3. Returns any validation errors

#### MultipartForm example

```go
type UploadForm struct {
	Title      string                `form:"title"`
	TextUpload *multipart.FileHeader `form:"txtUpload"`
}

func(w http.ResponseWriter, r *http.Request) {
	uploadForm := UploadForm{}
	errs := binding.MultipartForm(&uploadForm, r)
	
	if 
	file, err := uf.TextUpload.Open()
	...
}

func main() {
	m := macaron.Classic()
	m.Post("/", binding.MultipartForm(UploadForm{}), uploadHandler(uf UploadForm) string {
		file, err := uf.TextUpload.Open()
		// ... you can now read the uploaded file
	})
	m.Run()
}
```

#### Structs and slices example

*Html post values*
```javascript
author.id = "1"
author.name = "J. Smith"
reviewers.1.id = "2"
reviewers.1.name = "A. Jolie"
reviewers.2.id = "3"
reviewers.2.name = "M. boke"
```

*Structure*

```go
type Person struct {
	Id			int			`form:"id"`
	Name		string 		`form:"name"`
}

type Book struct {
	Author    	Person		`form:"author"`
	Reviewers 	[]Person	`form:"reviewers"`
}

func(_ http.ResponseWriter, req *http.Request) {
	book := Book{}
	errs := binding.Bind(&book, req)
}
```

### Json

`binding.Json` deserializes JSON data in the payload of the request. It does the following things:

 1. Deserializes request data into a struct
 2. Performs validation with `binding.Validate`
 3. Returns any validation errors

