package render

func (this *ResultFixture) TestStatusCodeResult() {
	result := StatusCodeResult{
		StatusCode: 456,
		Message:    "Status 456",
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Status 456")
	this.assertHasHeader(contentTypeHeader, plaintextContentType)
}
func (this *ResultFixture) TestStatusCodeResult_StatusCodeDefaultsTo200() {
	result := StatusCodeResult{
		StatusCode: 0,
		Message:    "Status OK",
	}

	this.render(result)

	this.assertStatusCode(200)
	this.assertContent("Status OK")
}
