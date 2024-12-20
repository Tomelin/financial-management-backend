package web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/internal/core/service"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type TenantHandlerHttpInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByFilterMany(c *gin.Context)
	GetByFilterOne(c *gin.Context)
}

type TenantHandlerHttp struct {
	Service service.ITenantService
}

func NewTenantHandlerHttp(svc *service.ITenantService, routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) TenantHandlerHttpInterface {

	lab := &TenantHandlerHttp{
		Service: *svc,
	}

	lab.handlers(routerGroup, middleware...)

	return lab

}

func (c *TenantHandlerHttp) handlers(routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) {
	middlewareList := make([]gin.HandlerFunc, len(middleware))
	for i, mw := range middleware {
		middlewareList[i] = mw
	}

	routerGroup.POST("/tenant", append(middlewareList, c.Create)...)
	routerGroup.GET("/tenant/:id", append(middlewareList, c.GetById)...)
	routerGroup.GET("/tenant/search", append(middlewareList, c.GetByFilterMany)...)
	routerGroup.GET("/tenant/filter", append(middlewareList, c.GetByFilterOne)...)
	routerGroup.GET("/tenant", append(middlewareList, c.Get)...)
	routerGroup.PUT("/tenant/:id", append(middlewareList, c.Update)...)
	routerGroup.DELETE("/tenant/:id", append(middlewareList, c.Delete)...)
}

// CreateTenantResponse    godoc
// @Summary     create a new lab destroy
// @Tags        Tenant
// @Accept       json
// @Produce     json
// @Description create a new lab destroy
// @Success     200 {object} entity.TenantResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /tenant [post]
func (obj *TenantHandlerHttp) Create(c *gin.Context) {

	var tenant entity.TenantResponse
	if err := c.BindJSON(&tenant); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, "body is required")
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	u, err := entity.NewTenant(&tenant)
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

// GetAllTenantResponse    godoc
// @Summary     get all lab destroy
// @Tags        Tenant
// @Accept       json
// @Produce     json
// @Description get all lab destroy
// @Success     200 {object} []entity.TenantResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /tenant [get]
func (obj *TenantHandlerHttp) Get(c *gin.Context) {

	response, err := obj.Service.Get()
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
// @Tags        Tenant
// @Accept      json
// @Produce     json
// @Param       id path string true "found"
// @Description get a lab destroy by ID
// @Success     200 {object} entity.TenantResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /tenant/{id} [get]
func (obj *TenantHandlerHttp) GetById(c *gin.Context) {

	tenantId := c.Param("id")
	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	response, err := obj.Service.GetById(&tenantId)
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
// @Tags        Tenant
// @Accept      json
// @Produce     json
// @Param       email query string false "email@domain.com"
// @Param       project query string false "project_name"
// @Param       available query boolean false "string default" default(false)
// @Description get a lab destroy by email or project
// @Success     200 {object} []entity.TenantResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /tenant/search [get]
func (obj *TenantHandlerHttp) GetByFilterMany(c *gin.Context) {

	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key and value are required"})
		c.Abort()
		return
	}

	filter := entity.QueryDB{
		Key:       key,
		Value:     value,
		Condition: string(entity.QueryFirebaseEqual),
	}

	response, err := obj.Service.GetByFilterMany(context.Background(), []entity.QueryDB{filter})
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Abort()
		return
	}

	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"messagem": "not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByFilterOne  godoc
// @Summary     get a lab destroy by email or project
// @Tags        Tenant
// @Accept      json
// @Produce     json
// @Param       email query string false "email@domain.com"
// @Param       project query string false "project_name"
// @Param       available query boolean false "string default" default(false)
// @Description get a lab destroy by email or project
// @Success     200 {object} []entity.TenantResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /tenant/search [get]
func (obj *TenantHandlerHttp) GetByFilterOne(c *gin.Context) {

	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key and value are required"})
		c.Abort()
		return
	}

	filter := entity.QueryDB{
		Key:       key,
		Value:     value,
		Condition: string(entity.QueryFirebaseEqual),
	}
	ctx := context.Background()
	response, err := obj.Service.GetByFilterOne(ctx, []entity.QueryDB{filter})
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Abort()
		return
	}

	if response == nil || response.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}

func (obj *TenantHandlerHttp) Update(c *gin.Context) {

	tenantId := c.Param("id")
	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	err := utils.ValidateUUID(&tenantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var tenant entity.TenantResponse
	if err := c.ShouldBindJSON(&tenant); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is required"})
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tenant.ID != tenantId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id in body must be equal to id in path"})
		c.Abort()
		return
	}

	data, err := obj.Service.Update(&tenant)
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

func (obj *TenantHandlerHttp) Delete(c *gin.Context) {

	tenantId := c.Param("id")
	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	err := obj.Service.Delete(&tenantId)

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
