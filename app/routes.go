package main

import (
	"roommates/controller"
	"roommates/docs"
	"roommates/gintemplrenderer"
	g "roommates/globals"
	"roommates/locales"
	"roommates/middleware"

	"github.com/gin-gonic/gin"
	"github.com/invopop/ctxi18n"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitGinEngine(c *controller.Controller) *gin.Engine {
	e := gin.Default()
	e.SetTrustedProxies(nil)
	// e.Use(middleware.NewCSRFMiddleware())
	ginHtmlRenderer := e.HTMLRender
	e.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	if err := ctxi18n.Load(locales.Content); err != nil {
		log.Error().Err(err).Msg("error loading i18n content")
		panic(err)
	}

	InitRoutes(e, c)
	return e
}

func InitRoutes(r *gin.Engine, c *controller.Controller) {
	// TODO(low prio, effort not worth benefit): read the paths, add to constants and change @Route of swagger API doc in associated controllers
	// TODO: rate limit, if not for all then auth endpoints for sure
	authMw := middleware.NewAuthenticationMiddleware(c, true)
	authInfoMwUnblocking := middleware.NewAuthenticationMiddleware(c, false)
	i18nMw := middleware.NewLanguageMiddleware()

	// API endpoints
	var v1 = r.Group(docs.SwaggerInfo.BasePath)
	{
		// setup swagger API
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		authentication := v1.Group("/auth")
		{
			authentication.POST("/sign-in", c.SignIn)
			authentication.GET("/sign-out", c.SignOut)
		}

		// TODO: API point for websocket -- https://github.com/gin-gonic/examples/blob/master/websocket/server/server.go#L16
	}

	// --- html endpoints ---

	// public endpoints
	var public = r.Group("")
	{
		public.Use(i18nMw)
		public.Use(authInfoMwUnblocking)

		public.GET(g.RLogin, c.PageLogin)
		public.POST(g.RLogin, c.PageLogin)

		public.GET(g.RRegister, c.PageRegister)
		public.POST(g.RRegister, c.PageRegister)
	}

	// protected endpoints
	var p = r.Group("")
	{
		p.Use(i18nMw)
		p.Use(authMw)

		p.GET("/", c.PageMain)
		p.GET(g.RProfile, c.PageProfile)
		p.GET(g.RPayments, c.PagePayments)
		p.GET(g.RNotes, c.PageNotes)
		p.GET(g.RMessaging, c.PageMessaging)

		p.GET(g.RHouses, c.PageHouses)
		p.GET(g.RHouseID, c.PageHouse)

		p.GET(g.RHtmxRoomateSearch, c.HtmxRoomateSearch)
		p.POST(g.RHtmxRoomateSearch, c.HtmxRoomateSearch)

		p.GET(g.RHtmxHouseForm, c.GetHtmxHouseModal)
		p.POST(g.RHtmxHouseForm, c.PostHtmxHouseForm)
		p.PUT(g.RHtmxHouseForm, c.PutHtmxHouseForm)
		p.DELETE(g.RHtmxHouseForm, c.DeleteHouse)

		p.GET(g.RHtmxHouseResidentsBadge, c.HtmxHouseCardResidentsBadge)
	}

	r.Static("/assets", "./assets/public")
	r.StaticFile("/favicon.ico", "./assets/favicon/favicon.ico")
}
