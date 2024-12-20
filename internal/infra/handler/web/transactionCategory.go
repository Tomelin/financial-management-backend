package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	middleware "github.com/synera-br/financial-management/src/backend/internal/infra/handler/middleware/authorization"
	"github.com/synera-br/financial-management/src/backend/pkg/observability"
)

type TransactionCategoryHandlerHttpInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetById(c *gin.Context)
	GetByFilterMany(c *gin.Context)
	GetByFilterOne(c *gin.Context)
	// GetWalletByIdAndUserID(c *gin.Context)
	// Update(c *gin.Context)
	// Delete(c *gin.Context)
}

type TransactionCategoryHandlerHttp struct {
	Service entity.ITransactionCategory
	Wallet  entity.IWallet
	Trace   *observability.Tracer
}

func NewTransactionCategoryHandlerHttp(trace *observability.Tracer, svc *entity.ITransactionCategory, wallet *entity.IWallet, routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) TransactionCategoryHandlerHttpInterface {

	lab := &TransactionCategoryHandlerHttp{
		Service: *svc,
		Wallet:  *wallet,
		Trace:   trace,
	}

	lab.handlers(routerGroup, middleware...)

	return lab

}

func (c *TransactionCategoryHandlerHttp) handlers(routerGroup *gin.RouterGroup, middlewares ...func(g *gin.Context)) {

	middlewareList := make([]gin.HandlerFunc, len(middlewares))
	for i, mw := range middlewares {

		middlewareList[i] = mw
	}

	routerGroup.POST("/category", append(middlewareList, c.Create)...)
	routerGroup.GET("/category", append(middlewareList, c.Get)...)
	routerGroup.GET("/category/:id", append(middlewareList, c.GetById)...)
	routerGroup.GET("/category/search", append(middlewareList, c.GetByFilterMany)...)
	routerGroup.GET("/category/filter", append(middlewareList, c.GetByFilterOne)...)
	// routerGroup.GET("/category/:id", append(middlewareList, c.GetWalletByIdAndUserID)...)
	// routerGroup.PUT("/category/:id", append(middlewareList, c.Update)...)
	// routerGroup.DELETE("/category/:id", append(middlewareList, c.Delete)...)
}

// CreateWalletResponse    godoc
// @Summary     create a new lab destroy
// @Tags        Wallet
// @Accept       json
// @Produce     json
// @Description create a new lab destroy
// @Success     200 {object} entity.WalletResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /category [post]
func (obj *TransactionCategoryHandlerHttp) Create(c *gin.Context) {
	ctx, span := obj.Trace.Trace.Start(context.Background(), fmt.Sprintf("%s.Create", string(entity.ApplicationLayerHandler)))
	defer span.End()

	var category entity.TransactionCategory
	if err := c.BindJSON(&category); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, entity.Error("body is required", "transaction_category", "Create", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest))
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, entity.Error(err.Error(), "wallet", "transaction_category", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest))
		c.Abort()
		return
	}

	data, mErr := entity.NewTransactionCategory(&category)
	if mErr != nil {
		e := entity.ResponseMessage(mErr.Code)
		c.JSON(int(entity.ResponseMessageToCode(e)), mErr)
		c.Abort()
		return
	}

	email, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, entity.Error(err.Error(), "create", "transaction_category", entity.ApplicationLayerHandler, entity.ResponseCodeUnauthorized))
		c.Abort()
		return
	}

	result, mErr := obj.Service.Create(ctx, email, data)
	if mErr != nil {
		c.JSON(int(mErr.Code), mErr)
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusAccepted, result)
}

