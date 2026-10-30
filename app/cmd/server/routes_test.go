package main

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"path"
	"roommates/web/assets"
	"testing"

	"github.com/gin-gonic/gin"
)

// only the hashed bundle may be cached for a year, other assets keep their names across deploys
// and a bundle missing now can exist after the next one
func TestAssetCaching(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerAssets(r)
	bundles, err := fs.Glob(assets.PublicAssetsFS, "public/js/*.js")
	if err != nil || len(bundles) == 0 {
		t.Fatalf("no bundle: %v", err)
	}
	for _, tc := range []struct {
		url    string
		status int
		cached bool
	}{
		{"/assets/js/" + path.Base(bundles[0]), http.StatusOK, true},
		{"/assets/js/missing.js", http.StatusNotFound, false},
		{"/assets/date-picker.js", http.StatusOK, false},
	} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.url, nil))
		if w.Code != tc.status || (w.Header().Get("Cache-Control") != "") != tc.cached {
			t.Errorf("%s: status %d, Cache-Control %q", tc.url, w.Code, w.Header().Get("Cache-Control"))
		}
	}
}
