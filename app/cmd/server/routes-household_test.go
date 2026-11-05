package main

import (
	"crypto/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	g "roommates/internal/config"
	"roommates/internal/db"
	"roommates/internal/db/dbqueries"
	"roommates/internal/events"
	"roommates/internal/handlers"
	"roommates/internal/locales"
	"roommates/internal/mailer"
	"roommates/internal/push"
	"roommates/internal/services"
	"roommates/internal/sessions"
	"roommates/web/renderer"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func randomUUID(t *testing.T) string {
	t.Helper()
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatal(err)
	}
	return pgtype.UUID{Bytes: b, Valid: true}.String()
}

// TestHouseholdRoutesHideFromStrangers walks every route under /h/:householdID, so a route
// added later without the membership check fails here
func TestHouseholdRoutesHideFromStrangers(t *testing.T) {
	testDB := os.Getenv("TEST_POSTGRES_DB")
	if testDB == "" {
		t.Skip("TEST_POSTGRES_DB not set")
	}
	t.Setenv(g.EVK_DatabaseName.String(), testDB)
	t.Setenv(g.EVK_SessionHashKey.String(), strings.Repeat("h", 32))
	t.Setenv(g.EVK_SessionBlockKey.String(), strings.Repeat("b", 32))
	pool, err := pgxpool.New(t.Context(), g.GetDatabaseURL())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	db.LockTestDB(t, pool)
	db.MigrateToLatest(pool)

	locales.LoadLocales()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	renderer.ReplaceGinRenderer(r)
	r.Use(sessions.NewPostgreSessionStoreHandler(pool))
	// a 404 never reaches the store, so none is needed
	service := services.NewDefaultService(pool, mailer.LogMailer{}, nil, push.LogPusher{})
	initRoutes(r, handlers.NewDefaultHandler(service, events.NewHub()), service)

	// users get random emails and are left behind, like in the services tests
	const password = "correct horse battery"
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	signIn := func(name string) (email string, userID pgtype.UUID) {
		email = strings.ToLower(rand.Text()) + "@example.com"
		id, err := service.RegisterUser(ginCtx, dbqueries.InsertUserParams{Email: email, Username: name, Password: password, Language: nil})
		if err != nil {
			t.Fatal(err)
		}
		return email, *id
	}
	// login goes through the real route, and returns the session cookies it sets
	login := func(email string) []*http.Cookie {
		form := url.Values{"email": {email}, "password": {password}}
		w := send(r, http.MethodPost, g.RouteLogin, strings.NewReader(form.Encode()), nil)
		if len(w.Result().Cookies()) == 0 {
			t.Fatalf("login of %s set no session cookie, status %d", email, w.Code)
		}
		return w.Result().Cookies()
	}

	memberEmail, memberID := signIn("member")
	strangerEmail, _ := signIn("stranger")
	householdID, err := service.CreateHousehold(t.Context(), memberID, "Flat")
	if err != nil {
		t.Fatal(err)
	}
	household := householdID.String()
	memberCookies, strangerCookies := login(memberEmail), login(strangerEmail)

	if w := send(r, http.MethodGet, g.RouteHousehold+"/"+household, nil, memberCookies); w.Code == http.StatusNotFound {
		t.Fatalf("member got 404 on their household home: the test would pass with everything broken")
	}

	param := regexp.MustCompile(`:\w+`)
	prefix := g.RouteHousehold + "/:" + g.ParamHouseholdID
	covered := 0
	for _, route := range r.Routes() {
		if !strings.HasPrefix(route.Path, prefix) {
			continue
		}
		covered++
		path := param.ReplaceAllStringFunc(route.Path, func(p string) string {
			if p == ":"+g.ParamHouseholdID {
				return household
			}
			return randomUUID(t)
		})
		t.Run(route.Method+" "+route.Path, func(t *testing.T) {
			if w := send(r, route.Method, path, nil, strangerCookies); w.Code != http.StatusNotFound {
				t.Errorf("stranger got %d, want 404", w.Code)
			}
		})
	}
	if covered == 0 {
		t.Fatalf("no route starts with %s", prefix)
	}
}

// send makes the request a same-origin browser would, so CrossOriginProtection lets it through
func send(r *gin.Engine, method, path string, body *strings.Reader, cookies []*http.Cookie) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req = httptest.NewRequest(method, path, http.NoBody)
	}
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
