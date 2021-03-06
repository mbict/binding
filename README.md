[![wercker status](https://app.wercker.com/status/69476f936eb47a347f8aa60af7c7c84e/s "wercker status")](https://app.wercker.com/project/bykey/69476f936eb47a347f8aa60af7c7c84e)
[![Build Status](https://travis-ci.org/mbict/go-binding.png?branch=master)](https://travis-ci.org/mbict/go-binding)
[![GoDoc](https://godoc.org/github.com/mbict/go-binding?status.png)](http://godoc.org/github.com/mbict/go-binding)
[![GoCover](http://gocover.io/_badge/github.com/mbict/go-binding)](http://gocover.io/github.com/mbict/go-binding)
[![GoReportCard](http://goreportcard.com/badge/mbict/go-binding)](http://goreportcard.com/report/mbict/go-binding)

binding
=======


### Installation

	go get github.com/mbict/binding
	
## Features

 - Automatically converts data from a request into a struct
 - Supports form, JSON, and multipart form data (including file uploads)

## Usage

#### Getting form data from a request

Suppose you have a contact form on your site. We'll need a struct to receive the data:

```go
type ContactForm struct {
	Name    string `form:"name"`
	Email   string `form:"email"`
	Message string `form:"message"`
}
```

In your http handle add 

```go
func(w http.ResponseWriter, r *http.Request) {
	contactForm := ContactForm{}
	err := binding.Bind(&contactForm, r)
	...
}
```

or in case you need to use the pointer variant. Yes `bind` support (nil) pointer structs

```go
func(w http.ResponseWriter, r *http.Request) {
	contactForm := (*ContactForm)(nil))
	err := binding.Bind(&contactForm, r)
	...
}
```

#### Getting JSON data from a request

To get data from JSON payloads, simply use the `json:` struct tags instead of `form:`. Pro Tip: Use [JSON-to-Go](http://mholt.github.io/json-to-go/) to correctly convert JSON to a Go type definition. It's useful if you're new to this or the structure is large/complex.

### Bind

`binding.Bind` is a convenient wrapper over the other handlers in this package.

Content-Type will be used to know how to deserialize the requests.

### Form

`binding.Form` deserializes form data from the request, whether in the query string or as a form-urlencoded payload.

### MultipartForm and file uploads

Like `binding.Form`, `binding.MultipartForm` deserializes form data from a request into the struct you pass in. Additionally, this will deserialize a POST request that has a form of *enctype="multipart/form-data"*. If the bound struct contains a field of type [`*multipart.FileHeader`](http://golang.org/pkg/mime/multipart/#FileHeader) (or `[]*multipart.FileHeader`), you also can read any uploaded files that were part of the form.

#### MultipartForm example

```go
type UploadForm struct {
	Title      string                `form:"title"`
	TextUpload *multipart.FileHeader `form:"txtUpload"`
}

func(w http.ResponseWriter, r *http.Request) {
	uploadForm := UploadForm{}
	err := binding.MultipartForm(&uploadForm, r)
	
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
	err := binding.Bind(&book, req)
}
```

### Json

`binding.Json` deserializes JSON data in the payload of the request to a provided structure.


