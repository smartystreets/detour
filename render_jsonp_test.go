package detour

import "net/http"

func (this *ResultFixture) TestJSONPResult() {
	this.setRequestURLCallback("maybe")
	result := JSONPResult{
		StatusCode: 123,
		Content:    map[string]string{"key": "value"},
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`maybe({"key":"value"})`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestJSONPResultIndented() {
	this.setRequestURLCallback("maybe")
	result := JSONPResult{
		StatusCode: 123,
		Content:    map[string]string{"key": "value"},
		Indent:     "  ",
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`maybe({
  "key": "value"
})`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestJSONPResult_WithCustomContentType() {
	this.setRequestURLCallback("maybe")
	result := JSONPResult{
		StatusCode:  123,
		ContentType: "application/custom-json",
		Content:     map[string]string{"key": "value"},
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`maybe({"key":"value"})`)
	this.assertHasHeader(contentTypeHeader, "application/custom-json")
}
func (this *ResultFixture) TestJSONPResult_SerializationFailure_HTTP500WithErrorMessage() {
	this.setRequestURLCallback("maybe")
	result := JSONPResult{
		StatusCode: 123,
		Content:    new(BadJSON),
	}

	this.render(result)

	this.assertStatusCode(500)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
	this.assertContent(`[{"fields":["HTTP Response"],"message":"Marshal failure"}]`)
}
func (this *ResultFixture) TestJSONPResult_StatusCodeDefaultsTo200() {
	this.setRequestURLCallback("maybe")
	result := JSONPResult{
		StatusCode: 0,
		Content:    42,
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
}
func (this *ResultFixture) TestJSONPResult_NoCallback_SerializesAsPlainOldJSON() {
	this.setRequestURLCallback("")
	result := JSONPResult{
		StatusCode: 123,
		Content:    map[string]string{"key": "value"},
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`{"key":"value"}`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestJSONPResultHeadersAddedToResponse() {
	this.response.Header().Set("Key", "already-added")
	headers := http.Header{"Key": []string{"value"}}
	result := JSONPResult{Header: headers}

	this.render(result)

	this.assertHasHeader("Key", "value")
	this.assertHasHeader("Key", "already-added")
}
