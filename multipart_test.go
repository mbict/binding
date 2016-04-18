package binding

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strconv"

	. "gopkg.in/check.v1"
)

type multipartSuite struct{}

var _ = Suite(&multipartSuite{})

func (s *multipartSuite) Test_NotByReference(c *C) {
	post := Post{}
	req := newMultipartRequest(bytes.NewBufferString(""), "")
	err := MultipartForm.Bind(post, req)

	c.Assert(err, NotNil)
	c.Assert(err, DeepEquals, ErrorInputNotByReference)
}

func (s *multipartSuite) Test_NotAStruct(c *C) {
	test := int(1)
	req := newMultipartRequest(bytes.NewBufferString(""), "")
	err := MultipartForm.Bind(&test, req)

	c.Assert(err, NotNil)
	c.Assert(err, DeepEquals, ErrorInputIsNotStructure)
}

func (s *multipartSuite) Test_HappyPath(c *C) {
	blogPost := BlogPost{
		Post: Post{Title: "Glorious Post Title"}, Id: 1,
		Author:       Person{Name: "Matt Holt"},
		Coauthor:     &Person{Name: "The other guy"},
		Readers:      []Person{Person{Name: "Person a"}, Person{Name: "Person b", Email: "b@test.com"}},
		Contributors: []*Person{&Person{Name: "Michael Boke", Email: "mb@test.com"}, &Person{Name: "The other guy"}},
	}
	b, w := makeMultipartPayload(blogPost)
	req := newMultipartRequest(b, w.FormDataContentType())
	w.Close()

	response := BlogPost{}
	err := MultipartForm.Bind(&response, req)

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, blogPost)
}

func (s *multipartSuite) Test_FormValueCalledBeforeReader(c *C) {
	blogPost := BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}}
	b, w := makeMultipartPayload(blogPost)
	req := newMultipartRequest(b, w.FormDataContentType())
	w.Close()
	req.FormValue("foo") //called before multipart form
	response := BlogPost{}
	err := MultipartForm.Bind(&response, req)

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, blogPost)
}

func (s *multipartSuite) Test_EmptyBody(c *C) {
	blogPost := BlogPost{}
	b, w := makeMultipartPayload(blogPost)
	req := newMultipartRequest(b, w.FormDataContentType())
	w.Close()
	response := BlogPost{}
	err := MultipartForm.Bind(&response, req)

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, blogPost)
}

func (s *multipartSuite) Test_MissingRequiredFieldId(c *C) {
	blogPost := BlogPost{Post: Post{Title: "Glorious Post Title"}, Author: Person{Name: "Matt Holt"}}
	b, w := makeMultipartPayload(blogPost)
	req := newMultipartRequest(b, w.FormDataContentType())
	w.Close()
	response := BlogPost{}
	err := MultipartForm.Bind(&response, req)

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, blogPost)
}

func (s *multipartSuite) Test_MultipleValues(c *C) {
	blogPost := BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}, Ratings: []int{3, 5, 4}}
	b, w := makeMultipartPayload(blogPost)
	req := newMultipartRequest(b, w.FormDataContentType())
	w.Close()
	response := BlogPost{}
	err := MultipartForm.Bind(&response, req)

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, blogPost)
}

func (s *multipartSuite) Test_BadEncoding(c *C) {
	b, w := makeMalformedMultipartPayload()
	req := newMultipartRequest(b, "multipart/form-data")
	w.Close()
	response := BlogPost{}
	err := MultipartForm.Bind(&response, req)

	c.Assert(err, DeepEquals, ErrorDeserialization)
	c.Assert(response, DeepEquals, BlogPost{})
}

func makeMalformedMultipartPayload() (*bytes.Buffer, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	body.Write([]byte(`--` + writer.Boundary() + `\nContent-Disposition: form-data; name="foo"\n\n--` + writer.Boundary() + `--`))
	return body, writer
}

func makeMultipartPayload(blogPost BlogPost) (*bytes.Buffer, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("title", blogPost.Title)
	writer.WriteField("content", blogPost.Content)
	writer.WriteField("id", strconv.Itoa(blogPost.Id))
	writer.WriteField("ignored", blogPost.Ignored)
	for _, value := range blogPost.Ratings {
		writer.WriteField("rating", strconv.Itoa(value))
	}
	writer.WriteField("author.name", blogPost.Author.Name)
	writer.WriteField("author.email", blogPost.Author.Email)
	if blogPost.Coauthor != nil {
		writer.WriteField("coauthor.name", blogPost.Coauthor.Name)
		writer.WriteField("coauthor.email", blogPost.Coauthor.Email)
	}
	for key, value := range blogPost.Contributors {
		writer.WriteField("contributors."+strconv.Itoa(key)+".name", value.Name)
		writer.WriteField("contributors."+strconv.Itoa(key)+".email", value.Email)
	}
	for key, value := range blogPost.Readers {
		writer.WriteField("readers."+strconv.Itoa(key)+".name", value.Name)
		writer.WriteField("readers."+strconv.Itoa(key)+".email", value.Email)
	}
	return body, writer
}

func newMultipartRequest(multipart *bytes.Buffer, contentType string) *http.Request {
	req, err := http.NewRequest("POST", "", multipart)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", contentType)
	return req
}
