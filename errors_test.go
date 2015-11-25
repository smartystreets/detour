package detour

import (
	"errors"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ErrorFixture struct {
	*gunit.Fixture

	problems Errors
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
