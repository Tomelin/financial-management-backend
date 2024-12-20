package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Tomelin/financial-management-backend/internal/core/entity"
	"github.com/Tomelin/financial-management-backend/pkg/observability"
	"github.com/Tomelin/financial-management-backend/pkg/utils"
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
	Service entity.IUser
	trace   *observability.Tracer
}

func NewUserHandlerHttp(svc *entity.IUser, trace *observability.Tracer, routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) UserHandlerHttpInterface {

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

	routerGroup.POST("/user", c.Create)
	routerGroup.GET("/user/:id", append(middlewareList, c.GetById)...)
	routerGroup.GET("/user/search", append(middlewareList, c.GetByFilterMany)...)
	routerGroup.GET("/user/filter", append(middlewareList, c.GetByFilterOne)...)
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
	ctx, span := obj.trace.Trace.Start(c.Request.Context(), fmt.Sprintf("%s.Create", string(entity.ApplicationLayerHandler)))
	defer span.End()

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

	result, err := obj.Service.Create(ctx, u)
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
	ctx := c.Request.Context()
	response, err := obj.Service.Get(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if len(response) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
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

	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}
	ctx := c.Request.Context()
	response, err := obj.Service.GetById(ctx, &userId)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
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

	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key and value are required"})
		c.Abort()
		return
	}

	filter := []entity.QueryDBClause{
		{
			Clause: entity.QueryClauseAnd,
			Queries: []entity.QueryDB{
				{
					Key:       key,
					Value:     value,
					Condition: string(entity.QueryFirebaseEqual),
				},
			},
		},
	}
	ctx := c.Request.Context()
	response, err := obj.Service.GetByFilterMany(ctx, filter)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
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

	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key and value are required"})
		c.Abort()
		return
	}

	filter := []entity.QueryDBClause{
		{
			Clause: entity.QueryClauseAnd,
			Queries: []entity.QueryDB{
				{
					Key:       key,
					Value:     value,
					Condition: string(entity.QueryFirebaseEqual),
				},
			},
		},
	}

	ctx := c.Request.Context()
	response, err := obj.Service.GetByFilterOne(ctx, filter)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}

func (obj *UserHandlerHttp) Update(c *gin.Context) {

	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	err := utils.ValidateUUID(&userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var user entity.AccountUser
	if err := c.ShouldBindJSON(&user); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is required"})
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ID != userId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id in body must be equal to id in path"})
		c.Abort()
		return
	}

	ctx := c.Request.Context()
	data, err := obj.Service.Update(ctx, &user)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, data)
}

func (obj *UserHandlerHttp) Delete(c *gin.Context) {

	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}
	ctx := c.Request.Context()
	err := obj.Service.Delete(ctx, &userId)

	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "deleted"})
}
