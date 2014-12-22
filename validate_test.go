package binding

import (
	"net/http"

	. "gopkg.in/check.v1"
)

type validateSuite struct{}

var _ = Suite(&validateSuite{})

func (s *validateSuite) Test_NoErrors(c *C) {
	errs := validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
	}, dummyRequest())
	c.Assert(errs, IsNil)
}

func (s *validateSuite) Test_IdRequired(c *C) {
	errs := validate(BlogPost{
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
	}, dummyRequest())

	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{Error{FieldNames: []string{"Id"}, Classification: RequiredError, Message: "Required"}})
}

func (s *validateSuite) Test_EmbeddedStructFieldRequired(c *C) {
	errs := validate(BlogPost{
		Id: 1,
		Post: Post{
			Content: "Content given, but title is required",
		},
		Author: Person{
			Name: "Matt Holt",
		},
	}, dummyRequest())
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Title"},
			Classification: RequiredError,
			Message:        "Required",
		},
		Error{
			FieldNames:     []string{"Title"},
			Classification: "LengthError",
			Message:        "Life is too short",
		},
	})
}

func (s *validateSuite) Test_NestedStructFieldRequired(c *C) {
	errs := validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
	}, dummyRequest())
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Name"},
			Classification: RequiredError,
			Message:        "Required",
		},
	})
}

func (s *validateSuite) Test_RequiredFieldMissingInNestedStructPointer(c *C) {
	errs := validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
		Coauthor: &Person{},
	}, dummyRequest())
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Name"},
			Classification: RequiredError,
			Message:        "Required",
		},
	})
}

func (s *validateSuite) Test_AllRequiredFieldsSpecifiedInNestedStructPointer(c *C) {
	errs := validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
		Coauthor: &Person{
			Name: "Jeremy Saenz",
		},
	}, dummyRequest())
	c.Assert(errs, IsNil)
}

func (s *validateSuite) Test_CustomStructValidation(c *C) {
	errs := validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Too short",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
	}, dummyRequest())
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Title"},
			Classification: "LengthError",
			Message:        "Life is too short",
		},
	})
}

func (s *validateSuite) Test_ListValidation(c *C) {
	errs := validate([]BlogPost{
		BlogPost{
			Id: 1,
			Post: Post{
				Title:   "First Post",
				Content: "And some content",
			},
			Author: Person{
				Name: "Leeor Aharon",
			},
		},
		BlogPost{
			Id: 2,
			Post: Post{
				Title:   "Second Post",
				Content: "And some content",
			},
			Author: Person{
				Name: "Leeor Aharon",
			},
		},
	}, dummyRequest())
	c.Assert(errs, IsNil)
}

func (s *validateSuite) Test_ListValidationErrors(c *C) {
	errs := validate([]BlogPost{
		BlogPost{
			Id: 1,
			Post: Post{
				Title:   "First Post",
				Content: "And some content",
			},
			Author: Person{
				Name: "Leeor Aharon",
			},
		},
		BlogPost{
			Id: 2,
			Post: Post{
				Title:   "Too Short",
				Content: "And some content",
			},
			Author: Person{
				Name: "Leeor Aharon",
			},
		},
	}, dummyRequest())
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Title"},
			Classification: "LengthError",
			Message:        "Life is too short",
		},
	})
}

func (s *validateSuite) Test_ListOfInvalidCustomValidations(c *C) {
	errs := validate([]SadForm{
		SadForm{
			AlphaDash:    ",",
			AlphaDashDot: ",",
			MinSize:      ",",
			MinSizeSlice: []string{",", ","},
			MaxSize:      ",,",
			MaxSizeSlice: []string{",", ","},
			Email:        ",",
			Url:          ",",
			UrlEmpty:     "",
			Range:        3,
			InInvalid:    "4",
			NotIn:        "1",
			Include:      "def",
			Exclude:      "abc",
		},
	}, dummyRequest())
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"AlphaDash"},
			Classification: AlphaDashError,
			Message:        "AlphaDash",
		},
		Error{
			FieldNames:     []string{"AlphaDashDot"},
			Classification: AlphaDashDotError,
			Message:        "AlphaDashDot",
		},
		Error{
			FieldNames:     []string{"MinSize"},
			Classification: MinSizeError,
			Message:        "MinSize",
		},
		Error{
			FieldNames:     []string{"MinSizeSlice"},
			Classification: MinSizeError,
			Message:        "MinSize",
		},
		Error{
			FieldNames:     []string{"MaxSize"},
			Classification: MaxSizeError,
			Message:        "MaxSize",
		},
		Error{
			FieldNames:     []string{"MaxSizeSlice"},
			Classification: MaxSizeError,
			Message:        "MaxSize",
		},
		Error{
			FieldNames:     []string{"Email"},
			Classification: EmailError,
			Message:        "Email",
		},
		Error{
			FieldNames:     []string{"Url"},
			Classification: UrlError,
			Message:        "Url",
		},
		Error{
			FieldNames:     []string{"Range"},
			Classification: RangeError,
			Message:        "Range",
		},
		Error{
			FieldNames:     []string{"In"},
			Classification: DefaultError,
			Message:        "Default",
		},
		Error{
			FieldNames:     []string{"InInvalid"},
			Classification: InError,
			Message:        "In",
		},
		Error{
			FieldNames:     []string{"NotIn"},
			Classification: NotInError,
			Message:        "NotIn",
		},
		Error{
			FieldNames:     []string{"Include"},
			Classification: IncludeError,
			Message:        "Include",
		},
		Error{
			FieldNames:     []string{"Exclude"},
			Classification: ExcludeError,
			Message:        "Exclude",
		},
	})
}

func (s *validateSuite) Test_ListOfValidCustomValidations(c *C) {
	errs := validate([]SadForm{
		SadForm{
			AlphaDash:    "123-456",
			AlphaDashDot: "123.456",
			MinSize:      "12345",
			MinSizeSlice: []string{"1", "2", "3", "4", "5"},
			MaxSize:      "1",
			MaxSizeSlice: []string{"1"},
			Email:        "123@456.com",
			Url:          "http://123.456",
			Range:        2,
			In:           "1",
			InInvalid:    "1",
			Include:      "abc",
		},
	}, dummyRequest())
	c.Assert(errs, IsNil)
}

func dummyRequest() *http.Request {
	req, _ := http.NewRequest("GET", "", nil)
	return req
}
