package main

import (
	"io/fs"
	"net/http"
	"roommates/assets"
	"roommates/gintemplrenderer"
	g "roommates/globals"
	"roommates/handlers"
	"roommates/middleware"
	"roommates/services"
	"roommates/sessions"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initRoutes(r *gin.Engine, h *handlers.DefaultHandler) {
	mwI18n := middleware.NewInternationalization()
	r.Use(middleware.NewSecurityHeader())

	registerAssets(r)
	public := r.Group("")
	{
		public.Use(middleware.NewGuestOnly())
		public.Use(mwI18n)

		public.GET(g.RouteLogin, h.PageLogin)
		public.POST(g.RouteLogin, h.PageLogin)

		public.GET(g.RouteRegister, h.PageRegister)
		public.POST(g.RouteRegister, h.PageRegister)
	}

	// --- --- ---

	// endpoints that need authentication
	// authorization should be checked by the handler
	priv := r.Group("")
	{
		priv.Use(middleware.NewOnlyAuthenticated())
		priv.Use(mwI18n)

		priv.GET("/")
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

// https://github.com/a-h/templ/blob/main/examples/integration-gin/main.go#L16
func ginTemplIntegration(r *gin.Engine) {
	ginHTMLRenderer := r.HTMLRender // store previous one
	r.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHTMLRenderer: ginHTMLRenderer}
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
	gintemplrenderer.ReplaceGinRenderer(r)
	ginTemplIntegration(r)

	// TODOS:
	// - https://gin-gonic.com/en/docs/logging/write-log/#log-rotation-in-production
	// - https://templ.guide/integrations/web-frameworks/#csrf-protection

	r.Use(sessions.NewPostgreSessionStoreHandler(dbpool))
	service := services.NewDefaultService(dbpool)
	handler := handlers.NewDefaultHandler(service)
	initRoutes(r, handler)
	return r
}
