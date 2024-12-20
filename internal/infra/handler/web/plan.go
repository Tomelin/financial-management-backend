package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type PlanHandlerHttpInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByFilterMany(c *gin.Context)
	GetByFilterOne(c *gin.Context)
}

type PlanHandlerHttp struct {
	Service entity.IPlan
}

func NewPlanHandlerHttp(svc *entity.IPlan, routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) PlanHandlerHttpInterface {

	lab := &PlanHandlerHttp{
		Service: *svc,
	}

	lab.handlers(routerGroup, middleware...)

	return lab

}

func (c *PlanHandlerHttp) handlers(routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) {
	middlewareList := make([]gin.HandlerFunc, len(middleware))
	for i, mw := range middleware {
		middlewareList[i] = mw
	}

	routerGroup.POST("/plan", append(middlewareList, c.Create)...)
	routerGroup.GET("/plan/:id", append(middlewareList, c.GetById)...)
	routerGroup.GET("/plan/search", append(middlewareList, c.GetByFilterMany)...)
	routerGroup.GET("/plan/filter", append(middlewareList, c.GetByFilterOne)...)
	routerGroup.GET("/plan", append(middlewareList, c.Get)...)
	routerGroup.PUT("/plan/:id", append(middlewareList, c.Update)...)
	routerGroup.DELETE("/plan/:id", append(middlewareList, c.Delete)...)
}

// CreatePlanResponse    godoc
// @Summary     create a new lab destroy
// @Tags        Plan
// @Accept       json
// @Produce     json
// @Description create a new lab destroy
// @Success     200 {object} entity.PlanResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /plan [post]
func (obj *PlanHandlerHttp) Create(c *gin.Context) {

	var plan entity.PlanResponse
	if err := c.BindJSON(&plan); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, "body is required")
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	u, err := entity.NewPlan(&plan)
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

// GetAllPlanResponse    godoc
// @Summary     get all lab destroy
// @Tags        Plan
// @Accept       json
// @Produce     json
// @Description get all lab destroy
// @Success     200 {object} []entity.PlanResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /plan [get]
func (obj *PlanHandlerHttp) Get(c *gin.Context) {

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
// @Tags        Plan
// @Accept      json
// @Produce     json
// @Param       id path string true "found"
// @Description get a lab destroy by ID
// @Success     200 {object} entity.PlanResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /plan/{id} [get]
func (obj *PlanHandlerHttp) GetById(c *gin.Context) {

	planId := c.Param("id")
	if planId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	response, err := obj.Service.GetById(&planId)
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
// @Tags        Plan
// @Accept      json
// @Produce     json
// @Param       email query string false "email@domain.com"
// @Param       project query string false "project_name"
// @Param       available query boolean false "string default" default(false)
// @Description get a lab destroy by email or project
// @Success     200 {object} []entity.PlanResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /plan/search [get]
func (obj *PlanHandlerHttp) GetByFilterMany(c *gin.Context) {

	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key and value are required"})
		c.Abort()
		return
	}

	response, err := obj.Service.GetByFilterMany(key, &value)
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
// @Tags        Plan
// @Accept      json
// @Produce     json
// @Param       email query string false "email@domain.com"
// @Param       project query string false "project_name"
// @Param       available query boolean false "string default" default(false)
// @Description get a lab destroy by email or project
// @Success     200 {object} []entity.PlanResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /plan/search [get]
func (obj *PlanHandlerHttp) GetByFilterOne(c *gin.Context) {

	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key and value are required"})
		c.Abort()
		return
	}

	response, err := obj.Service.GetByFilterOne(key, &value)
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

func (obj *PlanHandlerHttp) Update(c *gin.Context) {

	planId := c.Param("id")
	if planId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	err := utils.ValidateUUID(&planId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var plan entity.PlanResponse
	if err := c.ShouldBindJSON(&plan); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is required"})
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if plan.ID != planId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id in body must be equal to id in path"})
		c.Abort()
		return
	}

	data, err := obj.Service.Update(&plan)
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

func (obj *PlanHandlerHttp) Delete(c *gin.Context) {

	planId := c.Param("id")
	if planId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	err := obj.Service.Delete(&planId)

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
