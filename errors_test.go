package detour

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestErrorFixture(t *testing.T) {
	gunit.Run(new(ErrorFixture), t)
}

type ErrorFixture struct {
	*gunit.Fixture

	problems Errors
}

func (this *ErrorFixture) TestErrorAppendIf() {
	this.problems = this.problems.AppendIf(errors.New(""), true)
	this.problems = this.problems.AppendIf(errors.New(""), false)
	this.problems = this.problems.AppendIf(errors.New(""), true)
	this.problems = this.problems.AppendIf(errors.New(""), false)
	this.So(len(this.problems), should.Equal, 2)
}

func (this *ErrorFixture) TestErrorAggregation() {
	this.So(this.problems, should.BeNil)

	this.problems = this.problems.Append(errors.New("hi"))
	this.problems = this.problems.Append(errors.New("bye"))

	this.So(len(this.problems), should.Equal, 2)
}

func (this *ErrorFixture) TestErrorSerialization() {
	this.problems = this.problems.Append(SimpleInputError("Hello", "World"))

	this.So(this.problems.Error(), should.Equal, `[{"fields":["World"],"message":"Hello"}]`)
}

func (this *ErrorFixture) TestInputErrorMarshaled() {
	err := &InputError{HTTPStatusCode: 400, Message: "Message", Fields: []string{"Field1"}}
	rendered := err.Error()
	this.So(rendered, should.Equal, `{"fields":["Field1"],"message":"Message"}`)
}
