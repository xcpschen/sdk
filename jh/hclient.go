package jh

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Hclient struct {
	http.Client
}
type IhttpReq interface {
	// ToHTTPReq return a http.Request pointer and error when is failure
	ToHTTPReq() (req *http.Request, err error)
}

// MaxRetry http request retry max time
const MaxRetry = 2

// DoReq Do http request
func (h *Hclient) DoReq(r IhttpReq) (b []byte, err error) {
	req, err := r.ToHTTPReq()
	if err != nil {
		return
	}
	retry := 0
DoRetry:
	resp, err := h.Do(req)
	if err != nil {
		if retry >= MaxRetry {
			return
		}
		retry++
		time.Sleep(100 * time.Microsecond)
		goto DoRetry
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()
	return
}
