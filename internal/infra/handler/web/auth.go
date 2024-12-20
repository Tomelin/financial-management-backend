package web

import (
	"context"
	"net/http"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"

	"github.com/Tomelin/financial-management-backend/internal/core/entity"
	"github.com/Tomelin/financial-management-backend/internal/core/service"
	"github.com/Tomelin/financial-management-backend/pkg/authProvider"
)

type IAuthHandlerHttp interface {
	Callback(c *gin.Context)
	Logout(c *gin.Context)
	Auth(c *gin.Context)
	IsLoggedIn(c *gin.Context)
	ValidateToken(c *gin.Context)
}

type AuthHandlerHttp struct {
	User         service.IUserService
	AuthProvider authProvider.IAuthProvider
}

func NewAuthHandlerHttp(ap authProvider.IAuthProvider, user service.IUserService, routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) IAuthHandlerHttp {

	lab := &AuthHandlerHttp{
		User:         user,
		AuthProvider: ap,
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
	routerGroup.GET("/v1/auth/:provider", append(middlewareList, c.Auth)...)
	routerGroup.GET("/v1/auth/:provider/is_logged_in", append(middlewareList, c.IsLoggedIn)...)
}

func (obj *AuthHandlerHttp) Callback(c *gin.Context) {
	c.Set("provider", c.Param("provider"))

	userFromProvider, err := obj.AuthProvider.Callback(c.Writer, c.Request)
	if err != nil {
		obj.AuthProvider.Logout(c, c.Writer, c.Request)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	b, err := json.Marshal(userFromProvider)
	if err != nil {
		obj.AuthProvider.Logout(c, c.Writer, c.Request)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user entity.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		obj.AuthProvider.Logout(c, c.Writer, c.Request)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := entity.NewUser(&user)
	if err != nil {
		obj.AuthProvider.Logout(c, c.Writer, c.Request)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	_, err = obj.User.Create(response)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Redirect(c.Writer, c.Request, "http://localhost:5173", http.StatusTemporaryRedirect)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		http.Redirect(c.Writer, c.Request, "http://localhost:5173/error", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(c.Writer, c.Request, "http://localhost:5173/form", http.StatusTemporaryRedirect)
}

func (obj *AuthHandlerHttp) Logout(c *gin.Context) {
	err := obj.AuthProvider.Logout(c, c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Writer.Header().Set("Location", "http://localhost:5173")
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (obj *AuthHandlerHttp) Auth(c *gin.Context) {

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Writer.Flush()
		c.Abort()
		return
	}

	if c.GetHeader("Authorization") != user.AccessToken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not authorized"})
		c.Writer.Flush()
		c.Abort()
		return
	}

	c.Next()
}
