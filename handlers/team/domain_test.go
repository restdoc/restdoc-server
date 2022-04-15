package team_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"restdoc/internal/route"
)

func TestNotLogin(t *testing.T) {
	r := route.InitRouter()
	Convey("Get request to /domain should redirect", t, func() {
		req, _ := http.NewRequest("GET", "/domain", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		So(resp.Code, ShouldEqual, 307)
	})
}
