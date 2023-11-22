package url

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimiro1/health"
)

func Test_Checker_Check_Up(t *testing.T) {
	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	checker := NewChecker(fmt.Sprintf("%s/up/", server.URL))

	handler := health.NewHandler()
	handler.AddChecker("Up", checker)

	mux.Handle("/health/", handler)
	mux.HandleFunc("/up/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "UP")
	})

	resp, _ := http.Get(fmt.Sprintf("%s/health/", server.URL))

	wants := `{"Up":{"code":200,"status":"UP"},"status":"UP"}`

	check(t, resp, wants, http.StatusOK)
}

func Test_Checker_Check_Down(t *testing.T) {
	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	checker := NewChecker(fmt.Sprintf("%s/down/", server.URL))

	handler := health.NewHandler()
	handler.AddChecker("Down", checker)

	mux.Handle("/health/", handler)
	mux.HandleFunc("/down/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Down")
	})

	resp, _ := http.Get(fmt.Sprintf("%s/health/", server.URL))

	wants := `{"Down":{"code":500,"status":"DOWN"},"status":"DOWN"}`

	check(t, resp, wants, http.StatusServiceUnavailable)
}

func Test_Checker_Check_Down_invalid(t *testing.T) {
	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	checker := NewChecker("")

	handler := health.NewHandler()
	handler.AddChecker("Down", checker)

	mux.Handle("/health/", handler)

	resp, _ := http.Get(fmt.Sprintf("%s/health/", server.URL))

	wants := `{"Down":{"code":400,"status":"DOWN"},"status":"DOWN"}`
	check(t, resp, wants, http.StatusServiceUnavailable)
}

func check(t *testing.T, resp *http.Response, wants string, code int) {
	jsonbytes, _ := ioutil.ReadAll(resp.Body)
	jsonstring := strings.TrimSpace(string(jsonbytes))

	if jsonstring != wants {
		t.Errorf("jsonstring == %s, wants %s", jsonstring, wants)
	}

	contentType := resp.Header.Get("Content-Type")
	wants = "application/json"

	if contentType != wants {
		t.Errorf("type == %s, wants %s", contentType, wants)
	}

	if resp.StatusCode != code {
		t.Errorf("resp.StatusCode == %d, wants %d", resp.StatusCode, code)
	}
}
