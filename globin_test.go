package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/podre-Henrique/globin/ts"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	app := setupApp()

	t.Run("success", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"original": "https://google.com",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result URL
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Len(t, result.ShortURL, URL_LEN)
		assert.Empty(t, result.Original)
	})

	t.Run("invalid_url", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"original": "not-a-valid-url",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("bad_request_body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"original":`)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestRedirectURL(t *testing.T) {
	app := setupApp()

	originalURL := "https://example.com/"
	shortURL := "abcdef"
	db.URLs[shortURL] = URLDB{Original: originalURL, start: ts.Timestamp()}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
		assert.Equal(t, originalURL, resp.Header.Get("Location"))
	})

	t.Run("not_found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/123456", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("invalid_length", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/abc", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid_char", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/12345$", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDeleteExpiredURLs(t *testing.T) {
	db = DB{URLs: make(map[string]URLDB)}
	ts.StartTimeStampUpdater()
	time.Sleep(1 * time.Second)

	now := ts.Timestamp()
	expiredStart := now - (URL_LIFETIME + 60)

	db.URLs["expiredKey"] = URLDB{Original: "https://expired.com", start: expiredStart}
	db.URLs["activeKey"] = URLDB{Original: "https://active.com", start: now}

	deleteExpiredURLs()

	db.RLock()
	defer db.RUnlock()
	_, existsExpired := db.URLs["expiredKey"]
	_, existsActive := db.URLs["activeKey"]

	assert.False(t, existsExpired)
	assert.True(t, existsActive)
}

func TestValidURL(t *testing.T) {
	testCases := []struct {
		url      string
		expected bool
	}{
		{"https://google.com", true},
		{"http://example.com/path?query=1", true},
		{"ftp://files.net", true},
		{"invalid-url", false},
		{"www.google.com", false},
		{"", false},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, validURL(tc.url), "URL: "+tc.url)
	}
}
