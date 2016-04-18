package binding

import . "gopkg.in/check.v1"

type bindSuite struct{}

var _ = Suite(&bindSuite{})

func (s *bindSuite) Test_UnkownContenttype(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, `BoGuS`)
	err := Bind(&post, req)

	c.Assert(err, NotNil)
	c.Assert(err, DeepEquals, ErrorUnsupportedContentType)
	c.Assert(post, DeepEquals, Post{})
}

func (s *bindSuite) Test_EmptyContentType(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``)
	err := Bind(&post, req)

	c.Assert(err, NotNil)
	c.Assert(err, DeepEquals, ErrorEmptyContentType)
	c.Assert(post, DeepEquals, Post{})
}

func (s *bindSuite) Test_Form(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, formContentType)
	err := Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *bindSuite) Test_FormGET(c *C) {
	post := Post{}
	req := newRequest(`GET`, `?title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``, formContentType)
	err := Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *bindSuite) Test_FormGETWithoutContentType(c *C) {
	post := Post{}
	req := newRequest(`GET`, `?title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``, ``)
	err := Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *bindSuite) Test_Multipart(c *C) {
	blogPost := BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}}
	b, w := makeMultipartPayload(blogPost)
	req := newMultipartRequest(b, w.FormDataContentType())
	w.Close()
	response := BlogPost{}
	err := Bind(&response, req)

	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, blogPost)
}

func (s *bindSuite) Test_Json(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title": "Glorious Post Title", "content": "Lorem ipsum dolor sit amet"}`, jsonContentType)
	err := Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}
