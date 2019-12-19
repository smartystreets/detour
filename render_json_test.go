package detour

import "net/http"

func (this *ResultFixture) TestJSONResult() {
	result := JSONResult{
		StatusCode: 123,
		Content:    map[string]string{"key": "value"},
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`{"key":"value"}`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestJSONResultIndented() {
	result := JSONResult{
		StatusCode: 123,
		Content:    map[string]string{"key": "value"},
		Indent:     "  ",
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`{
  "key": "value"
}`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestJSONResult_WithCustomContentType() {
	result := JSONResult{
		StatusCode:  123,
		ContentType: "application/custom-json",
		Content:     map[string]string{"key": "value"},
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`{"key":"value"}`)
	this.assertHasHeader(contentTypeHeader, "application/custom-json")
}
func (this *ResultFixture) TestJSONResult_SerializationFailure_HTTP500WithErrorMessage() {
	result := JSONResult{
		StatusCode: 123,
		Content:    new(BadJSON),
	}
	this.render(result)

	this.assertStatusCode(500)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
	this.assertContent(`[{"fields":["HTTP Response"],"message":"Marshal failure"}]`)
}
func (this *ResultFixture) TestJSONResult_StatusCodeDefaultsTo200() {
	result := JSONResult{
		StatusCode: 0,
		Content:    42,
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
}
func (this *ResultFixture) TestJSONResultHeadersCopiedToResponse() {
	this.response.Header().Set("Key", "already-added")
	header := http.Header{"Key": []string{"value"}}
	result := JSONResult{Header: header}

	this.render(result)

	this.assertHasHeader("Key", "value")
	this.assertHasHeader("Key", "already-added")
}
