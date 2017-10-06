package ninazu

import (
	"net/url"
	"net/http"
)

type reCaptcha struct {
	key string
}

type response struct {
	Id    string
	Token string
}

func ReCaptcha(key string) (*reCaptcha) {
	return &reCaptcha{
		key: key,
	}
}

func sendRequest(url string) (*http.Response) {
	resp, err := http.Get(url)

	if err != nil {
		return nil
	}

	return resp
}

func (r *reCaptcha) buildUrl(script string, v *url.Values) string {
	v.Set("key", r.key)
	v.Set("json", "1")

	return "http://rucaptcha.com/" + script + "?" + v.Encode()
}

func (r *reCaptcha) ReportBadCaptcha(id string) {
	v := url.Values{}
	v.Set("method", "userrecaptcha")
	v.Set("action", "get")
	v.Set("id", id)

	sendRequest(r.buildUrl("res.php", &v))
}

func (r *reCaptcha) GetCaptcha(captchaKey, pageUrl string) *response {
	v := url.Values{}
	v.Set("method", "userrecaptcha")
	v.Set("googlekey", captchaKey)
	v.Set("pageurl", pageUrl)

	sendRequest(r.buildUrl("in.php", &v))

	return &response{
		Id:    "1",
		Token: "Hello",
	}
}
