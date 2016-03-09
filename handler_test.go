package health

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Handler_ServeHTTP_Down(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.AddChecker("DownChecker", downTestChecker{})

	h.ServeHTTP(w, r)

	jsonbytes, _ := ioutil.ReadAll(w.Body)
	jsonstring := strings.TrimSpace(string(jsonbytes))

	wants := `{"status":"down","info":{"DownChecker":{"status":"down","info":"down"}}}`

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

	wants := `{"status":"up","info":{"UpChecker":{"status":"up","info":"up"}}}`

	if jsonstring != wants {
		t.Errorf("jsonReturned == %s, wants %s", jsonstring, wants)
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