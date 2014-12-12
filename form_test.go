package binder

import (
	"net/http"
	"strings"

	. "gopkg.in/check.v1"
)

type formSuite struct{}

var _ = Suite(&formSuite{})

func (s *formSuite) Test_NotByReference(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, ``, formContentType)
	errs := Form(post, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs[0], DeepEquals, ErrorInputNotByReference)
}

func (s *formSuite) Test_NotAStruct(c *C) {
	test := int(1)
	req := newRequest(`POST`, ``, ``, formContentType)
	errs := Form(&test, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs[0], DeepEquals, ErrorInputIsNotStructure)
}

func (s *formSuite) Test_HappyPath(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, formContentType)
	errs := Form(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_HappyPathWithPointer(c *C) {
	post := &Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, formContentType)
	errs := Form(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, &Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_HappyPathWithNullPointer(c *C) {
	post := (*Post)(nil)
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, formContentType)
	errs := Form(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, &Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_EmptyBody(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, ``, formContentType)
	errs := Form(&post, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs[0], DeepEquals, Error{FieldNames: []string{"Title"}, Classification: "RequiredError", Message: "Required"})
	c.Assert(errs[1], DeepEquals, Error{FieldNames: []string{"Title"}, Classification: "LengthError", Message: "Life is too short"})
	c.Assert(post, DeepEquals, Post{})
}

func (s *formSuite) Test_EmptyContentType(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``)
	errs := Form(&post, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs[0], DeepEquals, Error{FieldNames: []string{"Title"}, Classification: "RequiredError", Message: "Required"})
	c.Assert(errs[1], DeepEquals, Error{FieldNames: []string{"Title"}, Classification: "LengthError", Message: "Life is too short"})
	c.Assert(post, DeepEquals, Post{})
}

func (s *formSuite) Test_MalformedBody(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=%2`, formContentType)
	errs := Form(&post, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 3)
	c.Assert(errs[0], DeepEquals, Error{FieldNames: []string{}, Classification: DeserializationError, Message: `invalid URL escape "%2"`})
	c.Assert(errs[1], DeepEquals, Error{FieldNames: []string{"Title"}, Classification: "RequiredError", Message: "Required"})
	c.Assert(errs[2], DeepEquals, Error{FieldNames: []string{"Title"}, Classification: "LengthError", Message: "Life is too short"})
	c.Assert(post, DeepEquals, Post{})
}

func (s *formSuite) Test_NestedEmbeddedStructs(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&id=1&name=Matt+Holt`, formContentType)
	errs := Form(&blogPost, req)

	c.Assert(errs, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *formSuite) Test_RequiredEmbeddedStructFieldNotSpecified(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `id=1&name=Matt+Holt`, formContentType)
	errs := Form(&blogPost, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs[0], DeepEquals, Error{FieldNames: []string{"Title"}, Classification: "RequiredError", Message: "Required"})
	c.Assert(errs[1], DeepEquals, Error{FieldNames: []string{"Title"}, Classification: "LengthError", Message: "Life is too short"})
	c.Assert(blogPost, DeepEquals, BlogPost{Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *formSuite) Test_RequiredNestedStructFieldNotSpecified(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&id=1`, formContentType)
	errs := Form(&blogPost, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 1)

	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1})
}

func (s *formSuite) Test_MultipleValuesIntoSlice(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&id=1&name=Matt+Holt&rating=4&rating=3&rating=5`, formContentType)
	errs := Form(&blogPost, req)

	c.Assert(errs, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}, Ratings: []int{4, 3, 5}})
}

func (s *formSuite) Test_UnexportedField(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&id=1&name=Matt+Holt&unexported=foo`, formContentType)
	errs := Form(&blogPost, req)

	c.Assert(errs, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *formSuite) Test_QueryString(c *C) {
	post := Post{}
	req := newRequest(`POST`, `?title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``, formContentType)
	errs := Form(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_QueryStringWithoutContentTypeGET(c *C) {
	post := Post{}
	req := newRequest(`GET`, `?title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``, formContentType)
	errs := Form(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_EmbedStructPointer(c *C) {
	embedPerson := EmbedPerson{}
	req := newRequest(`GET`, `?name=Glorious+Post+Title&email=Lorem+ipsum+dolor+sit+amet`, ``, formContentType)
	errs := Form(&embedPerson, req)

	c.Assert(errs, IsNil)
	c.Assert(embedPerson, DeepEquals, EmbedPerson{&Person{Name: "Glorious Post Title", Email: "Lorem ipsum dolor sit amet"}})
}

func (s *formSuite) Test_EmbedStructPointerRemainNilIfNotBinded(c *C) {
	embedPerson := EmbedPerson{}
	req := newRequest(`GET`, `?`, ``, formContentType)
	errs := Form(&embedPerson, req)

	c.Assert(errs, IsNil)
	c.Assert(embedPerson, DeepEquals, EmbedPerson{nil})
}

func newRequest(method, query, body, contentType string) *http.Request {
	req, err := http.NewRequest(method, query, strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", contentType)
	return req
}
