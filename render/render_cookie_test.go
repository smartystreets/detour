package render

import (
	"net/http"

	"github.com/smartystreets/assertions/should"
)

func (this *ResultFixture) TestCookieResult() {
	result := CookieResult{
		Cookie1: &http.Cookie{Name: "a", Value: "1"},
		Cookie2: &http.Cookie{Name: "b", Value: "2"},
		Cookie3: nil,
		Cookie4: &http.Cookie{Name: "d", Value: "4"},
	}

	this.render(result)

	this.assertStatusCode(200)
	this.So(this.response.Header()["Set-Cookie"], should.Resemble, []string{"a=1", "b=2", "d=4"})
	this.assertContent("")
}
