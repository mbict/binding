package binding

import (
	"bytes"
	"mime/multipart"
	"net/http"

	. "gopkg.in/check.v1"
)

type fileSuite struct{}

type fileInfo struct {
	fieldName string
	fileName  string
	data      string
}

var _ = Suite(&fileSuite{})

func (s *fileSuite) Test_SingleFile(c *C) {

	blogPost := BlogPost{}
	req := buildRequestWithFile([]fileInfo{
		fileInfo{
			fieldName: "headerImage",
			fileName:  "message.txt",
			data:      "All your binding are belong to us",
		},
	})
	MultipartForm(&blogPost, req)

	c.Assert(blogPost.HeaderImage, NotNil)
	c.Assert(blogPost.HeaderImage.Filename, Equals, "message.txt")
	c.Assert(unpackFileData(blogPost.HeaderImage), Equals, "All your binding are belong to us")

	c.Assert(blogPost.Pictures, HasLen, 0)
}

func (s *fileSuite) Test_MultipleFiles(c *C) {
	blogPost := BlogPost{}
	req := buildRequestWithFile([]fileInfo{
		fileInfo{
			fieldName: "picture",
			fileName:  "cool-gopher-fact.txt",
			data:      "Did you know? https://plus.google.com/+MatthewHolt/posts/GmVfd6TPJ51",
		},
		fileInfo{
			fieldName: "picture",
			fileName:  "gophercon2014.txt",
			data:      "@bradfitz has a Go time machine: https://twitter.com/mholt6/status/459463953395875840",
		},
	})
	MultipartForm(&blogPost, req)

	c.Assert(blogPost.HeaderImage, IsNil)

	c.Assert(blogPost.Pictures, HasLen, 2)
	c.Assert(blogPost.Pictures[0].Filename, Equals, "cool-gopher-fact.txt")
	c.Assert(unpackFileData(blogPost.Pictures[0]), Equals, "Did you know? https://plus.google.com/+MatthewHolt/posts/GmVfd6TPJ51")
	c.Assert(blogPost.Pictures[1].Filename, Equals, "gophercon2014.txt")
	c.Assert(unpackFileData(blogPost.Pictures[1]), Equals, "@bradfitz has a Go time machine: https://twitter.com/mholt6/status/459463953395875840")
}

func (s *fileSuite) Test_SingleFileAndMultipleFiles(c *C) {
	blogPost := BlogPost{}
	req := buildRequestWithFile([]fileInfo{
		fileInfo{
			fieldName: "headerImage",
			fileName:  "social media.txt",
			data:      "Hey, you should follow @mholt6 (Twitter) or +MatthewHolt (Google+)",
		},
		fileInfo{
			fieldName: "picture",
			fileName:  "thank you!",
			data:      "Also, thanks to all the contributors of this package!",
		},
		fileInfo{
			fieldName: "picture",
			fileName:  "btw...",
			data:      "This tool translates JSON into Go structs: http://mholt.github.io/json-to-go/",
		},
	})
	MultipartForm(&blogPost, req)

	c.Assert(blogPost.HeaderImage, NotNil)
	c.Assert(blogPost.HeaderImage.Filename, Equals, "social media.txt")
	c.Assert(unpackFileData(blogPost.HeaderImage), Equals, "Hey, you should follow @mholt6 (Twitter) or +MatthewHolt (Google+)")

	c.Assert(blogPost.Pictures, HasLen, 2)
	c.Assert(blogPost.Pictures[0].Filename, Equals, "thank you!")
	c.Assert(unpackFileData(blogPost.Pictures[0]), Equals, "Also, thanks to all the contributors of this package!")
	c.Assert(blogPost.Pictures[1].Filename, Equals, "btw...")
	c.Assert(unpackFileData(blogPost.Pictures[1]), Equals, "This tool translates JSON into Go structs: http://mholt.github.io/json-to-go/")
}

func buildRequestWithFile(files []fileInfo) *http.Request {
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	for _, file := range files {
		formFile, err := w.CreateFormFile(file.fieldName, file.fileName)
		if err != nil {
			panic("Could not create FormFile (multiple files): " + err.Error())
		}
		formFile.Write([]byte(file.data))
	}

	err := w.Close()
	if err != nil {
		panic("Could not close multipart writer: " + err.Error())
	}

	req, err := http.NewRequest("POST", "", b)
	if err != nil {
		panic("Could not create file upload request: " + err.Error())
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func unpackFileData(fh *multipart.FileHeader) string {
	if fh == nil {
		return ""
	}

	f, err := fh.Open()
	if err != nil {
		panic("Could not open file header:" + err.Error())
	}
	defer f.Close()

	var fb bytes.Buffer
	_, err = fb.ReadFrom(f)
	if err != nil {
		panic("Could not read from file header:" + err.Error())
	}

	return fb.String()
}
