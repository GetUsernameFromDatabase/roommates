package main

import (
	"bytes"
	"context"
	"io/fs"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/events"
	"roommates/internal/handlers"
	"roommates/internal/mailer"
	"roommates/internal/middleware"
	"roommates/internal/push"
	"roommates/internal/routes"
	"roommates/internal/services"
	"roommates/internal/sessions"
	"roommates/internal/storage"
	"roommates/web/assets"
	"roommates/web/renderer"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// limitBody caps every request body, so no route reads an endless one
//
// form parsing hides the cut-off from handlers, so a declared size above max
// gets its 413 here and a chunked body is cut by the reader (a 400 in most handlers)
func limitBody(max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > max {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, max)
		c.Next()
	}
}

func initRoutes(r *gin.Engine, h *handlers.DefaultHandler, service *services.DefaultService) {
	mwI18n := middleware.NewInternationalization()
	// separate counters per kind, flatmates on one Wi-Fi share an IP
	// each attempt costs a bcrypt hash (~250ms): login, register, new password
	mwAuthLimit := middleware.NewRateLimit(10, time.Minute)
	// each sends an email, so mailboxes cannot be flooded
	mwEmailLimit := middleware.NewRateLimit(10, time.Minute)
	// invite codes are short enough to be worth guessing
	mwInviteLimit := middleware.NewRateLimit(20, time.Minute)
	r.Use(limitBody(services.MaxFileSize + 1<<20))
	r.Use(middleware.NewSecurityHeader())
	r.Use(middleware.NewCrossOriginProtection())
	// before every page, error pages too, so all of them render in the chosen theme
	r.Use(middleware.NewAppearance())

	registerAssets(r)
	public := r.Group("")
	{
		public.Use(middleware.NewGuestOnly())
		public.Use(mwI18n)

		public.GET(routes.Login, h.PageLogin)
		public.POST(routes.Login, mwAuthLimit, h.PageLogin)

		public.GET(routes.Register, h.PageRegister)
		public.POST(routes.Register, mwAuthLimit, h.PageRegister)

		public.GET(routes.ForgotPassword, h.PageForgotPassword)
		public.POST(routes.ForgotPassword, mwEmailLimit, h.PageForgotPassword)
	}

	// --- --- ---

	// email links work whether signed in or not
	// their tokens have 130 random bits, guessing needs no rate limit
	open := r.Group("")
	{
		open.Use(mwI18n)

		open.GET(routes.VerifyEmail, h.VerifyEmail)
		open.GET(routes.ResetPassword, h.PageResetPassword)
		open.POST(routes.ResetPassword, mwAuthLimit, h.PageResetPassword)
		open.GET(routes.ConfirmEmail, h.ConfirmEmailChange)
		open.POST(routes.Appearance, h.SetAppearance)
		open.POST(routes.Language, h.SetLanguage)
		open.GET(routes.Manifest, h.Manifest)
		open.GET(routes.Offline, h.PageOffline)
		open.GET(routes.AccountDeleted, h.PageAccountDeleted)
	}

	// --- --- ---

	// endpoints that need authentication
	// authorization should be checked by the handler
	priv := r.Group("")
	{
		priv.Use(middleware.NewOnlyAuthenticated())
		priv.Use(mwI18n)

		priv.GET(routes.Home, h.PageHome)
		priv.POST(routes.Households, h.PageHome)
		priv.POST(routes.Logout, h.Logout)
		priv.GET(routes.Settings, h.PageSettings)
		accountRoutes(priv.Group(routes.Settings), h, mwAuthLimit, mwEmailLimit)
		priv.GET(routes.DataExportDownload, h.DownloadDataExport)
		priv.GET(routes.UserAvatar, h.DownloadAvatar)
		priv.POST(routes.VerifyEmailResend, mwEmailLimit, h.ResendVerificationEmail)
		priv.POST(routes.PushSubscribe, h.SubscribePush)
		priv.POST(routes.PushUnsubscribe, h.UnsubscribePush)

		// outside the household group, whose read-only check would refuse it
		priv.POST(routes.HouseholdRestore, h.RestoreHousehold)

		priv.GET(routes.Invite, mwInviteLimit, h.PageInvite)
		priv.POST(routes.Invite, mwInviteLimit, h.PageInvite)

		// everything under it is scoped to one household the user is a member of
		household := priv.Group("")
		household.Use(middleware.NewHouseholdMember(service))
		{
			household.GET(routes.Household, h.PageHousehold)
			household.GET(routes.HouseholdPeople, h.PagePeople)
			household.GET(routes.HouseholdSettings, h.PageSettings)
			household.GET(routes.HouseholdActivity, h.PageActivity)
			accountRoutes(household.Group(routes.HouseholdSettings), h, mwAuthLimit, mwEmailLimit)
			household.POST(routes.HouseholdPaymentShow, h.SetShowPaymentDetails)
			household.POST(routes.HouseholdNotifications, h.SetNotification)
			household.POST(routes.HouseholdUpdate, h.UpdateHousehold)
			household.POST(routes.HouseholdDelete, h.DeleteHousehold)
			household.POST(routes.HouseholdStepDown, h.StepDown)
			household.POST(routes.MemberTransfer, h.TransferOwnership)
			household.POST(routes.Invites, h.CreateInvite)
			household.POST(routes.InvitesEmail, mwEmailLimit, h.SendInviteEmail)
			household.POST(routes.InviteDelete, h.DeleteInvite)
			household.POST(routes.MemberRemove, h.RemoveMember)
			household.POST(routes.MemberPromote, h.PromoteMember)
			household.POST(routes.HouseholdLeave, h.LeaveHousehold)

			household.POST(routes.ShoppingLists, h.CreateShoppingList)
			household.GET(routes.ShoppingList, h.PageShoppingList)
			household.POST(routes.ShoppingListDelete, h.DeleteShoppingList)
			household.POST(routes.ShoppingListName, h.RenameShoppingList)
			household.POST(routes.ShoppingListCategory, h.SetShoppingListCategory)
			household.POST(routes.ShoppingItems, h.AddShoppingItem)
			household.POST(routes.ShoppingListClear, h.ClearCheckedShoppingItems)
			household.POST(routes.ShoppingItemClaim, h.SetShoppingItemClaim)
			household.POST(routes.ShoppingItemCheck, h.SetShoppingItemChecked)
			household.POST(routes.ShoppingItemDelete, h.DeleteShoppingItem)
			household.GET(routes.ShoppingItem, h.PageShoppingItem)
			household.POST(routes.ShoppingItemDescription, h.SetShoppingItemDescription)
			household.POST(routes.ShoppingItemText, h.SetShoppingItemText)
			household.POST(routes.ShoppingItem+routes.FileRefUpload, h.UploadShoppingItemFile)
			h.FileRefRoutes(household.Group(routes.ShoppingItem), h.ShoppingItemFileTarget)

			household.GET(routes.Chores, h.PageChores)
			household.POST(routes.Chores, h.CreateChore)
			household.GET(routes.Chore, h.PageChore)
			household.POST(routes.Chore, h.UpdateChore)
			household.POST(routes.ChoreDelete, h.DeleteChore)
			household.POST(routes.ChoreDone, h.DoneChore)
			household.POST(routes.ChoreSkip, h.SkipChore)
			household.POST(routes.ChorePass, h.PassChore)
			household.POST(routes.ChoreAssign, h.AssignChore)
			household.POST(routes.ChoreMemberMove, h.MoveChoreMember)

			household.GET(routes.Calendar, h.PageCalendar)
			household.GET(routes.CalendarAgenda, h.CalendarAgenda)
			household.POST(routes.CalendarEvents, h.CreateEvent)
			household.GET(routes.CalendarEvent, h.PageEvent)
			household.POST(routes.CalendarEvent, h.UpdateEvent)
			household.POST(routes.CalendarEventDelete, h.DeleteEvent)
			household.POST(routes.CalendarEvent+routes.FileRefUpload, h.UploadEventFile)
			h.FileRefRoutes(household.Group(routes.CalendarEvent), h.EventFileTarget)

			household.GET(routes.Expenses, h.PageExpenses)
			household.GET(routes.ExpenseNew, h.PageNewExpense)
			household.POST(routes.Expenses, h.CreateExpense)
			household.GET(routes.ExpensesExport, h.ExportExpenses)
			household.GET(routes.ExpensesBalances, h.ExportBalances)
			household.GET(routes.Recurring, h.PageRecurring)
			household.GET(routes.RecurringExpense, h.PageEditRecurring)
			household.POST(routes.RecurringExpense, h.UpdateRecurring)
			household.POST(routes.RecurringExpenseDelete, h.DeleteRecurring)
			household.POST(routes.RecurringExpenseSkip, h.SkipRecurring)
			household.GET(routes.Expense, h.PageExpense)
			household.GET(routes.ExpenseEdit, h.PageEditExpense)
			household.POST(routes.Expense, h.UpdateExpense)
			household.POST(routes.ExpenseDelete, h.DeleteExpense)
			household.POST(routes.ExpenseSharePaid, h.SetExpenseSharePaid)
			household.POST(routes.Expense+routes.FileRefUpload, h.UploadExpenseFile)
			h.FileRefRoutes(household.Group(routes.Expense), h.ExpenseFileTarget)
			household.POST(routes.Settlements, h.CreateSettlement)
			household.POST(routes.SettlementConfirm, h.ConfirmSettlement)
			household.POST(routes.SettlementDelete, h.DeleteSettlement)

			household.GET(routes.Notes, h.PageNotes)
			household.GET(routes.NoteNew, h.PageNewNote)
			household.POST(routes.Notes, h.CreateNote)
			household.POST(routes.NoteRender, h.RenderNoteBody)
			household.GET(routes.Note, h.PageNote)
			household.GET(routes.NoteEdit, h.PageEditNote)
			household.POST(routes.NoteEdit, h.UpdateNote)
			household.POST(routes.NotePin, h.SetNotePinned)
			household.POST(routes.NoteDelete, h.DeleteNote)
			household.POST(routes.Note+routes.FileRefUpload, h.UploadNoteFile)
			h.FileRefRoutes(household.Group(routes.Note), h.NoteFileTarget)

			household.GET(routes.HouseholdEvents, h.Events)
			household.GET(routes.Chat, h.PageChat)
			household.POST(routes.Chat, h.PostMessage)
			household.GET(routes.ChatBadge, h.ChatBadge)
			household.GET(routes.ChatFiles, h.ChatFilePicks)
			household.POST(routes.ChatMessageDelete, h.DeleteMessage)

			household.GET(routes.Files, h.PageFiles)
			household.POST(routes.Files, h.UploadFile)
			household.GET(routes.File, h.PageFile)
			household.GET(routes.FileRaw, h.DownloadFile)
			household.GET(routes.FileThumb, h.DownloadFileThumb)
			household.POST(routes.FileName, h.RenameFile)
			household.POST(routes.FileKeep, h.KeepFile)
			household.POST(routes.FileDelete, h.DeleteFile)
		}
	}
}

