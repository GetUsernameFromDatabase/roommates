package controller

import (
	"net/http"
	"roommates/components"
	"roommates/db/dbqueries"
	g "roommates/globals"
	"roommates/locales"
	"roommates/middleware"
	"roommates/models"
	"roommates/rdb"
	"roommates/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (c *Controller) shouldUserBeSignedIn(ctx *gin.Context, req SignInRequest) (*dbqueries.GetUserCredentialsRow, error) {
	credsInDb, err := c.DB.GetUserCredentials(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, g.ErrorInvalidCredential
		}
		return nil, err
	}

	// no need to check email as it's used to get stored credentials
	if err = bcrypt.CompareHashAndPassword([]byte(credsInDb.Password),
		[]byte(req.Password)); err != nil {
		return nil, g.ErrorInvalidCredential
	}

	return &credsInDb, nil
}

func (c *Controller) signUserIn(ctx *gin.Context, sessionValue rdb.UserSessionValue) uuid.UUID {
	token := c.RH.CreateUserSession(ctx, sessionValue)
	ctx.SetCookie(
		string(g.CSessionToken),
		token.String(),
		int(rdb.EUserSession.Seconds()),
		"/",
		"localhost",
		true,
		true,
	)
	return token
}

// will register user and also respond with server errors
//
//   - on all errors string will be empty
//   - if error is nil and string is empty then response has been sent
//   - response has not been sent if error is ErrorAccountAlreadyExists
func (c *Controller) registerUser(ctx *gin.Context, user dbqueries.InsertUserParams) (pgtype.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil { // this should never be an issue
		HandleServerError(ctx, err, "error processing password")
		return pgtype.UUID{}, nil
	}

	userID, err := c.DB.InsertUser(ctx, dbqueries.InsertUserParams{
		Email:    user.Email,
		Username: user.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		switch err := err.(type) {
		case *pgconn.PgError:
			// https://www.postgresql.org/docs/current/errcodes-appendix.html
			if err.Code == "23505" {
				return pgtype.UUID{}, g.ErrorAccountAlreadyExists
			}
		default:
			HandleServerError(ctx, err, "error registering user")
		}
		return pgtype.UUID{}, nil
	}

	return userID, nil
}

//------------------------------------------------------------------------------

type SignInRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type SignInResponse struct {
	Token string `json:"token" binding:"required"`
}

// SignIn godoc
//
//	@Summary      User login
//	@Description  Authenticates a user with email and password
//	@Tags         auth
//
//	@Accept  json
//	@Param    SignIn  body  SignInRequest  true  "Information used for sign-in"
//
//	@Produce  json
//	@Success  200  {object}  SignInResponse
//	@Failure  400  {object}  utils.HTTPError
//	@Failure  401  {object}  utils.HTTPError
//	@Failure  500  {object}  utils.HTTPError
//
//	@Security  ApiKeyAuth
//	@Router    /api/v1/auth/sign-in [post]
func (c *Controller) SignIn(ctx *gin.Context) {
	var req SignInRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	credsInDb, err := c.shouldUserBeSignedIn(ctx, req)
	if err != nil {
		if errors.Is(err, g.ErrorInvalidCredential) {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, err)
			return
		}
		HandleServerError(ctx, err, "error fetching user credentials")
		return
	}

	token := c.signUserIn(ctx, rdb.UserSessionValue{
		UserID:   credsInDb.ID,
		Username: credsInDb.Username,
	})
	ctx.JSON(http.StatusOK, SignInResponse{Token: token.String()})
}

//------------------------------------------------------------------------------

// SignOut godoc
//
//	@Summary      User login
//	@Description  Authenticates a user with email and password
//	@Tags         auth
//
//	@Accept   json
//
//	@Produce  json
//	@Success  200  {object}  SimpleResponse
//	@Failure  400  {object}  utils.HTTPError
//	@Failure  401  {object}  utils.HTTPError
//	@Failure  500  {object}  utils.HTTPError
//
//	@Security  ApiKeyAuth
//	@Router    /api/v1/auth/sign-out [get]
func (c *Controller) SignOut(ctx *gin.Context) {
	utils.DeleteCookie(ctx, g.CSessionToken)
	token := utils.GetAuthTokenFromHeader(ctx)
	c.RH.DeleteUserSession(ctx, token)
	ctx.JSON(http.StatusOK, SimpleResponse{Message: "successfully signed out"})
}

//------------------------------------------------------------------------------

func (c *Controller) PageLogin(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)
	if authInfo != nil {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	method := ctx.Request.Method
	render := func(model models.Login) {
		page := components.PageLogin(model)
		RenderTempl(ctx, page)
	}

	switch method {
	case http.MethodGet:
		render(models.Login{ModelBase: models.ModelBase{Initial: true}})
	case http.MethodPost:
		var model models.Login
		ctx.ShouldBind(&model)

		isValid, _ := model.IsValid()
		if !isValid {
			render(model)
			return
		}

		credsInDb, err := c.shouldUserBeSignedIn(ctx, SignInRequest{
			Email:    model.Email,
			Password: model.Password,
		})
		if err != nil {
			if errors.Is(err, g.ErrorInvalidCredential) {
				model.Error = utils.T(
					ctx.Request.Context(),
					locales.LKFormsErrorInvalidCredential,
					// this error is safe to output publically
					err.Error(),
				)
				render(model)
				return
			}
			HandleServerError(ctx, err, "error fetching user credentials")
			return
		}

		c.signUserIn(ctx, rdb.UserSessionValue{
			UserID:   credsInDb.ID,
			Username: credsInDb.Username,
		})

		utils.Redirect(ctx, "/")
	default:
		ctx.String(http.StatusMethodNotAllowed, "method %s not allowed", method)
	}
}

func (c *Controller) PageRegister(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)
	if authInfo != nil {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	method := ctx.Request.Method
	render := func(model models.Register) {
		page := components.PageRegister(model)
		RenderTempl(ctx, page)
	}

	switch method {
	case http.MethodGet:
		// this is ridiculous
		render(models.Register{Login: models.Login{ModelBase: models.ModelBase{Initial: true}}})
	case http.MethodPost:
		var model models.Register
		ctx.ShouldBind(&model)

		isValid, _ := model.IsValid()
		if !isValid {
			render(model)
			return
		}

		userID, err := c.registerUser(ctx, dbqueries.InsertUserParams{
			Email:            model.Email,
			Password:         model.Password,
			Username:         model.Username,
			FullName:         &model.FullName,
			IsFullNamePublic: model.IsFullNamePublic,
		})
		if errors.Is(err, g.ErrorAccountAlreadyExists) {
			model.Error = utils.T(
				ctx.Request.Context(),
				locales.LKFormsErrorAlreadyExists,
				// this error is safe to output publically
				err.Error(),
			)
			render(model)
			return
		}
		if !userID.Valid {
			return
		}

		c.signUserIn(ctx, rdb.UserSessionValue{
			UserID:   userID,
			Username: model.Username,
		})
		utils.Redirect(ctx, "/")
	default:
		ctx.String(http.StatusMethodNotAllowed, "method %s not allowed", method)
	}
}
