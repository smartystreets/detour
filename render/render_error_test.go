package render

import (
	"errors"
	"net/http"
)

func (this *ResultFixture) TestErrorResult() {
	result := ErrorResult{
		StatusCode: 409,
		Error1:     SimpleInputError("message1", "field1"),
		Error2:     SimpleInputError("message2", "field2"),
		Error3:     nil,
		Error4:     CompoundInputError("message3", "field3", "field4"),
	}

	this.render(result)

	this.assertStatusCode(409)
	this.assertContent(`[{"fields":["field1"],"message":"message1"},{"fields":["field2"],"message":"message2"},{"fields":["field3","field4"],"message":"message3"}]`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestErrorResult_StatusCodeDefaultsTo200() {
	result := ErrorResult{
		StatusCode: 0,
		Error1:     errors.New("ok"),
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
}
