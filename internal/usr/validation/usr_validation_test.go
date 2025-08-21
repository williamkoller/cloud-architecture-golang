package validation

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type errResp struct {
	Error string `json:"error"`
}

func TestRespondValidationError_Returns400AndJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	RespondValidationError(c, errors.New("bad input"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status code: got %d, want %d", w.Code, http.StatusBadRequest)
	}

	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content-type: got %q, want prefix %q", ct, "application/json")
	}

	var body errResp
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Error != "bad input" {
		t.Fatalf(`body.error: got %q, want %q`, body.Error, "bad input")
	}
}

func TestRespondValidationError_WithNilError_Panics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when err == nil, but did not panic")
		}
	}()

	RespondValidationError(c, nil)
}