// accountRoutes are the account forms of the settings page, at [routes.Settings] and [routes.HouseholdSettings]
func accountRoutes(settings *gin.RouterGroup, h *handlers.DefaultHandler, mwAuthLimit, mwEmailLimit gin.HandlerFunc) {
	settings.POST(routes.AccountAvatar, h.SetAvatar)
	settings.POST(routes.AccountAvatarDelete, h.RemoveAvatar)
	settings.POST(routes.AccountUsername, h.ChangeUsername)
	// bcrypt
	settings.POST(routes.AccountPassword, mwAuthLimit, h.ChangePassword)
	// bcrypt and a mail
	settings.POST(routes.AccountEmail, mwAuthLimit, mwEmailLimit, h.ChangeEmail)
	settings.POST(routes.AccountSignOutEverywhere, h.SignOutEverywhere)
	settings.POST(routes.AccountPayment, h.SetPaymentDetails)
	// bcrypt and a mail
	settings.POST(routes.AccountDelete, mwAuthLimit, mwEmailLimit, h.DeleteAccount)
	// a mail when it is built
	settings.POST(routes.AccountDataExport, mwEmailLimit, h.RequestDataExport)
}

// uses [assets.PublicAssetsFS]
//
//	"/assets/public" is served at [routes.Assets]
//	"/favicon/favicon.ico" is served at [routes.Favicon]
//	"/favicon/*.png" and "/favicon/favicon.svg" are served at [routes.Icons]
//	"/pwa/sw.js" is served at [routes.ServiceWorker]
//
// and [routes.JSSource] at [routes.RoutesJS] (pages load [routes.JSURL])
func registerAssets(r *gin.Engine) {
	assetFS := assets.PublicAssetsFS
	publicFS, err := fs.Sub(assetFS, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("could not get public assets")
	}

	assetRoutes := r.Group(routes.Assets, func(c *gin.Context) {
		cacheHeaders(c, "public/"+strings.TrimPrefix(c.Request.URL.Path, routes.Assets+"/"))
	})
	assetRoutes.StaticFS("/", http.FS(publicFS))
	file := func(route, name string) {
		h := func(c *gin.Context) {
			cacheHeaders(c, name)
			c.FileFromFS(name, http.FS(assetFS))
		}
		r.GET(route, h)
		r.HEAD(route, h)
	}
	file(routes.Favicon, "favicon/favicon.ico")
	for _, icon := range []string{"web-app-manifest-192x192.png", "web-app-manifest-512x512.png", "apple-touch-icon.png", "favicon.svg"} {
		file(routes.Icons+"/"+icon, "favicon/"+icon)
	}
	file(routes.ServiceWorker, "pwa/sw.js")

	serveJS := func(c *gin.Context) {
		c.Header("ETag", `"`+routes.JSHash+`"`)
		// pages load [routes.JSURL], sw.js the plain path, which must be checked for changes
		if c.Query("v") == routes.JSHash {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			c.Header("Cache-Control", "no-cache")
		}
		c.Header("Content-Type", "text/javascript; charset=utf-8")
		// answers If-None-Match with a 304
		http.ServeContent(c.Writer, c.Request, "", time.Time{}, bytes.NewReader(routes.JSSource))
	}
	r.GET(routes.RoutesJS, serveJS)
	r.HEAD(routes.RoutesJS, serveJS)
}

