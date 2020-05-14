package render

import "net/http"

type CookieResult struct {
	Cookie1 *http.Cookie
	Cookie2 *http.Cookie
	Cookie3 *http.Cookie
	Cookie4 *http.Cookie
}

func (this CookieResult) Render(response http.ResponseWriter, _ *http.Request) {
	for _, cookie := range []*http.Cookie{this.Cookie1, this.Cookie2, this.Cookie3, this.Cookie4} {
		if cookie != nil {
			http.SetCookie(response, cookie)
		}
	}
}
