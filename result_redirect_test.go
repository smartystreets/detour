package detour

import (
	"net/http"

	"github.com/smartystreets/assertions/should"
)

func (this *ResultFixture) TestRedirectResult() {
	result := RedirectResult{
		Location:   "http://www.google.com",
		StatusCode: http.StatusMovedPermanently,
	}

	this.render(result)

	this.assertStatusCode(http.StatusMovedPermanently)
	this.So(this.response.Header().Get("Location"), should.Equal, "http://www.google.com")
	this.assertContent(`<a href="http://www.google.com">Moved Permanently</a>.`)
}
