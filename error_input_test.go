package detour

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestInputErrorFixture(t *testing.T) {
	gunit.Run(new(InputErrorFixture), t)
}

type InputErrorFixture struct {
	*gunit.Fixture
}

func (this *InputErrorFixture) TestInputErrorMarshaled() {
	err := &InputError{HTTPStatusCode: 400, Message: "Message", Fields: []string{"Field1"}}
	rendered := err.Error()
	this.So(rendered, should.Equal, `{"fields":["Field1"],"message":"Message"}`)
}

func (this *InputErrorFixture) TestStatusCodeForErrors_DefaultsToZeroIfNotSpecifiedByAnyContainedError() {
	var err Errors
	err = err.Append(&InputError{})
	this.So(err.StatusCode(), should.Equal, 0)
}

func (this *InputErrorFixture) TestStatusCodeForErrors_UsesFirstNonZeroProvidedStatusCodeFromInnerErrors() {
	var err Errors
	err = err.Append(&InputError{HTTPStatusCode: 0})
	err = err.Append(&InputError{HTTPStatusCode: 200})
	err = err.Append(&InputError{HTTPStatusCode: 201})
	err = err.Append(&InputError{HTTPStatusCode: 202})
	this.So(err.StatusCode(), should.Equal, 200)
}
