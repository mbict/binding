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
		Title   string `form:"title" json:"title" validate:"Required"`
		Content string `form:"content" json:"content"`
	}

	// To be used as a nested struct (with a required field)
	Person struct {
		Name  string `form:"name" json:"name" validate:"Required"`
		Email string `form:"email" json:"email"`
	}

	// For advanced test cases: multiple values, embedded
	// and nested structs, an ignored field, and single
	// and multiple file uploads
	BlogPost struct {
		Post
		Id           int                     `form:"id" validate:"Required"`
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

	SadForm struct {
		AlphaDash    string   `form:"AlphaDash" validate:"AlphaDash"`
		AlphaDashDot string   `form:"AlphaDashDot" validate:"AlphaDashDot"`
		MinSize      string   `form:"MinSize" validate:"MinSize(5)"`
		MinSizeSlice []string `form:"MinSizeSlice" validate:"MinSize(5)"`
		MaxSize      string   `form:"MaxSize" validate:"MaxSize(1)"`
		MaxSizeSlice []string `form:"MaxSizeSlice" validate:"MaxSize(1)"`
		Email        string   `form:"Email" validate:"Email"`
		Url          string   `form:"Url" validate:"Url"`
		UrlEmpty     string   `form:"UrlEmpty" validate:"Url"`
		Range        int      `form:"Range" validate:"Range(1,2)"`
		RangeInvalid int      `form:"RangeInvalid" validate:"Range(1)"`
		In           string   `form:"In" validate:"Default(0);In(1,2,3)"`
		InInvalid    string   `form:"InInvalid" validate:"In(1,2,3)"`
		NotIn        string   `form:"NotIn" validate:"NotIn(1,2,3)"`
		Include      string   `form:"Include" validate:"Include(a)"`
		Exclude      string   `form:"Exclude" validate:"Exclude(a)"`
	}
)

func (p Post) Validate(errors Errors) Errors {
	if len(p.Title) < 10 {
		errors = append(errors, Error{
			FieldNames:     []string{"Title"},
			Classification: "LengthError",
			Message:        "Life is too short",
		})
	}
	return errors
}

func (p EmbedPerson) Validate(errors Errors) Errors {
	if len(p.Email) <= 0 {
		errors = append(errors, Error{
			FieldNames:     []string{"Email"},
			Classification: "LengthError",
			Message:        "Email is too short",
		})
	}
	return errors
}

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
