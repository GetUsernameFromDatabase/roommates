package main

import (
	"io/fs"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/handlers"
	"roommates/internal/mailer"
	"roommates/internal/middleware"
	"roommates/internal/services"
	"roommates/internal/sessions"
	"roommates/web/assets"
	"roommates/web/renderer"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initRoutes(r *gin.Engine, h *handlers.DefaultHandler, service *services.DefaultService) {
	mwI18n := middleware.NewInternationalization()
	// separate counters per kind, flatmates on one Wi-Fi share an IP
	// each attempt costs a bcrypt hash (~250ms): login, register, new password
	mwAuthLimit := middleware.NewRateLimit(10, time.Minute)
	// each sends an email, so mailboxes cannot be flooded
	mwEmailLimit := middleware.NewRateLimit(10, time.Minute)
	// invite codes are short enough to be worth guessing
	mwInviteLimit := middleware.NewRateLimit(20, time.Minute)
	r.Use(middleware.NewSecurityHeader())
	r.Use(middleware.NewCrossOriginProtection())
	// before every page, error pages too, so all of them render in the chosen theme
	r.Use(middleware.NewAppearance())

	registerAssets(r)
	public := r.Group("")
	{
		public.Use(middleware.NewGuestOnly())
		public.Use(mwI18n)

		public.GET(g.RouteLogin, h.PageLogin)
		public.POST(g.RouteLogin, mwAuthLimit, h.PageLogin)

		public.GET(g.RouteRegister, h.PageRegister)
		public.POST(g.RouteRegister, mwAuthLimit, h.PageRegister)

		public.GET(g.RouteForgotPassword, h.PageForgotPassword)
		public.POST(g.RouteForgotPassword, mwEmailLimit, h.PageForgotPassword)
	}

	// --- --- ---

	// email links work whether signed in or not
	// their tokens have 130 random bits, guessing needs no rate limit
	open := r.Group("")
	{
		open.Use(mwI18n)

		open.GET(g.RouteVerifyEmail, h.VerifyEmail)
		open.GET(g.RouteResetPassword, h.PageResetPassword)
		open.POST(g.RouteResetPassword, mwAuthLimit, h.PageResetPassword)
		open.POST(g.RouteAppearance, h.SetAppearance)
	}

	// --- --- ---

	// endpoints that need authentication
	// authorization should be checked by the handler
	priv := r.Group("")
	{
		priv.Use(middleware.NewOnlyAuthenticated())
		priv.Use(mwI18n)

		priv.GET("/", h.PageHome)
		priv.POST(g.RouteHouseholds, h.PageHome)
		priv.POST(g.RouteLogout, h.Logout)
		priv.POST(g.RouteVerifyEmailResend, mwEmailLimit, h.ResendVerificationEmail)

		priv.GET(g.RouteInvite, mwInviteLimit, h.PageInvite)
		priv.POST(g.RouteInvite, mwInviteLimit, h.PageInvite)

		// everything under it is scoped to one household the user is a member of
		household := priv.Group(g.RouteHousehold + "/:" + g.ParamHouseholdID)
		household.Use(middleware.NewHouseholdMember(service))
		{
			household.GET("", h.PageHousehold)
			household.GET("/people", h.PagePeople)
			household.POST("/invites", h.CreateInvite)
			household.POST("/invites/email", mwEmailLimit, h.SendInviteEmail)
			household.POST("/invites/:"+g.ParamInviteID+"/delete", h.DeleteInvite)
			household.POST("/members/:"+g.ParamUserID+"/remove", h.RemoveMember)
			household.POST("/members/:"+g.ParamUserID+"/promote", h.PromoteMember)
			household.POST("/leave", h.LeaveHousehold)

			household.POST("/lists", h.CreateShoppingList)
			list := household.Group("/lists/:" + g.ParamListID)
			list.GET("", h.PageShoppingList)
			list.POST("/delete", h.DeleteShoppingList)
			list.POST("/items", h.AddShoppingItem)
			list.POST("/clear", h.ClearCheckedShoppingItems)
			list.POST("/items/:"+g.ParamItemID+"/claim", h.SetShoppingItemClaim)
			list.POST("/items/:"+g.ParamItemID+"/check", h.SetShoppingItemChecked)
			list.POST("/items/:"+g.ParamItemID+"/delete", h.DeleteShoppingItem)

			household.GET("/notes", h.PageNotes)
			household.GET("/notes/new", h.PageNewNote)
			household.POST("/notes", h.CreateNote)
			note := household.Group("/notes/:" + g.ParamNoteID)
			note.GET("", h.PageNote)
			note.GET("/edit", h.PageEditNote)
			note.POST("/edit", h.UpdateNote)
			note.POST("/pin", h.SetNotePinned)
			note.POST("/delete", h.DeleteNote)
		}
	}
}

// uses [assets.PublicAssetsFS]
//
//	"/assets/public" is served at "/assets"
//	"/favicon/favicon.ico" is served at "/favicon.ico"
func registerAssets(r *gin.Engine) {
	// TODO: maybe change this at some point
	// https://templ.guide/hosting-and-deployment/hosting-using-docker#adding-static-content-to-the-docker-container
	assetFS := assets.PublicAssetsFS
	publicFS, err := fs.Sub(assetFS, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("could not get public assets")
	}

	r.StaticFS("/assets", http.FS(publicFS))
	r.StaticFileFS("/favicon.ico", "favicon/favicon.ico", http.FS(assetFS))
}

// -----------------------------------------------------------------------------

func InitGinEngine(dbpool *pgxpool.Pool) *gin.Engine {
	var r *gin.Engine
	if g.IsDevelopment() {
		r = gin.Default()
	} else {
		r = gin.New()
		r.Use(gin.Recovery())          // gin.Default() already has this
		r.Use(middleware.NewZerolog()) // gin.Default() has logger
	}

	// https://gin-gonic.com/en/docs/faq/#why-am-i-getting-404-not-found-instead-of-405-method-not-allowed
	r.HandleMethodNotAllowed = true
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error with SetTrustedProxies")
	}
	renderer.ReplaceGinRenderer(r)

	// TODOS:
	// - https://gin-gonic.com/en/docs/logging/write-log/#log-rotation-in-production
	// - https://templ.guide/integrations/web-frameworks/#csrf-protection

	r.Use(sessions.NewPostgreSessionStoreHandler(dbpool))
	// panics now if APP_BASE_URL is missing in release, not on the first email
	_ = g.GetAppBaseURL()
	service := services.NewDefaultService(dbpool, mailer.New())
	handler := handlers.NewDefaultHandler(service)
	initRoutes(r, handler, service)
	return r
}
