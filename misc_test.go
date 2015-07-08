package binding

import . "gopkg.in/check.v1"

type Everything struct {
	Integer    int     `form:"integer"`
	Integer8   int8    `form:"integer8"`
	Integer16  int16   `form:"integer16"`
	Integer32  int32   `form:"integer32"`
	Integer64  int64   `form:"integer64"`
	Uinteger   uint    `form:"uinteger"`
	Uinteger8  uint8   `form:"uinteger8"`
	Uinteger16 uint16  `form:"uinteger16"`
	Uinteger32 uint32  `form:"uinteger32"`
	Uinteger64 uint64  `form:"uinteger64"`
	Boolean_1  bool    `form:"boolean_1"`
	Boolean_2  bool    `form:"boolean_2"`
	Fl32_1     float32 `form:"fl32_1"`
	Fl32_2     float32 `form:"fl32_2"`
	Fl64_1     float64 `form:"fl64_1"`
	Fl64_2     float64 `form:"fl64_2"`
	Str        string  `form:"str"`
}

type miscSuite struct{}

var _ = Suite(&miscSuite{})

func (s *miscSuite) Test_AllTypes(c *C) {
	test := Everything{}
	req := newRequest(`POST`, ``, `integer=-1&integer8=-8&integer16=-16&integer32=-32&integer64=-64&uinteger=1&uinteger8=8&uinteger16=16&uinteger32=32&uinteger64=64&boolean_1=true&fl32_1=32.3232&fl64_1=-64.6464646464&str=string`, formContentType)
	errs := Form.Bind(&test, req)

	c.Assert(errs, IsNil)
	c.Assert(test, DeepEquals, Everything{
		Integer:    -1,
		Integer8:   -8,
		Integer16:  -16,
		Integer32:  -32,
		Integer64:  -64,
		Uinteger:   1,
		Uinteger8:  8,
		Uinteger16: 16,
		Uinteger32: 32,
		Uinteger64: 64,
		Boolean_1:  true,
		Fl32_1:     32.3232,
		Fl64_1:     -64.6464646464,
		Str:        "string",
	})
}

func (s *miscSuite) Test_AllTypesError(c *C) {
	test := Everything{}
	req := newRequest(`POST`, ``, `integer=&integer8=asdf&integer16=--&integer32=&integer64=dsf&uinteger=&uinteger8=asdf&uinteger16=+&uinteger32= 32 &uinteger64=+%20+&boolean_1=&boolean_2=asdf&fl32_1=asdf&fl32_2=&fl64_1=&fl64_2=asdfstr`, formContentType)
	errs := Form.Bind(&test, req)

	c.Assert(errs, IsNil)
	c.Assert(test, DeepEquals, Everything{})
}
