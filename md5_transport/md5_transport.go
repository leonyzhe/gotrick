package md5transport

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
)

type Transport struct {
	http.RoundTripper
}

func NewTransport(transport http.RoundTripper) http.RoundTripper {
	return &Transport{transport}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	//do something
	if req.Body == nil {
		return t.RoundTripper.RoundTrip(req)
	}

	b, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return t.RoundTripper.RoundTrip(req)
	}

	m := md5.Sum(b)
	req.Header.Set("X-Md5", hex.EncodeToString(m[:]))

	//由于ioutil。ReadAll方法会读取到EOF，所以需要重置Body
	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return t.RoundTripper.RoundTrip(req)
}
