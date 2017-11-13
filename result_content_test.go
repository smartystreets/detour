package detour

import (
	"net/http"

	"github.com/smartystreets/assertions/should"
)

func (this *ResultFixture) TestContentResult() {
	result := ContentResult{
		StatusCode: 456,
		Content:    "Hello, World!",
		Headers:    map[string]string{"custom-header": "custom-value"},
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader(contentTypeHeader, plaintextContentType)
	this.assertHasHeader("Custom-Header", "custom-value")
	this.So(len(this.response.HeaderMap), should.Equal, 2)
}
func (this *ResultFixture) TestContentResult_WithCustomContentType() {
	result := ContentResult{
		StatusCode:  456,
		ContentType: "application/custom-text",
		Content:     "Hello, World!",
		Headers:     map[string]string{"another-header": "another-value"},
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader(contentTypeHeader, "application/custom-text")
	this.assertHasHeader("Another-Header", "another-value")
	this.So(len(this.response.HeaderMap), should.Equal, 2)
}
func (this *ResultFixture) TestContentResult_WithCustomContentTypeListedTwice() {
	result := ContentResult{
		StatusCode:  456,
		ContentType: "application/custom-text",
		Content:     "Hello, World!",
		Headers:     map[string]string{"Content-Type": "ignored-content-type"},
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader(contentTypeHeader, "application/custom-text")
	this.So(len(this.response.HeaderMap), should.Equal, 1)
}
func (this *ResultFixture) TestContentResult_StatusCodeDefaultsTo200() {
	result := ContentResult{
		StatusCode: 0,
		Content:    "Status OK",
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
}
