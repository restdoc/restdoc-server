package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	//	"github.com/smartystreets/goconvey/convey"
	"github.com/gin-gonic/gin"

	"restdoc/config"
)

func TestGetHashedPassword(t *testing.T) {

	cases := []struct {
		password string
		uniq     int64
		want     string
	}{
		{"helloworld", 342872526663266958, "55c70c69d715b92fbddcb7f8cdab9a6f26c74216e043369d8757ed74ca6c634b"},
		{"hhlwujaetkyhvvis", 342897644386910212, "da7eb5bca185db27c58e6a22a8daf47041ba882d184345b47241b00e2eb300cd"},
	}

	for _, c := range cases {
		hashed := GetHashedPassword(c.uniq, c.password)
		if hashed != c.want {
			t.Errorf("GetHashedPassword(%d, %q) == %q, want(%q)", c.uniq, c.password, hashed, c.want)
		}
	}

}

func TestExtractSaaSInfo(t *testing.T) {
	w := httptest.NewRecorder()

	//c.Request, _ = http.NewRequest("GET", "https://www.hedwi.com/signup", nil)

	defaultSaaSDomain := "hedwi.com"

	cases := []struct {
		url        string
		saasDomain string
		want       string
	}{
		{"https://www.hedwi.com", defaultSaaSDomain, ""},
		{"http://www.hedwi.com", defaultSaaSDomain, ""},
		{"https://hedwi.com", defaultSaaSDomain, ""},
		{"http://hedwi.com", defaultSaaSDomain, ""},
		{"https://www.feetmail.com", defaultSaaSDomain, "hedwi.com"},
		{"http://www.feetmail.com", defaultSaaSDomain, ""},
		{"https://feetmail.com", defaultSaaSDomain, ""},
		{"http://feetmail.com", defaultSaaSDomain, ""},
		{"https://www.hedwi.com", "", ""},
		{"http://www.hedwi.com", "", ""},
		{"https://hedwi.com", "", ""},
		{"http://hedwi.com", "", ""},
		{"https://www.feetmail.com", "", "hedwi.com"},
		{"http://www.feetmail.com", "", ""},
		{"https://feetmail.com", "", ""},
		{"http://feetmail.com", "", ""},
	}

	for _, c := range cases {
		config.DefaultConfig.SaaSDomain = c.saasDomain
		req, _ := gin.CreateTestContext(w)
		req.Request, _ = http.NewRequest("GET", c.url, nil)
		result := ExtractSaaSInfo(req)
		t.Errorf("result %v", result)

	}
}

func TestFormatColor(t *testing.T) {

	cases := []struct {
		color int32
		want  string
	}{
		{0, "#427fed"},
		{1, "#000001"},
		{4358125, "#427fed"},
	}

	for _, c := range cases {
		result := FormatColor(c.color)
		if result != c.want {
			t.Errorf("result %v", result)
		}
	}
}
