package binding

import . "gopkg.in/check.v1"

type jsonSuite struct{}

var _ = Suite(&jsonSuite{})

func (s *jsonSuite) Test_HappyPath(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title": "Glorious Post Title", "content": "Lorem ipsum dolor sit amet"}`, jsonContentType)
	err := JSON.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *jsonSuite) Test_NilPayload(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `-nil-`, jsonContentType)
	err := JSON.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{})
}

func (s *jsonSuite) Test_EmptyPayload(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, ``, jsonContentType)
	err := JSON.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{})
}

func (s *jsonSuite) Test_EmptyContentType(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title": "Glorious Post Title", "content": "Lorem ipsum dolor sit amet"}`, ``)
	err := JSON.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *jsonSuite) Test_UnsupportedContentType(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title": "Glorious Post Title", "content": "Lorem ipsum dolor sit amet"}`, "BoGus")
	err := JSON.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *jsonSuite) Test_MalformedJson(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `{"title":"foo"`, jsonContentType)
	err := JSON.Bind(&post, req)

	c.Assert(err, NotNil)
	c.Assert(err, DeepEquals, ErrorDeserialization)
	c.Assert(post, DeepEquals, Post{})
}

func (s *jsonSuite) Test_DeserializationWithNestedAndEmbeddedStruct(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `{"title":"Glorious Post Title", "id":1, "author":{"name":"Matt Holt"}}`, jsonContentType)
	err := JSON.Bind(&blogPost, req)

	c.Assert(err, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *jsonSuite) Test_RequiredNestedStructFieldNotSpecified(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `{"title":"Glorious Post Title", "id":1, "author":{}}`, jsonContentType)
	err := JSON.Bind(&blogPost, req)

	c.Assert(err, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1})
}

func (s *jsonSuite) Test_RequiredEmbeddedStructFieldNotSpecified(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `{"id":1, "author":{"name":"Matt Holt"}}`, jsonContentType)
	err := JSON.Bind(&blogPost, req)

	c.Assert(err, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *jsonSuite) Test_SliceOfPosts(c *C) {
	posts := []Post{}
	req := newRequest(`POST`, ``, `[{"title": "First Post"}, {"title": "Second Post"}]`, jsonContentType)
	err := JSON.Bind(&posts, req)

	c.Assert(err, IsNil)
	c.Assert(posts, DeepEquals, []Post{Post{Title: "First Post"}, Post{Title: "Second Post"}})
}
