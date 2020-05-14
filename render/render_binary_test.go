package render

import "net/http"

func (this *ResultFixture) TestBinaryResult() {
	result := BinaryResult{
		StatusCode: 456,
		Content:    []byte("Hello, World!"),
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader(contentTypeHeader, octetStreamContentType)
}
func (this *ResultFixture) TestBinaryResult_WithCustomContentType() {
	result := BinaryResult{
		StatusCode:  456,
		ContentType: "application/custom-binary",
		Content:     []byte("Hello, World!"),
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader(contentTypeHeader, "application/custom-binary")
}
func (this *ResultFixture) TestBinaryResult_StatusCodeDefaultsTo200() {
	result := BinaryResult{
		StatusCode: 0,
		Content:    []byte("Status OK"),
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
}
