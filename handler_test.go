package health

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_NewHandler(t *testing.T) {
	// How can I test a function that returns a Struct ?
	// A better Idea? Please tell me!
	h := NewHandler()
	handler := &h

	if handler == nil {
		t.Error("&NewHandler() == nil, wants !nil")
	}
}

func Test_Handler_ServeHTTP_Down(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.AddChecker("DownChecker", downTestChecker{})

	h.ServeHTTP(w, r)

	jsonbytes, _ := ioutil.ReadAll(w.Body)
	jsonstring := strings.TrimSpace(string(jsonbytes))

	wants := `{"info":{"DownChecker":{"status":"DOWN"}},"status":"DOWN"}`

	if jsonstring != wants {
		t.Errorf("jsonReturned == %s, wants %s", jsonstring, wants)
	}

	contentType := w.Header().Get("Content-Type")
	wants = "application/json"

	if contentType != wants {
		t.Errorf("type == %s, wants %s", contentType, wants)
	}

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusServiceUnavailable)
	}
}

func Test_Handler_ServeHTTP_Up(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.AddChecker("UpChecker", upTestChecker{})

	h.ServeHTTP(w, r)
	jsonbytes, _ := ioutil.ReadAll(w.Body)
	jsonstring := strings.TrimSpace(string(jsonbytes))

	wants := `{"info":{"UpChecker":{"status":"UP"}},"status":"UP"}`

	if jsonstring != wants {
		t.Errorf("jsonstring == %s, wants %s", jsonstring, wants)
	}

	contentType := w.Header().Get("Content-Type")
	wants = "application/json"

	if contentType != wants {
		t.Errorf("type == %s, wants %s", contentType, wants)
	}

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}
}
