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

// a file is kept for a year only when its URL carries its content hash, as names stay the same
// across deploys; everything else revalidates with its ETag, and a missing file gets no headers
// as it can exist after the next deploy
func TestAssetCaching(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerAssets(r)
	bundles, err := fs.Glob(assets.PublicAssetsFS, "public/js/*.js")
	if err != nil || len(bundles) == 0 {
		t.Fatalf("no bundle: %v", err)
	}
	const immutable, revalidate = "public, max-age=31536000, immutable", "no-cache"
	etag := `"` + assets.Hash("public/date-picker.js") + `"`
	for _, tc := range []struct {
		url, ifNoneMatch string
		status           int
		cacheControl     string
	}{
		{"/assets/js/" + path.Base(bundles[0]), "", http.StatusOK, immutable},
		{"/assets/js/missing.js", "", http.StatusNotFound, ""},
		{assets.URL("date-picker.js"), "", http.StatusOK, immutable},
		{"/assets/date-picker.js", "", http.StatusOK, revalidate},
		{"/assets/date-picker.js?v=stale", "", http.StatusOK, revalidate},
		{"/assets/date-picker.js", etag, http.StatusNotModified, revalidate},
		{"/icons/favicon.svg", "", http.StatusOK, revalidate},
		{"/sw.js", "", http.StatusOK, revalidate},
	} {
		req := httptest.NewRequest(http.MethodGet, tc.url, nil)
		if tc.ifNoneMatch != "" {
			req.Header.Set("If-None-Match", tc.ifNoneMatch)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.status || w.Header().Get("Cache-Control") != tc.cacheControl || (tc.cacheControl != "") != (w.Header().Get("ETag") != "") {
			t.Errorf("%s: status %d, Cache-Control %q, ETag %q", tc.url, w.Code, w.Header().Get("Cache-Control"), w.Header().Get("ETag"))
		}
	}
}
