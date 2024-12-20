package web

import (
	"context"
	"net/http"

	"encoding/json"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	middleware "github.com/synera-br/financial-management/src/backend/internal/infra/handler/middleware/authorization"
	"github.com/synera-br/financial-management/src/backend/pkg/authProvider"
	"github.com/synera-br/financial-management/src/backend/pkg/logger"
)

type IAuthHandlerHttp interface {
	Callback(c *gin.Context)
	Logout(c *gin.Context)
	Login(c *gin.Context)
	IsLoggedIn(c *gin.Context)
	ValidateToken(c *gin.Context)
}

type AuthHandlerHttp struct {
	User         entity.IUser
	AuthProvider authProvider.IAuthProvider
	Log          logger.Logger
	tokenJWT     entity.IAuthorization
}

func NewAuthenticationHandlerHttp(ap authProvider.IAuthProvider, l logger.Logger, tokenJWT entity.IAuthorization, user entity.IUser, routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) IAuthHandlerHttp {

	lab := &AuthHandlerHttp{
		User:         user,
		AuthProvider: ap,
		tokenJWT:     tokenJWT,
		Log:          l,
	}

	lab.handlers(routerGroup, middleware...)

	return lab
}

func (c *AuthHandlerHttp) handlers(routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) {
	middlewareList := make([]gin.HandlerFunc, len(middleware))
	for i, mw := range middleware {
		middlewareList[i] = mw
	}

	routerGroup.GET("/v1/auth/:provider/callback", append(middlewareList, c.Callback)...)
	routerGroup.GET("/v1/auth/:provider/logout", append(middlewareList, c.Logout)...)
	routerGroup.GET("/v1/auth/:provider", append(middlewareList, c.Login)...)
	routerGroup.GET("/v1/auth/:provider/is_logged_in", append(middlewareList, c.IsLoggedIn)...)
}

func (obj *AuthHandlerHttp) Callback(c *gin.Context) {
	c.Set("provider", c.Param("provider"))

	userFromProvider, err := obj.AuthProvider.Callback(c.Writer, c.Request)
	if err != nil {
		obj.AuthProvider.Logout(c, c.Writer, c.Request)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: err.Error(), Code: logger.ResponseCodeUnauthorized}).Error(),
		})
		return
	}
	if userFromProvider == nil {
		obj.AuthProvider.Logout(c, c.Writer, c.Request)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: "user not found", Code: logger.ResponseCodeUnauthorized}).Error(),
		})
		return
	}

	b, err := json.Marshal(userFromProvider)
	if err != nil {
		obj.AuthProvider.Logout(c, c.Writer, c.Request)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: err.Error(), Code: logger.ResponseCodeInternalServer}).Error(),
		})
		return
	}

	var user entity.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		obj.AuthProvider.Logout(c, c.Writer, c.Request)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: err.Error(), Code: logger.ResponseCodeInternalServer}).Error(),
		})
		return
	}

	response, err := entity.NewUser(&user)
	if err != nil {
		obj.AuthProvider.Logout(c, c.Writer, c.Request)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: err.Error(), Code: logger.ResponseCodeInternalServer}).Error(),
		})
		return
	}

	ctx := c.Request.Context()
	getuser, err := obj.User.GetByEmail(ctx, &user.Email)
	if err != nil {
		if err.Error() != "not found" {
			obj.AuthProvider.Logout(c, c.Writer, c.Request)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	if getuser != nil {
		getuser.ID = response.ID
	}

	token, err := obj.tokenJWT.GenerateTokenJWT(context.Background(), &entity.AuthorizationClaims{
		Email:     response.Email,
		UserID:    response.ID,
		Roles:     []entity.AccountRoles{},
		IsRevoked: false,
		// StandardClaims: jwt.StandardClaims{
		// 	ExpiresAt: userFromProvider.ExpiresAt.Unix() - 300, // diminuir 5 minutos
		// },
		Username: response.Name,
	}, response)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: err.Error(), Code: logger.ResponseCodeInternalServer}).Error(),
		})
		http.Redirect(c.Writer, c.Request, "http://localhost:5173/error", http.StatusTemporaryRedirect)
		return
	}

	store := sessions.NewCookieStore([]byte(*token))
	c.SetCookie("Authorization", *token, 3600, "/", "", true, true) // Ajuste os parâmetros do cookie conforme necessário

	sessions.Default(c).Set("Authorization", store)
	sessions.Default(c).Save()

	_, err = obj.User.Create(ctx, response)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Redirect(c.Writer, c.Request, "http://localhost:5173", http.StatusTemporaryRedirect)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: err.Error(), Code: logger.ResponseCodeInternalServer}).Error(),
		})
		http.Redirect(c.Writer, c.Request, "http://localhost:5173/error", http.StatusTemporaryRedirect)
		return
	}

	err = obj.tokenJWT.StoreTokenJWT(context.Background(), []byte(*token), &response.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: err.Error(), Code: logger.ResponseCodeInternalServer}).Error(),
		})
		http.Redirect(c.Writer, c.Request, "http://localhost:5173/error", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(c.Writer, c.Request, "http://localhost:5173/form", http.StatusTemporaryRedirect)
}

func (obj *AuthHandlerHttp) Logout(c *gin.Context) {
	err := obj.AuthProvider.Logout(c, c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: err.Error(), Code: logger.ResponseCodeUnauthorized}).Error(),
		})
		return
	}
	cookie, _ := c.Cookie("Authorization")

	middleware.OpenTokenJWT(&cookie)

	session := sessions.Default(c)
	c.SetCookie("Authorization", "", -1, "/", "", true, true)

	session.Delete("Authorization")
	session.Save()

	c.Writer.Header().Set("Location", "http://localhost:5173")
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (obj *AuthHandlerHttp) Login(c *gin.Context) {

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))
	c.Set("provider", c.Param("provider"))

	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"user": gothUser,
		})
	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func (obj *AuthHandlerHttp) IsLoggedIn(c *gin.Context) {

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))
	c.Set("provider", c.Param("provider"))

	user, err := obj.AuthProvider.IsLoggedIn(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"logged_in": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logged_in":    true,
		"user":         user,
		"access_token": user.AccessToken,
	})
}

func (obj *AuthHandlerHttp) ValidateToken(c *gin.Context) {

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))
	c.Set("provider", c.Param("provider"))

	user, err := obj.AuthProvider.IsLoggedIn(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: "invalid token", Code: logger.ResponseCodeBadRequest}).Error(),
		})
		c.Writer.Flush()
		c.Abort()
		return
	}

	if c.GetHeader("Authorization") != user.AccessToken {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": obj.Log.Error(&logger.Message{Body: "invalid token", Code: logger.ResponseCodeBadRequest}).Error(),
		})
		c.Writer.Flush()
		c.Abort()
		return
	}

	c.Next()
}
