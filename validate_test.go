package binding

/*
import . "gopkg.in/check.v1"

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
	})
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
	})

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
	})
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
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Author.Name"},
			Classification: RequiredError,
			Message:        "Required",
		},
	})
}

func (s *validateSuite) Test_StructureSliceRequired(c *C) {
	errs := validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
		Readers:      []Person{Person{Name: ""}, Person{Name: "Person b", Email: "b@test.com"}},
		Contributors: []*Person{&Person{Name: "Michael Boke", Email: "mb@test.com"}, &Person{Name: ""}},
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Readers.0.Name"},
			Classification: RequiredError,
			Message:        "Required",
		},
		Error{
			FieldNames:     []string{"Contributors.1.Name"},
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
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Coauthor.Name"},
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
	})
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
	})
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
	})
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
	})
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
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"0.AlphaDash"},
			Classification: AlphaDashError,
			Message:        "AlphaDash",
		},
		Error{
			FieldNames:     []string{"0.AlphaDashDot"},
			Classification: AlphaDashDotError,
			Message:        "AlphaDashDot",
		},
		Error{
			FieldNames:     []string{"0.MinSize"},
			Classification: MinSizeError,
			Message:        "MinSize",
		},
		Error{
			FieldNames:     []string{"0.MinSizeSlice"},
			Classification: MinSizeError,
			Message:        "MinSize",
		},
		Error{
			FieldNames:     []string{"0.MaxSize"},
			Classification: MaxSizeError,
			Message:        "MaxSize",
		},
		Error{
			FieldNames:     []string{"0.MaxSizeSlice"},
			Classification: MaxSizeError,
			Message:        "MaxSize",
		},
		Error{
			FieldNames:     []string{"0.Email"},
			Classification: EmailError,
			Message:        "Email",
		},
		Error{
			FieldNames:     []string{"0.Url"},
			Classification: UrlError,
			Message:        "Url",
		},
		Error{
			FieldNames:     []string{"0.Range"},
			Classification: RangeError,
			Message:        "Range",
		},
		Error{
			FieldNames:     []string{"0.In"},
			Classification: DefaultError,
			Message:        "Default",
		},
		Error{
			FieldNames:     []string{"0.InInvalid"},
			Classification: InError,
			Message:        "In",
		},
		Error{
			FieldNames:     []string{"0.NotIn"},
			Classification: NotInError,
			Message:        "NotIn",
		},
		Error{
			FieldNames:     []string{"0.Include"},
			Classification: IncludeError,
			Message:        "Include",
		},
		Error{
			FieldNames:     []string{"0.Exclude"},
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
	})
	c.Assert(errs, IsNil)
}
*/
