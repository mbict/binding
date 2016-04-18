package binding

import . "gopkg.in/check.v1"

type formSuite struct{}

var _ = Suite(&formSuite{})

func (s *formSuite) Test_NotByReference(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, ``, formContentType)
	err := Form.Bind(post, req)

	c.Assert(err, NotNil)
	c.Assert(err, DeepEquals, ErrorInputNotByReference)
}

func (s *formSuite) Test_NotAStruct(c *C) {
	test := int(1)
	req := newRequest(`POST`, ``, ``, formContentType)
	err := Form.Bind(&test, req)

	c.Assert(err, NotNil)
	c.Assert(err, DeepEquals, ErrorInputIsNotStructure)
}

func (s *formSuite) Test_HappyPath(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, formContentType)
	err := Form.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_HappyPathWithPointer(c *C) {
	post := &Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, formContentType)
	err := Form.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, &Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_HappyPathWithNullPointer(c *C) {
	post := (*Post)(nil)
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, formContentType)
	err := Form.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, &Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_EmptyBody(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, ``, formContentType)
	err := Form.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{})
}

func (s *formSuite) Test_EmptyContentType(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``)
	err := Form.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{})
}

func (s *formSuite) Test_MalformedBody(c *C) {
	post := Post{}
	req := newRequest(`POST`, ``, `title=%2`, formContentType)
	err := Form.Bind(&post, req)

	c.Assert(err, NotNil)
	c.Assert(err, DeepEquals, ErrorDeserialization)
	c.Assert(post, DeepEquals, Post{})
}

func (s *formSuite) Test_NestedEmbeddedStructs(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&id=1&author.name=Matt+Holt`, formContentType)
	err := Form.Bind(&blogPost, req)

	c.Assert(err, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *formSuite) Test_RequiredEmbeddedStructFieldNotSpecified(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `id=1&author.name=Matt+Holt`, formContentType)
	err := Form.Bind(&blogPost, req)

	c.Assert(err, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *formSuite) Test_RequiredNestedStructFieldNotSpecified(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&id=1`, formContentType)
	err := Form.Bind(&blogPost, req)

	c.Assert(err, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1})
}

func (s *formSuite) Test_MultipleValuesIntoSlice(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&id=1&author.name=Matt+Holt&rating=4&rating=3&rating=5`, formContentType)
	err := Form.Bind(&blogPost, req)

	c.Assert(err, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}, Ratings: []int{4, 3, 5}})
}

func (s *formSuite) Test_UnexportedField(c *C) {
	blogPost := BlogPost{}
	req := newRequest(`POST`, ``, `title=Glorious+Post+Title&id=1&author.name=Matt+Holt&unexported=foo`, formContentType)
	err := Form.Bind(&blogPost, req)

	c.Assert(err, IsNil)
	c.Assert(blogPost, DeepEquals, BlogPost{Post: Post{Title: "Glorious Post Title"}, Id: 1, Author: Person{Name: "Matt Holt"}})
}

func (s *formSuite) Test_QueryString(c *C) {
	post := Post{}
	req := newRequest(`POST`, `?title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``, formContentType)
	err := Form.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_QueryStringWithoutContentTypeGET(c *C) {
	post := Post{}
	req := newRequest(`GET`, `?title=Glorious+Post+Title&content=Lorem+ipsum+dolor+sit+amet`, ``, formContentType)
	err := Form.Bind(&post, req)

	c.Assert(err, IsNil)
	c.Assert(post, DeepEquals, Post{Title: "Glorious Post Title", Content: "Lorem ipsum dolor sit amet"})
}

func (s *formSuite) Test_EmbedStructPointer(c *C) {
	embedPerson := EmbedPerson{}
	req := newRequest(`GET`, `?name=Glorious+Post+Title&email=Lorem+ipsum+dolor+sit+amet`, ``, formContentType)
	err := Form.Bind(&embedPerson, req)

	c.Assert(err, IsNil)
	c.Assert(embedPerson, DeepEquals, EmbedPerson{&Person{Name: "Glorious Post Title", Email: "Lorem ipsum dolor sit amet"}})
}

func (s *formSuite) Test_EmbedStructPointerPtr(c *C) {
	embedPerson := (*EmbedPerson)(nil)
	req := newRequest(`GET`, `?name=Glorious+Post+Title&email=Lorem+ipsum+dolor+sit+amet`, ``, formContentType)
	err := Form.Bind(&embedPerson, req)

	c.Assert(err, IsNil)
	c.Assert(embedPerson, DeepEquals, &EmbedPerson{&Person{Name: "Glorious Post Title", Email: "Lorem ipsum dolor sit amet"}})
}
