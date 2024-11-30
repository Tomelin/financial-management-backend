package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/internal/core/service"
)

type UserHandlerHttpInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByFilterMany(c *gin.Context)
	GetByFilterOne(c *gin.Context)
}

type UserHandlerHttp struct {
	Service service.IUserService
}

func NewUserHandlerHttp(svc *service.IUserService, routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) UserHandlerHttpInterface {

	lab := &UserHandlerHttp{
		Service: *svc,
	}

	lab.handlers(routerGroup, middleware...)

	return lab

}

func (c *UserHandlerHttp) handlers(routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) {
	middlewareList := make([]gin.HandlerFunc, len(middleware))
	for i, mw := range middleware {
		middlewareList[i] = mw
	}

	routerGroup.POST("/user", append(middlewareList, c.Create)...)
	routerGroup.GET("/user/:id", append(middlewareList, c.GetById)...)
	routerGroup.GET("/user/search", append(middlewareList, c.GetByFilterMany)...)
	routerGroup.GET("/user", append(middlewareList, c.Get)...)
	routerGroup.PUT("/user/:id", append(middlewareList, c.Update)...)
	routerGroup.DELETE("/user/:id", append(middlewareList, c.Delete)...)
}

// CreateUserResponse    godoc
// @Summary     create a new lab destroy
// @Tags        User
// @Accept       json
// @Produce     json
// @Description create a new lab destroy
// @Success     200 {object} entity.EntityResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /User [post]
func (obj *UserHandlerHttp) Create(c *gin.Context) {

	var user entity.User
	if err := c.BindJSON(&user); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, "body is required")
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	u, err := entity.NewUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	result, err := obj.Service.Create(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	c.JSON(http.StatusAccepted, result)
}

// GetAllUserResponse    godoc
// @Summary     get all lab destroy
// @Tags        User
// @Accept       json
// @Produce     json
// @Description get all lab destroy
// @Success     200 {object} []entity.EntityResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /User [get]
func (obj *UserHandlerHttp) Get(c *gin.Context) {

	users, err := obj.Service.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	c.JSON(http.StatusAccepted, users)
}

// GetById      godoc
// @Summary     get a lab destroy by ID
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path string true "found"
// @Description get a lab destroy by ID
// @Success     200 {object} entity.EntityResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /User/{id} [get]
func (obj *UserHandlerHttp) GetById(c *gin.Context) {

	id := c.Param("id")
	user, err := obj.Service.GetById(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	c.JSON(http.StatusAccepted, user)
}

// GetByFilterMany  godoc
// @Summary     get a lab destroy by email or project
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       email query string false "email@domain.com"
// @Param       project query string false "project_name"
// @Param       username query string false "username"
// @Param       available query boolean false "string default" default(false)
// @Description get a lab destroy by email or project
// @Success     200 {object} []entity.EntityResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /User/search [get]
func (obj *UserHandlerHttp) GetByFilterMany(c *gin.Context) {

	c.JSON(http.StatusAccepted, "result")
}

// GetByFilterOne  godoc
// @Summary     get a lab destroy by email or project
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       email query string false "email@domain.com"
// @Param       project query string false "project_name"
// @Param       username query string false "username"
// @Param       available query boolean false "string default" default(false)
// @Description get a lab destroy by email or project
// @Success     200 {object} []entity.EntityResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /User/search [get]
func (obj *UserHandlerHttp) GetByFilterOne(c *gin.Context) {

	c.JSON(http.StatusAccepted, "result")
}

func (obj *UserHandlerHttp) Update(c *gin.Context) {

	var user entity.UserResponse
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	result, err := obj.Service.Update(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	c.JSON(http.StatusAccepted, result)
}

func (obj *UserHandlerHttp) Delete(c *gin.Context) {

	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, "id not found")
		c.Abort()
		return
	}

	id := c.Param("id")
	err := obj.Service.Delete(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	c.JSON(http.StatusAccepted, "removed")
}
