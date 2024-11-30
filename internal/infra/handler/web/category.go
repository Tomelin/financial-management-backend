package web

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type HandlerHttpInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByFilterMany(c *gin.Context)
	GetByFilterOne(c *gin.Context)
}

type CategoryHandlerHttp struct {
	Service entity.ICategory
}

func NewCategoryHandlerHttp(svc entity.ICategory, routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) HandlerHttpInterface {

	load := &CategoryHandlerHttp{
		Service: svc,
	}

	load.handlers(routerGroup, middleware...)

	return load
}

func (cat *CategoryHandlerHttp) handlers(routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) {
	middlewareList := make([]gin.HandlerFunc, len(middleware))
	for i, mw := range middleware {
		middlewareList[i] = mw
	}

	routerGroup.POST("/category", append(middlewareList, cat.Create)...)
	routerGroup.GET("/category/", append(middlewareList, cat.Get)...)
	routerGroup.GET("/category/:id", append(middlewareList, cat.GetById)...)
	routerGroup.PUT("/category/:id", append(middlewareList, cat.Update)...)
	routerGroup.DELETE("/category/:id", append(middlewareList, cat.Delete)...)
	routerGroup.GET("/category/search", append(middlewareList, cat.GetByFilterMany)...)
	routerGroup.GET("/category/filter", append(middlewareList, cat.GetByFilterOne)...)
}

func (cat *CategoryHandlerHttp) Create(c *gin.Context) {
	var category entity.Category
	err := c.BindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("the name is required")})
		c.Abort()
		return
	}

	catResponse, err := entity.NewCategory(&category.Name, &category.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	response, err := cat.Service.Create(catResponse)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (cat *CategoryHandlerHttp) Get(c *gin.Context) {
	response, err := cat.Service.Get()
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

func (cat *CategoryHandlerHttp) GetById(c *gin.Context) {

	categoryId := c.Param("id")
	if categoryId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	response, err := cat.Service.GetById(&categoryId)
	if err != nil {
		log.Println(err)
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

func (cat *CategoryHandlerHttp) Update(c *gin.Context) {
	categoryId := c.Param("id")
	if categoryId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	err := utils.ValidateUUID(&categoryId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var category entity.CategoryResponse
	if err := c.ShouldBindJSON(&category); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is required"})
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if category.ID != categoryId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id in body must be equal to id in path"})
		c.Abort()
		return
	}

	data, err := cat.Service.Update(&category)
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

func (cat *CategoryHandlerHttp) Delete(c *gin.Context) {
	categoryId := c.Param("id")
	if categoryId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	err := cat.Service.Delete(&categoryId)

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

func (cat *CategoryHandlerHttp) GetByFilterMany(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key and value are required"})
		c.Abort()
		return
	}

	response, err := cat.Service.GetByFilterMany(key, &value)
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

func (cat *CategoryHandlerHttp) GetByFilterOne(c *gin.Context) {

	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key and value are required"})
		c.Abort()
		return
	}

	response, err := cat.Service.GetByFilterOne(key, &value)
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