// GetAllWalletResponse    godoc
// @Summary     get all lab destroy
// @Tags        Wallet
// @Accept       json
// @Produce     json
// @Description get all lab destroy
// @Success     200 {object} []entity.WalletResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /category [get]
func (obj *TransactionCategoryHandlerHttp) Get(c *gin.Context) {
	ctx, span := obj.Trace.Trace.Start(context.Background(), fmt.Sprintf("%s.Get", string(entity.ApplicationLayerHandler)))
	defer span.End()

	email, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, entity.Error(err.Error(), "get", "transaction_category", entity.ApplicationLayerHandler, entity.ResponseCodeUnauthorized))
		c.Abort()
		return
	}

	walletID := c.Query("wallet_id")

	response, mErr := obj.Service.Get(ctx, email, &walletID)
	if mErr != nil {
		c.JSON(int(mErr.Code), mErr)
		c.Abort()
		return
	}

	if len(response) == 0 {
		c.JSON(http.StatusNotFound, entity.Error("category not found", "get", "transaction_category", entity.ApplicationLayerHandler, entity.ResponseCodeNotFound))
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetAllWalletResponse    godoc
// @Summary     get all lab destroy
// @Tags        Wallet
// @Accept       json
// @Produce     json
// @Description get all lab destroy
// @Success     200 {object} entity.WalletResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /category/{id} [get]
func (obj *TransactionCategoryHandlerHttp) GetById(c *gin.Context) {
	ctx, span := obj.Trace.Trace.Start(context.Background(), fmt.Sprintf("%s.GetById", string(entity.ApplicationLayerHandler)))
	defer span.End()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, entity.Error("id is required", "get", "transaction_category", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest))
		c.Abort()
		return
	}

	email, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, entity.Error(err.Error(), "get", "transaction_category", entity.ApplicationLayerHandler, entity.ResponseCodeUnauthorized))
		c.Abort()
		return
	}

	response, mErr := obj.Service.GetById(ctx, email, &id)
	if mErr != nil {
		if mErr.Code == entity.ResponseCodeNotFound {
			c.JSON(int(mErr.Code), "category not found")
		} else {
			c.JSON(int(mErr.Code), mErr)
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}

// 	c.JSON(http.StatusOK, response)
// }

// GetWalletByIdAndUserID      godoc
// @Summary     get a lab destroy by ID
// @Tags        Wallet
// @Accept      json
// @Produce     json
// @Param       id path string true "found"
// @Description get a lab destroy by ID
// @Success     200 {object} entity.WalletResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /category/{id} [get]
// func (obj *TransactionCategoryHandlerHttp) GetWalletByIdAndUserID(c *gin.Context) {

// 	walletId := c.Param("id")
// 	if walletId == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
// 		c.Abort()
// 		return
// 	}

// 	email, err := middleware.GetEmailFromToken(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		c.Abort()
// 		return
// 	}

// 	response, mErr := obj.Service.GetWalletByIdAndUserID(context.Background(), email, &walletId)
// 	if mErr != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": mErr})
// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// GetByFilterMany  godoc
// @Summary     get a lab destroy by email or project
// @Tags        Wallet
// @Accept      json
// @Produce     json
// @Param       email query string false "email@domain.com"
// @Param       project query string false "project_name"
// @Param       available query boolean false "string default" default(false)
// @Description get a lab destroy by email or project
// @Success     200 {object} []entity.WalletResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /category/search [get]
func (obj *TransactionCategoryHandlerHttp) GetByFilterMany(c *gin.Context) {
	ctx, span := obj.Trace.Trace.Start(context.Background(), fmt.Sprintf("%s.GetByFilterMany", string(entity.ApplicationLayerHandler)))
	defer span.End()

	key := c.Query("key")
	value := c.Query("value")
	condition := entity.QueryFirebaseString(c.Query("condition"))

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("key and value are required", "transactionCategory", "GetByFilterMany", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	query := []entity.QueryDB{
		{
			Key:       key,
			Value:     value,
			Condition: condition,
		},
	}

	email, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "transactionCategory", "GetByFilterMany", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	response, mErr := obj.Service.GetByFilterMany(ctx, email, query)
	if mErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": mErr})
		c.Abort()
		return
	}

	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"messagem": entity.Error("not found", "transactionCategory", "GetByFilterMany", entity.ApplicationLayerHandler, entity.ResponseCodeNotFound)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByFilterOne  godoc
// @Summary     get a lab destroy by email or project
// @Tags        Wallet
// @Accept      json
// @Produce     json
// @Param       email query string false "email@domain.com"
// @Param       project query string false "project_name"
// @Param       available query boolean false "string default" default(false)
// @Description get a lab destroy by email or project
// @Success     200 {object} []entity.WalletResponse
// @Failure     404 {object} string
// @Failure     500 {object} string
// @Router      /category/search [get]
func (obj *TransactionCategoryHandlerHttp) GetByFilterOne(c *gin.Context) {
	ctx, span := obj.Trace.Trace.Start(context.Background(), fmt.Sprintf("%s.GetByFilterOne", string(entity.ApplicationLayerHandler)))
	defer span.End()

	key := c.Query("key")
	value := c.Query("value")
	condition := entity.QueryFirebaseString(c.Query("condition"))

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("key and value are required", "transactionCategory", "GetByFilterOne", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	query := []entity.QueryDB{
		{
			Key:       key,
			Value:     value,
			Condition: condition,
		},
	}

	email, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "transactionCategory", "GetByFilterOne", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	response, mErr := obj.Service.GetByFilterOne(ctx, email, query)
	if mErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": mErr})
		c.Abort()
		return
	}

	if response == nil || response.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": entity.Error("not found", "transactionCategory", "GetByFilterOne", entity.ApplicationLayerHandler, entity.ResponseCodeNotFound)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}

// func (obj *TransactionCategoryHandlerHttp) Update(c *gin.Context) {

// 	walletId := c.Param("id")
// 	if walletId == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("id is required", "wallet", "Update", entity.ApplicationLayerHandler, "")})
// 		c.Abort()
// 		return
// 	}

// 	err := utils.ValidateUUID(&walletId)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerHandler, "validate uuid")})
// 		c.Abort()
// 		return
// 	}

// 	var wallet entity.WalletResponse
// 	if err := c.ShouldBindJSON(&wallet); err != nil {
// 		if err.Error() == "EOF" {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("body is required", "wallet", "Update", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
// 			c.Abort()
// 			return
// 		}
// 		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
// 		return
// 	}

// 	if wallet.ID != walletId {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("id in body must be equal to id in path", "wallet", "Update", entity.ApplicationLayerHandler, "validate id")})
// 		c.Abort()
// 		return
// 	}

// 	email, err := middleware.GetEmailFromToken(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerHandler, "")})
// 		c.Abort()
// 		return
// 	}

// 	data, mErr := obj.Service.Update(context.Background(), email, &wallet)
// 	if err != nil {
// 		if err.Error() == "not found" {
// 			c.JSON(http.StatusNotFound, gin.H{"error3": mErr})
// 		} else {
// 			c.JSON(http.StatusBadRequest, gin.H{"error4": mErr})
// 		}
// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusOK, data)
// }

// func (obj *TransactionCategoryHandlerHttp) Delete(c *gin.Context) {

// 	walletId := c.Param("id")
// 	if walletId == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("id is required", "wallet", "Delete", entity.ApplicationLayerHandler, "")})
// 		c.Abort()
// 		return
// 	}

// 	authId := obj.getAccessToken(c)
// 	err := obj.Service.Delete(context.Background(), authId, &walletId)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err})
// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusNoContent, gin.H{"message": "deleted"})
// }

// func (obj *TransactionCategoryHandlerHttp) getAccessToken(c *gin.Context) *string {

// 	authId := c.GetHeader("Authorization")
// 	if authId == "" {
// 		return nil
// 	}
// 	claims, err := middleware.OpenTokenJWT(&authId)
// 	if err != nil {
// 		return nil
// 	}

// 	return &claims.Email
// }
