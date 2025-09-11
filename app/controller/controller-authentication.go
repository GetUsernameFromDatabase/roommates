package controller

import (
	"net/http"
	"roommates/db/dbqueries"
	g "roommates/globals"
	"roommates/rdb"
	"roommates/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
// returns user id, if empty then 500 error occured
func (c *Controller) registerUser(ctx *gin.Context, user RegisterAccountRequest) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil { // this should never be an issue
		HandleServerError(ctx, err, "error processing password")
		return ""
	}

	userID, err := c.DB.InsertUser(ctx, dbqueries.InsertUserParams{
		Email:    user.Email,
		Username: user.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		// this should never be an issue
		HandleServerError(ctx, err, "error registering user")
		return ""
	}

	return userID.String()
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
		UserID:   credsInDb.ID.String(),
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

type RegisterAccountRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
	Username string `form:"username" json:"username" binding:"required"`
}

type RegisterAccountResponse struct {
	SignInResponse
	Message string `json:"username" binding:"required"`
}

// RegisterAccount godoc
//
//	@Summary      RegisterAccount new account
//	@Description  Creates a new user account
//	@Tags         auth
//
//	@Accept  json
//	@Param   RegisterAccount  body  RegisterAccountRequest  true  "Info used to register account"
//
//	@Produce      json
//	@Success  201  {object}  RegisterAccountResponse
//	@Failure  400  {object}  utils.HTTPError
//	@Failure  500  {object}  utils.HTTPError
//
//	@Security  ApiKeyAuth
//	@Router    /api/v1/auth/register-account [post]
func (c *Controller) RegisterAccount(ctx *gin.Context) {
	var req RegisterAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	userID := c.registerUser(ctx, req)
	if userID == "" {
		return
	}

	token := c.signUserIn(ctx, rdb.UserSessionValue{
		UserID:   userID,
		Username: req.Username,
	})
	ctx.JSON(http.StatusCreated, RegisterAccountResponse{
		Message:        "account created",
		SignInResponse: SignInResponse{Token: token.String()},
	})
}
