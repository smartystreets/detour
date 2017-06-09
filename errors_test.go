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

	problems *Errors
}

func (this *ErrorFixture) Setup() {
	this.problems = new(Errors)
}

func (this *ErrorFixture) TestErrorAppendIf() {
	this.problems.AppendIf(errors.New(""), true)
	this.problems.AppendIf(errors.New(""), false)
	this.problems.AppendIf(errors.New(""), true)
	this.problems.AppendIf(errors.New(""), false)
	this.So(len(this.problems.errors), should.Equal, 2)
}

func (this *ErrorFixture) TestErrorAggregation() {
	this.So(this.problems.errors, should.BeNil)

	this.problems.Append(errors.New("hi"))
	this.problems.Append(errors.New("bye"))

	this.So(len(this.problems.errors), should.Equal, 2)
}

func (this *ErrorFixture) TestErrorSerialization() {
	this.problems.Append(SimpleInputError("Hello", "World"))

	this.So(this.problems.Error(), should.Equal, `[{"fields":["World"],"message":"Hello"}]`)
}
