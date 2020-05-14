package render

func (this *ResultFixture) TestValidationResult() {
	result := ValidationResult{
		Failure1: SimpleInputError("message1", "field1"),
		Failure2: SimpleInputError("message2", "field2"),
		Failure3: nil,
		Failure4: CompoundInputError("message3", "field3", "field4"),
	}

	this.render(result)

	this.assertStatusCode(422)
	this.assertContent(`[{"fields":["field1"],"message":"message1"},{"fields":["field2"],"message":"message2"},{"fields":["field3","field4"],"message":"message3"}]`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestValidationResult_SerializationFailure_HTTP500WithErrorMessage() {
	result := ValidationResult{
		Failure1: new(BadJSON),
	}

	this.render(result)

	this.assertStatusCode(500)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
	this.assertContent(`[{"fields":["HTTP Response"],"message":"Marshal failure"}]`)
}