// cacheHeaders gives an embedded file its content hash as ETag, so a repeat visit gets a 304,
// and lets the browser keep it for a year when the URL carries that hash ([assets.URL]) or the
// shadcn bundle's name does; a missing file gets nothing and stays uncached for the next deploy
func cacheHeaders(c *gin.Context, name string) {
	hash := assets.Hash(name)
	if hash == "" {
		return
	}
	c.Header("ETag", `"`+hash+`"`)
	if c.Query("v") == hash || strings.HasPrefix(name, "public/js/") {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		c.Header("Cache-Control", "no-cache")
	}
}

// -----------------------------------------------------------------------------

// InitGinEngine runs the background loops until ctx is done,
// main waits on the returned service's pending pushes before closing dbpool
func InitGinEngine(ctx context.Context, dbpool *pgxpool.Pool) (*gin.Engine, *services.DefaultService) {
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
	switch platform := g.EVK_TrustedPlatform.Get(""); platform {
	case "":
	case "cloudflare":
		r.TrustedPlatform = gin.PlatformCloudflare
	case "flyio":
		r.TrustedPlatform = gin.PlatformFlyIO
	default:
		log.Fatal().Str("value", platform).Msg(g.EVK_TrustedPlatform.String() + " must be cloudflare, flyio or unset")
	}
	renderer.ReplaceGinRenderer(r)

	// before the session middleware, so probes neither need a login nor create sessions
	healthz := func(c *gin.Context) {
		pingCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := dbpool.Ping(pingCtx); err != nil {
			c.String(http.StatusServiceUnavailable, "database unreachable")
			return
		}
		c.String(http.StatusOK, "ok")
	}
	r.GET(routes.Healthz, healthz)
	r.HEAD(routes.Healthz, healthz)

	r.Use(sessions.NewPostgreSessionStoreHandler(dbpool))
	// panics now if APP_BASE_URL is missing in release, not on the first email
	_ = g.GetAppBaseURL()
	// panics now if the object storage is missing, not on the first upload
	service := services.NewDefaultService(dbpool, mailer.New(), storage.New(ctx), push.New())
	go service.RunCleanup(ctx, time.Minute)
	hub := events.NewHub()
	go hub.Run(ctx, dbpool)
	handler := handlers.NewDefaultHandler(service, hub)
	initRoutes(r, handler, service)
	return r, service
}
