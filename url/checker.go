package url

import (
	"net/http"
	"time"

	"github.com/dimiro1/health"
)

// Checker is a checker that check a given URL
type Checker struct {
	URL     string
	Timeout time.Duration
}

// NewChecker returns a new url.Checker with the given URL
func NewChecker(url string) Checker {
	return Checker{URL: url, Timeout: 5 * time.Second}
}

// NewCheckerWithTimeout returns a new url.Checker with the given URL and given timeout
func NewCheckerWithTimeout(url string, timeout time.Duration) Checker {
	return Checker{URL: url, Timeout: timeout}
}

// Check makes a HEAD request to the given URL
// If the request returns something different than StatusOK,
// The status is set to StatusBadRequest and the URL is considered Down
func (u Checker) Check() health.Health {
	client := http.Client{
		Timeout: u.Timeout,
	}

	var (
		resp *http.Response
		err  error
	)

	health := health.NewHealth()

	resp, err = client.Head(u.URL)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		// retry with GET request to the given URL
		resp, err = client.Get(u.URL)
		if err != nil {
			health.Down().AddInfo("code", http.StatusBadRequest)
			return health
		}
	}

	if resp.StatusCode == http.StatusOK {
		health.Up()
	} else {
		health.Down()
	}

	health.AddInfo("code", resp.StatusCode)

	return health
}
