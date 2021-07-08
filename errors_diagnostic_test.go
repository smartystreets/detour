package detour

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDiagnosticErrorsFixture(t *testing.T) {
	gunit.Run(new(DiagnosticErrorsFixture), t)
}

type DiagnosticErrorsFixture struct {
	*gunit.Fixture
	request  *http.Request
	response *httptest.ResponseRecorder
}

func (this *DiagnosticErrorsFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/test", nil)
	this.response = httptest.NewRecorder()
}

func (this *DiagnosticErrorsFixture) TestSingleError_SingleBulletPoint() {
	var err DiagnosticErrors
	err = err.Append(errors.New("A"))
	this.So(err.Error(), should.Equal, "Errors:\n\n- A")
}

func (this *DiagnosticErrorsFixture) TestMultipleErrors_OrderedList() {
	var err DiagnosticErrors
	err = err.Append(errors.New("A"))
	err = err.Append(nil)
	err = err.AppendIf(errors.New("AA"), false)
	err = err.AppendIf(errors.New("B"), true)
	this.So(err.Error(), should.Equal, "Errors:\n\n1. A\n2. B")
}
