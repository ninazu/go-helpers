package ninazu

import (
	"net/url"
	"net/http"
)

type ReCaptcha struct {
	key string
}

func sendRequest(url string) (*http.Response) {
	resp, err := http.Get(url)

	if err != nil {
		return nil
	}

	return resp
}

func (r *ReCaptcha) buildUrl(script string, v *url.Values) string {
	v.Set("key", r.key)
	v.Set("json", "1")

	return "http://rucaptcha.com/" + script + "?" + v.Encode()
}

func (r *ReCaptcha) ReportBadCaptcha(id string) {
	v := url.Values{}
	v.Set("method", "userrecaptcha")
	v.Set("action", "get")
	v.Set("id", id)

	sendRequest(r.buildUrl("res.php", &v))
}

func (r *ReCaptcha) GetCaptcha(captchaKey, pageUrl string) {
	v := url.Values{}
	v.Set("method", "userrecaptcha")
	v.Set("googlekey", captchaKey)
	v.Set("pageurl", pageUrl)

	sendRequest(r.buildUrl("in.php", &v))
}
