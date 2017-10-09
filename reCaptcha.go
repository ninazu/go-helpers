package ninazu

/**
* @see https://rucaptcha.com/api-rucaptcha
*/

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
	"errors"
)

type reCaptcha struct {
	key      string
	turnSend int
	turnWait int
	turnGet  int
}

type ruCaptchaResponse struct {
	Status  int
	Request string
}

func ReCaptcha(key string) (*reCaptcha) {
	return &reCaptcha{
		key:      key,
		turnSend: 0,
		turnWait: 0,
		turnGet:  0,
	}
}

func (r *reCaptcha) GetCaptchaSolution(captchaKey, pageUrl string) (string, error) {
	r.turnGet++

	id, err := r.sendJob(captchaKey, pageUrl)

	if err != nil {
		return r.returnTurnGet("", err)
	}

	//Wait first solution
	time.Sleep(20 * time.Second)
	sol, err := r.waitSolution(id)

	if err != nil {

		return r.returnTurnGet("", err)
	} else if len(sol) == 0 {
		return r.GetCaptchaSolution(captchaKey, pageUrl)
	}

	return sol, nil
}

func (r *reCaptcha) trySendJob(url string) (int, error) {
	if r.turnSend > 10 {
		return r.returnTurnSend(0, errors.New("MAX_TURN_SEND"))
	}

	r.turnSend++

	d := sendRequest(url)

	if d.Status == 0 {
		switch d.Request {
		case "MAX_USER_TURN":
			time.Sleep(1 * time.Second)

			return r.trySendJob(url)

		case "ERROR_NO_SLOT_AVAILABLE":
			time.Sleep(5 * time.Second)

			return r.trySendJob(url)

		case "ERROR_ZERO_BALANCE":
			time.Sleep(60 * time.Second)

			return r.trySendJob(url)

		default:
			return r.returnTurnSend(0, errors.New(d.Request))
		}
	}

	i, err := strconv.Atoi(d.Request)

	if err != nil {
		r.returnTurnSend(0, err)
	}

	return r.returnTurnSend(i, nil)
}

func (r *reCaptcha) tryWaitSolution(url string) (string, error) {
	if r.turnWait > 360 {
		return r.returnTurnWait("", errors.New("MAX_TURN_WAIT"))
	}

	r.turnWait++

	d := sendRequest(url)

	if d.Status == 0 {
		switch d.Request {
		case "CAPCHA_NOT_READY":
			time.Sleep(5 * time.Second)

			return r.tryWaitSolution(url)

		case "ERROR_CAPTCHA_UNSOLVABLE":
		case "ERROR_BAD_DUPLICATES":
			return r.returnTurnWait("", nil)

		default:
			return r.returnTurnWait("", errors.New(d.Request))
		}
	}

	return r.returnTurnWait(d.Request, nil)
}

func (r *reCaptcha) buildUrl(script string, v *url.Values) string {
	v.Set("key", r.key)
	v.Set("json", "1")

	return "http://rucaptcha.com/" + script + "?" + v.Encode()
}

func (r *reCaptcha) sendJob(captchaKey, pageUrl string) (int, error) {
	v := url.Values{}
	v.Set("method", "userrecaptcha")
	v.Set("googlekey", captchaKey)
	v.Set("pageurl", pageUrl)

	return r.trySendJob(r.buildUrl("in.php", &v))
}

func (r *reCaptcha) waitSolution(id int) (string, error) {
	v := url.Values{}
	v.Set("method", "userrecaptcha")
	v.Set("action", "get")
	v.Set("id", strconv.Itoa(id))

	return r.tryWaitSolution(r.buildUrl("res.php", &v))
}

func (r *reCaptcha) returnTurnGet(s string, e error) (string, error) {
	r.turnGet = 0

	return s, e
}

func (r *reCaptcha) returnTurnWait(s string, e error) (string, error) {
	r.turnWait = 0

	return s, e
}

func (r *reCaptcha) returnTurnSend(s int, e error) (int, error) {
	r.turnSend = 0

	return s, e
}

func sendRequest(url string) (*ruCaptchaResponse) {
	resp, err := http.Get(url)

	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil
	}

	var data ruCaptchaResponse
	json.Unmarshal(body, &data)

	return &data
}
