package binding

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

// These types are mostly contrived examples, but they're used
// across many test cases. The idea is to cover all the scenarios
// that this binding package might encounter in actual use.
type (
	// For basic test cases with a required field
	Post struct {
		Title   string `form:"title" json:"title"`
		Content string `form:"content" json:"content"`
	}

	// To be used as a nested struct (with a required field)
	Person struct {
		Name  string `form:"name" json:"name"`
		Email string `form:"email" json:"email"`
	}

	// For advanced test cases: multiple values, embedded
	// and nested structs, an ignored field, and single
	// and multiple file uploads
	BlogPost struct {
		Post
		Id           int                     `form:"id"`
		Ignored      string                  `form:"-" json:"-"`
		Ratings      []int                   `form:"rating" json:"ratings"`
		Author       Person                  `json:"author"`
		Coauthor     *Person                 `json:"coauthor"`
		Readers      []Person                `schema:"readers"`
		Contributors []*Person               `schema:"contributors"`
		HeaderImage  *multipart.FileHeader   `form:"headerImage"`
		Pictures     []*multipart.FileHeader `form:"picture"`
		unexported   string                  `form:"unexported"`
	}

	EmbedPerson struct {
		*Person
	}
)

func newRequest(method, query, body, contentType string) *http.Request {

	var bodyReader io.Reader
	if body != "-nil-" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, query, bodyReader)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", contentType)
	return req
}

const (
	formContentType = "application/x-www-form-urlencoded"
)
