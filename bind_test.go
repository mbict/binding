package binder

import . "gopkg.in/check.v1"

type bindSuite struct{}

var _ = Suite(&bindSuite{})

func (s *bindSuite) Test_UnkownContenttype(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, `BoGuS`)
	errs := Bind(&post, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs, DeepEquals, Errors{Error{FieldNames: []string{}, Classification: "ContentTypeError", Message: "Unsupported Content-Type"}})
	c.Assert(post, DeepEquals, Post{})
}

func (s *bindSuite) Test_EmptyContentType(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``)
	errs := Bind(&post, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs, DeepEquals, Errors{Error{FieldNames: []string{}, Classification: "ContentTypeError", Message: "Empty Content-Type"}})
	c.Assert(post, DeepEquals, Post{})
}

func (s *bindSuite) Test_Form(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, formContentType)
	errs := Bind(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *bindSuite) Test_FormGET(c *C) {
	post := Post{}
	req := newRequest(`GET`, `?title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``, formContentType)
	errs := Bind(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *bindSuite) Test_FormGETWithoutContentType(c *C) {
	post := Post{}
	req := newRequest(`GET`, `?title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``, ``)
	errs := Bind(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *bindSuite) Test_Multipart(c *C) {
	blogPost := BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}}
	b, w := makeMultipartPayload(blogPost)
	req := newMultipartRequest(b, w.FormDataContentType())
	w.Close()
	response := BlogPost{}
	errs := Bind(&response, req)

	c.Assert(errs, IsNil)
	c.Assert(response, DeepEquals, blogPost)
}

func (s *bindSuite) Test_Json(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title": "Glorious Post Title", "content": "Lorem ipsum dolor sit amet"}`, jsonContentType)
	errs := Bind(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}
