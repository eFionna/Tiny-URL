package app

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/eFionna/Tiny-URL/internal/config"
	"github.com/stretchr/testify/assert"
)

func setupApp(t *testing.T) (*App, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		BaseURL:   "http://localhost:8080",
		RedisAddr: mr.Addr(),
	}
	tmpl := template.Must(template.New("index").Parse("<html>{{.ShortURL}}{{.Error}}</html>"))
	appInstance, err := NewApp(cfg, tmpl)
	if err != nil {
		t.Fatal(err)
	}

	return appInstance, func() { mr.Close() }
}

func TestHandleIndexGET(t *testing.T) {
	appInstance, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	appInstance.HandleIndex(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandleIndexPOSTSuccess(t *testing.T) {
	appInstance, cleanup := setupApp(t)
	defer cleanup()

	form := strings.NewReader("url=https://example.com")
	req := httptest.NewRequest(http.MethodPost, "/", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	appInstance.HandleIndex(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandleIndexPOSTInvalidForm(t *testing.T) {
	appInstance, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	appInstance.HandleIndex(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandleRedirect(t *testing.T) {
	appInstance, cleanup := setupApp(t)
	defer cleanup()

	code, err := appInstance.RDB.SetNX(appInstance.Ctx, "short:testcode", "https://example.com", 0).Result()
	assert.True(t, code)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/s/testcode", nil)
	w := httptest.NewRecorder()

	appInstance.HandleRedirect(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
	assert.Equal(t, "https://example.com", resp.Header.Get("Location"))
}

func TestHandleRedirectNotFound(t *testing.T) {
	appInstance, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/s/unknown", nil)
	w := httptest.NewRecorder()

	appInstance.HandleRedirect(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
