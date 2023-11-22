package url

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dimiro1/health"
)

// Checker is a checker that check a given URL
type Checker struct {
	URL                string
	Timeout            time.Duration
	Method             string // Defaults to GET
	ExpectBodyContains string // Disabled if empty
	ExpectStatusCode   int    // if 0 defaults to 200
}

// NewChecker returns a new url.Checker with the given URL
func NewChecker(url string) Checker {
	return NewCheckerWithTimeout(url, 5*time.Second)
}

// NewCheckerWithTimeout returns a new url.Checker with the given URL and given timeout
func NewCheckerWithTimeout(url string, timeout time.Duration) Checker {
	return Checker{URL: url, Timeout: timeout}
}

// Check makes a HEAD request to the given URL
// If the request returns something different than StatusOK,
// The status is set to StatusBadRequest and the URL is considered Down
func (u Checker) Check() health.Health {
	ctxWithTimeout, cancelFn := context.WithTimeout(context.Background(), u.Timeout)
	defer cancelFn()
	health := health.NewHealth()
	body := new(bytes.Buffer)
	req, err := http.NewRequestWithContext(ctxWithTimeout, u.Method, u.URL, body)
	if err != nil {
		health.Down().AddInfo("code", http.StatusBadRequest).AddInfo("error", err.Error())
		return health
	}

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		health.Down().AddInfo("code", http.StatusBadRequest).AddInfo("error", err.Error())
		return health
	}

	if u.ExpectStatusCode == 0 {
		u.ExpectStatusCode = http.StatusOK
	}

	if resp.StatusCode == u.ExpectStatusCode {
		health.Up().AddInfo("code", resp.StatusCode)
	} else {
		health.Down().AddInfo("code", resp.StatusCode).AddInfo("expectedCode", u.ExpectStatusCode)
		return health
	}

	if u.ExpectBodyContains != "" {
		respbody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			health.Down().AddInfo("error", err.Error())
			return health
		}
		if !bytes.ContainsAny(respbody, u.ExpectBodyContains) {
			health.Down().AddInfo("body", "does not contain: "+u.ExpectBodyContains)
		}
	}

	return health
}
