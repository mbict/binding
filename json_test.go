package binding

import . "gopkg.in/check.v1"

type jsonSuite struct{}

var _ = Suite(&jsonSuite{})

func (s *jsonSuite) Test_HappyPath(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title": "Glorious Post Title", "content": "Lorem ipsum dolor sit amet"}`, jsonContentType)
	errs := Json(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *jsonSuite) Test_NilPayload(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `-nil-`, jsonContentType)
	errs := Json(&post, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs, DeepEquals, Errors{
		Error{FieldNames: []string{"Title"}, Classification: "RequiredError", Message: "Required"},
		Error{FieldNames: []string{"Title"}, Classification: "LengthError", Message: "Life is too short"}})
	c.Assert(post, DeepEquals, Post{})
}

func (s *jsonSuite) Test_EmptyPayload(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, ``, jsonContentType)
	errs := Json(&post, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs, DeepEquals, Errors{
		Error{FieldNames: []string{"Title"}, Classification: "RequiredError", Message: "Required"},
		Error{FieldNames: []string{"Title"}, Classification: "LengthError", Message: "Life is too short"}})
	c.Assert(post, DeepEquals, Post{})
}

func (s *jsonSuite) Test_EmptyContentType(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title": "Glorious Post Title", "content": "Lorem ipsum dolor sit amet"}`, ``)
	errs := Json(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *jsonSuite) Test_UnsupportedContentType(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title": "Glorious Post Title", "content": "Lorem ipsum dolor sit amet"}`, "BoGus")
	errs := Json(&post, req)

	c.Assert(errs, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *jsonSuite) Test_MalformedJson(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title":"foo"`, jsonContentType)
	errs := Json(&post, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 3)
	c.Assert(errs, DeepEquals, Errors{
		Error{FieldNames: []string{}, Classification: "DeserializationError", Message: "unexpected EOF"},
		Error{FieldNames: []string{"Title"}, Classification: "RequiredError", Message: "Required"},
		Error{FieldNames: []string{"Title"}, Classification: "LengthError", Message: "Life is too short"}})

	c.Assert(post, DeepEquals, Post{})
}

func (s *jsonSuite) Test_DeserializationWithNestedAndEmbeddedStruct(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `{"title":"Glorious Post Title", "id":1, "author":{"name":"Matt Holt"}}`, jsonContentType)
	errs := Json(&blogPost, req)

	c.Assert(errs, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *jsonSuite) Test_RequiredNestedStructFieldNotSpecified(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `{"title":"Glorious Post Title", "id":1, "author":{}}`, jsonContentType)
	errs := Json(&blogPost, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs, DeepEquals, Errors{Error{FieldNames: []string{"Name"}, Classification: "RequiredError", Message: "Required"}})
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1})
}

func (s *jsonSuite) Test_RequiredEmbeddedStructFieldNotSpecified(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `{"id":1, "author":{"name":"Matt Holt"}}`, jsonContentType)
	errs := Json(&blogPost, req)

	c.Assert(errs, NotNil)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs, DeepEquals, Errors{
		Error{FieldNames: []string{"Title"}, Classification: "RequiredError", Message: "Required"},
		Error{FieldNames: []string{"Title"}, Classification: "LengthError", Message: "Life is too short"}})
	c.Assert(blogPost, DeepEquals, BlogPost{Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *jsonSuite) Test_SliceOfPosts(c *C) {
	posts := []Post{}
	req := newRequest(`POST`, ``, `[{"title": "First Post"}, {"title": "Second Post"}]`, jsonContentType)
	errs := Json(&posts, req)

	c.Assert(errs, IsNil)
	c.Assert(posts, DeepEquals, []Post{Post{Title: "First Post"}, Post{Title: "Second Post"}})
}
