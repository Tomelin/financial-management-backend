package web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	middleware "github.com/synera-br/financial-management/src/backend/internal/infra/handler/middleware/authorization"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type WalletHandlerHttpInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetWalletByIdAndUserID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByFilterMany(c *gin.Context)
	GetByFilterOne(c *gin.Context)
}

type WalletHandlerHttp struct {
	Service entity.IWallet
	User    entity.IUser
}

func NewWalletHandlerHttp(svc *entity.IWallet, user *entity.IUser, routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) WalletHandlerHttpInterface {

	lab := &WalletHandlerHttp{
		Service: *svc,
		User:    *user,
	}

	lab.handlers(routerGroup, middleware...)

	return lab

}

func (c *WalletHandlerHttp) handlers(routerGroup *gin.RouterGroup, middleware ...func(c *gin.Context)) {
	middlewareList := make([]gin.HandlerFunc, len(middleware))
	for i, mw := range middleware {
		middlewareList[i] = mw
	}

	routerGroup.POST("/wallet", append(middlewareList, c.Create)...)
	routerGroup.GET("/wallet/:id", append(middlewareList, c.GetWalletByIdAndUserID)...)
	routerGroup.GET("/wallet/search", append(middlewareList, c.GetByFilterMany)...)
	routerGroup.GET("/wallet/filter", append(middlewareList, c.GetByFilterOne)...)
	routerGroup.GET("/wallet", append(middlewareList, c.Get)...)
	routerGroup.PUT("/wallet/:id", append(middlewareList, c.Update)...)
	routerGroup.DELETE("/wallet/:id", append(middlewareList, c.Delete)...)
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
// @Router      /wallet [post]
func (obj *WalletHandlerHttp) Create(c *gin.Context) {

	var wallet entity.WalletResponse
	ctx := context.Background()
	if err := c.BindJSON(&wallet); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, entity.Error("body is required", "wallet", "Create", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest))
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, entity.Error(err.Error(), "wallet", "Create", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest))
		c.Abort()
		return
	}

	data, mErr := entity.NewWallet(&wallet)
	if mErr != nil {
		c.JSON(http.StatusBadGateway, mErr)
		c.Abort()
		return
	}

	authId, err := middleware.GetEmailFromToken(c)
	if mErr != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	user, err := obj.User.GetByEmail(ctx, authId)
	if mErr != nil {
		c.JSON(500, gin.H{"error2": err.Error()})
		c.Abort()
		return
	}

	if authId == nil || *authId == "" || user == nil {
		c.JSON(401, gin.H{"error": "token authorization is required"})
		c.Abort()
		return
	}

	if user.Email != *authId {
		c.JSON(401, gin.H{"error": user.Email + " is not equal to " + *authId})
		c.Abort()
		return
	}

	result, mErr := obj.Service.Create(context.Background(), &user.ID, data)
	if mErr != nil {
		c.JSON(http.StatusInternalServerError, mErr)
		c.Abort()
		return
	}
	c.JSON(http.StatusAccepted, result)
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
// @Router      /wallet [get]
func (obj *WalletHandlerHttp) Get(c *gin.Context) {

	authId, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(500, gin.H{"error": entity.Error(err.Error(), "wallet", "Get", entity.ApplicationLayerHandler, entity.ResponseCodeUnauthorized)})
		c.Abort()
		return
	}

	ctx := context.Background()
	user, err := obj.User.GetByEmail(ctx, authId)
	if err != nil {
		c.JSON(500, gin.H{"error": entity.Error(err.Error(), "wallet", "Get", entity.ApplicationLayerHandler, entity.ResponseCodeUnauthorized)})
		c.Abort()
		return
	}

	if authId == nil || *authId == "" || user == nil {
		c.JSON(401, gin.H{"error": "token authorization is required"})
		c.Abort()
		return
	}

	if user.Email != *authId {
		c.JSON(401, gin.H{"error": user.Email + " is not equal to " + *authId})
		c.Abort()
		return
	}

	response, mErr := obj.Service.Get(context.Background(), &user.ID)
	if mErr != nil {
		c.JSON(500, gin.H{"error": mErr})
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
// @Router      /wallet/{id} [get]
func (obj *WalletHandlerHttp) GetWalletByIdAndUserID(c *gin.Context) {

	walletId := c.Param("id")
	if walletId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		c.Abort()
		return
	}

	email, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	response, mErr := obj.Service.GetWalletByIdAndUserID(context.Background(), email, &walletId)
	if mErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": mErr})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}

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
// @Router      /wallet/search [get]
func (obj *WalletHandlerHttp) GetByFilterMany(c *gin.Context) {

	key := c.Query("key")
	value := c.Query("value")
	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("key and value are required", "wallet", "GetByFilterMany", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	query := []entity.QueryDB{
		{
			Key:       key,
			Value:     value,
			Condition: string(entity.QueryFirebaseEqual),
		},
	}

	email, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "wallet", "GetByFilterMany", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	response, mErr := obj.Service.GetByFilterMany(context.Background(), email, query)
	if mErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": mErr})
		c.Abort()
		return
	}

	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"messagem": entity.Error("not found", "wallet", "GetByFilterMany", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
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
// @Router      /wallet/search [get]
func (obj *WalletHandlerHttp) GetByFilterOne(c *gin.Context) {

	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("key and value are required", "wallet", "GetByFilterOne", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	query := []entity.QueryDB{
		{
			Key:       key,
			Value:     value,
			Condition: string(entity.QueryFirebaseEqual),
		},
	}

	email, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	response, mErr := obj.Service.GetByFilterOne(context.Background(), email, query)
	if mErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": mErr})
		c.Abort()
		return
	}

	if response == nil || response.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": entity.Error("not found", "wallet", "GetByFilterOne", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}

func (obj *WalletHandlerHttp) Update(c *gin.Context) {

	walletId := c.Param("id")
	if walletId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("id is required", "wallet", "Update", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	err := utils.ValidateUUID(&walletId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	var wallet entity.WalletResponse
	if err := c.ShouldBindJSON(&wallet); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("body is required", "wallet", "Update", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		return
	}

	if wallet.ID != walletId {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("id in body must be equal to id in path", "wallet", "Update", entity.ApplicationLayerHandler, entity.ResponseCodeUnauthorized)})
		c.Abort()
		return
	}

	email, err := middleware.GetEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	data, mErr := obj.Service.Update(context.Background(), email, &wallet)
	if mErr != nil {
		if mErr.Code == entity.ResponseCodeNotFound {
			c.JSON(int(mErr.Code), gin.H{"error": mErr})
		} else {
			c.JSON(int(mErr.Code), gin.H{"error": mErr})
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, data)
}

func (obj *WalletHandlerHttp) Delete(c *gin.Context) {

	walletId := c.Param("id")
	if walletId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": entity.Error("id is required", "wallet", "Delete", entity.ApplicationLayerHandler, entity.ResponseCodeBadRequest)})
		c.Abort()
		return
	}

	authId := obj.getAccessToken(c)
	err := obj.Service.Delete(context.Background(), authId, &walletId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "deleted"})
}

func (obj *WalletHandlerHttp) getAccessToken(c *gin.Context) *string {

	authId := c.GetHeader("Authorization")
	if authId == "" {
		return nil
	}
	claims, err := middleware.OpenTokenJWT(&authId)
	if err != nil {
		return nil
	}

	return &claims.Email
}
