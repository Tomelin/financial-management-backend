package http_server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

type Certificate struct {
	Key string `mapstructure:"key"`
	Crt string `mapstructure:"crt"`
}

type RestAPIConfig struct {
	Port           string `mapstructure:"port"`
	SSLEnabled     bool   `mapstructure:"ssl_enabled"`
	Host           string `mapstructure:"host"`
	Version        string `mapstructure:"version"`
	Name           string `mapstructure:"name"`
	CertificateCrt string `mapstructure:"certificate_crt"`
	CertificateKey string `mapstructure:"certificate_key"`
	Token          string `mapstructure:"token"`
}

type RestAPI struct {
	Config *RestAPIConfig
	Route  *gin.Engine
	*gin.RouterGroup
}

func NewRestApi(fields any, s *sessions.CookieStore) (*RestAPI, error) {

	b, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	var rest *RestAPIConfig
	err = json.Unmarshal(b, &rest)
	if err != nil {
		return nil, err
	}

	r, g := newRestAPI(rest, s)

	return &RestAPI{
		Config:      rest,
		Route:       r,
		RouterGroup: g,
	}, nil
}

var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "cloud-collector-resources3",
	Description:      "",
	InfoInstanceName: "swagger",
	// SwaggerTemplate:  docTemplate,
	// LeftDelim:  "{{",
	// RightDelim: "}}",
}

func newRestAPI(config *RestAPIConfig, s *sessions.CookieStore) (*gin.Engine, *gin.RouterGroup) {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("X-Authorization", *s))

	router.UseH2C = true

	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.2", "10.0.0.0/8"})

	routerGroupPath := fmt.Sprintf("/%s", config.Name)
	routerPath = router.Group(routerGroupPath)

	router.GET("/metrics", prometheusHandler())

	// Set swagger
	routerPath.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))
	routerPath.GET("/docs/swagger", func(c *gin.Context) {
		c.Redirect(301, fmt.Sprintf("%s/docs/swagger/index.html", routerGroupPath))
	})
	routerPath.GET("/docs", func(c *gin.Context) {

		c.Redirect(301, fmt.Sprintf("%s/docs/swagger/index.html", routerGroupPath))
	})
	routerPath.GET("/", func(c *gin.Context) {
		c.Redirect(301, fmt.Sprintf("%s/docs/swagger/index.html", routerGroupPath))
	})

	router.Use(setHeader)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(cors.New(corsConfig))

	return router, routerPath

}
