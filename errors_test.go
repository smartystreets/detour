package binding

import (
	"errors"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ErrorFixture struct {
	*gunit.Fixture

	problems ValidationErrors
}

func (this *ErrorFixture) TestErrorAggregation() {
	this.So(this.problems, should.BeNil)

	this.problems = this.problems.Append(errors.New("hi"))
	this.problems = this.problems.Append(errors.New("bye"))

	this.So(len(this.problems), should.Equal, 2)
}

func (this *ErrorFixture) TestValidationErrorSerialization() {
	this.problems = this.problems.Append(SimpleValidationError("Hello", "World"))

	this.So(this.problems.Error(), should.Equal, `[{"fields":["World"],"message":"Hello"}]`)
}
