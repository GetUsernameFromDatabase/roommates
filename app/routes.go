package main

import (
	"roomates/controller"
	"roomates/docs"
	"roomates/gintemplrenderer"
	"roomates/locales"
	"roomates/middleware"

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
	// TODO(low prio, complexity high): read the paths, add to constants and change @Route of swagger API doc in associated controllers
	// TODO: rate limit, if not for all then auth endpoints for sure
	autMw := middleware.NewAuthenticationMiddleware(c, true)
	autMwUnblocking := middleware.NewAuthenticationMiddleware(c, false)
	i18nMw := middleware.NewLanguageMiddleware()
	// --- API endpoints ---
	v1 := r.Group(docs.SwaggerInfo.BasePath)
	{
		// setup swagger API
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		house := v1.Group("/house")
		{
			house.Use(autMw)
			house.GET("", nil)

			housePayments := house.Group("/payments")
			{
				housePayments.GET("", nil)
				// TODO: post payment file https://gin-gonic.com/en/docs/examples/upload-file/single-file/
			}

			houseNotes := house.Group("/notes")
			{
				houseNotes.GET("", nil)
			}
		}

		authentication := v1.Group("/auth")
		{
			authentication.POST("/sign-in", c.SignIn)
			authentication.GET("/sign-out", c.SignOut)
			authentication.POST("/register-account", c.RegisterAccount)
		}

		// TODO: API point for websocket -- https://github.com/gin-gonic/examples/blob/master/websocket/server/server.go#L16
	}

	// --- html endpoints ---
	// protected endpoints
	pr := r.Group("")
	{
		pr.Use(i18nMw)
		pr.Use(autMw)

		pr.GET("/", c.PageMain)
	}

	// public endpoints
	pu := r.Group("")
	{
		pu.Use(i18nMw)
		pu.Use(autMwUnblocking)
		pu.GET("/login", c.PageLogin)
		pu.POST("/login", c.PageLogin)

		pu.GET("/register", c.PageRegister)
		pu.POST("/register", c.PageRegister)
	}

	r.Static("/assets", "./assets/public")
	r.StaticFile("/favicon.ico", "./assets/favicon/favicon.ico")
}
