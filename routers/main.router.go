package routers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"wms-server/constants"
	"wms-server/controllers"
	"wms-server/helpers"
	"strings"
	"time"

	"wms-server/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterInterface Abstract Class
type (
	RouterInterface interface {
		StartServer() error
		routerControllers()
	}

	// Router Actual Class implementation
	Router struct {
		address     string
		port        string
		Path        map[string]string
		Gin         *gin.Engine
		Controller  controllers.Controller
		Logs        *logs.Logger
		maskingData string
	}
)

// InitializeRouter return sturct that will implement the abs. class
func InitializeRouter(ctrl controllers.Controller, l *logs.Logger) RouterInterface {
	gin.SetMode(helpers.GetEnv("ROUTER_SETMODE"))
	return &Router{
		address: helpers.GetEnv("ROUTER_SERVER_ADDRESS"),
		port:    helpers.GetEnv("ROUTER_PORT"),
		Path: map[string]string{
			"PATH_VERSION": helpers.GetEnv("PATH_VERSION"),
			"PATH_MAIN":    helpers.GetEnv("PATH_MAIN"),
		},
		Gin:         gin.New(),
		Controller:  ctrl,
		Logs:        l,
		maskingData: helpers.GetEnv("MASKING_DATA"),
	}
}

// StartServer Start the server by initialize end-point & create the server
func (r *Router) StartServer() error {
	r.Logs.AddHook(helpers.NewTraceIDHook("-"))
	r.Logs.WithContext(context.Background()).Info("Starting Server on ", r.address+r.port)
	r.routerControllers()

	err := helpers.GinServerUp(r.address+r.port, r.Gin)
	if err != nil {
		r.Logs.WithContext(context.Background()).Error("[GinServerUp]Error: ", err)
		return err
	}

	return nil
}

// routerControllers ...
func (r *Router) routerControllers() {

	docs.SwaggerInfo.Title = helpers.GetEnv("SWAG_TITLE")
	docs.SwaggerInfo.Description = helpers.GetEnv("SWAG_DESC")
	docs.SwaggerInfo.Version = helpers.GetEnv("SWAG_VERSION")
	docs.SwaggerInfo.Host = helpers.GetEnv("SWAGGER_HOST")
	docs.SwaggerInfo.BasePath = r.Path["PATH_MAIN"]
	docs.SwaggerInfo.Schemes = strings.Split(helpers.GetEnv("SWAG_SCHEMES"), ",")

	r.Gin.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "PATCH", "PUT", "POST", "HEAD", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization", "access-control-allow-origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		//AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge: 12 * time.Hour,
	}))

	r.Gin.Use(r.GinLogger())
	r.Gin.Use(gin.Recovery())

	url := ginSwagger.URL("http://" + docs.SwaggerInfo.Host + "/swagger/doc.json") // The url pointing to API definition
	// }

	r.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	main := r.Gin.Group(r.Path["PATH_MAIN"])
	main.GET("/healthcheck", r.Controller.V1().HealthCheck)

	// BEGIN __INCLUDE_TEMPLATE__
	{
		// version 1
		v1 := main.Group("/v1")
		v1.POST("/login", r.Controller.V1().Login)
		v1.POST("/register-user", r.Controller.V1().RegisterUser)
		v1.POST("/register-tenant", r.Controller.V1().RegisterTenant)

		// customer := v1.Group("/customer")
		// customer.GET("/all", r.Controller.V1().BindQueryParam)
		// customer.GET("/:id", r.Controller.V1().GetData)
		// customer.POST("/", r.Controller.V1().BindBody)
		// customer.DELETE("/:id", r.Controller.V1().BindHeader)

		// product := v1.Group("/product")
		// product.GET("/all", r.Controller.V1().)
		// product.GET("/:id", r.Controller.V1().)
		// product.POST("/", r.Controller.V1().)
		// product.DELETE("/:id", r.Controller.V1().)

		// sales := v1.Group("/sales")
		// sales.GET("/all", r.Controller.V1().)
		// sales.GET("/:id", r.Controller.V1().)
		// sales.POST("/", r.Controller.V1().)
		// sales.DELETE("/:id", r.Controller.V1().)

		// purchase := v1.Group("/purchase")
		// purchase.GET("/all", r.Controller.V1().)
		// purchase.GET("/:id", r.Controller.V1().)
		// purchase.POST("/", r.Controller.V1().)
		// purchase.DELETE("/:id", r.Controller.V1().)

		// receive := v1.Group("/receive")
		// receive.GET("/all", r.Controller.V1().)
		// receive.GET("/:id", r.Controller.V1().)
		// receive.POST("/", r.Controller.V1().)
		// receive.DELETE("/:id", r.Controller.V1().)

		// report := v1.Group("/report")
		// report.GET("/sales", r.Controller.V1().)
		// report.GET("/purchase", r.Controller.V1().)


	}

	// Invalid URL
	r.Gin.NoRoute(func(c *gin.Context) {
		r.Logs.WithContext(c).WithFields(logs.Fields{"URL": r.address + r.port, "Method": c.Request.Method, "Path": c.Request.URL.Path}).Error("[NoRoute]Invalid")
		c.JSON(404, gin.H{"code": "404", "message": "Page not found"})
	})

	r.Logs.WithContext(context.Background()).Info("initializing Router done")
}

// ResponseBodyWriter ...
type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r ResponseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// GinLogger ...
func (r *Router) GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		timeDuration, errorDuration := time.ParseDuration(helpers.GetEnv("MAX_TIMEOUT"))
		if errorDuration != nil {
			timeDuration = 2 * time.Second
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeDuration)

		c.Request = c.Request.WithContext(ctx)
		w := &ResponseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		start := time.Now()
		traceID := helpers.GetTraceID(c)
		c.Set(constants.TRANSACTION_ID, traceID)

		var req, res interface{}
		bodyCopy := new(bytes.Buffer)
		io.Copy(bodyCopy, c.Request.Body)
		bodyData := bodyCopy.Bytes()
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyData))

		// Decode request body
		json.NewDecoder(bodyCopy).Decode(&req)

		// Channel untuk menangkap sinyal selesai
		done := make(chan struct{})

		go func() {
			// Melanjutkan ke handler berikutnya
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// Request selesai diproses
		case <-ctx.Done():
			// Timeout terjadi
			c.JSON(http.StatusGatewayTimeout, helpers.GetMetaResponse(constants.RC_TIMEOUT))
			c.Abort()
			// return
		}

		// Decode response body
		json.NewDecoder(w.body).Decode(&res)

		// Filter sensitive information
		filteredReq := r.filterSensitiveData(&req)
		filteredRes := r.filterSensitiveData(&res)

		defer func() {
			r.Logs.WithContext(c).WithFields(helpers.GettingResponseLog(c, filteredReq, filteredRes, time.Since(start))).Infof("API LOG")
			cancel()
		}()
	}
}

// filterSensitiveData removes sensitive fields from the data.
func (r *Router) filterSensitiveData(data interface{}) interface{} {
	maskingData := strings.Split(r.maskingData, ",")

	if data == nil {
		return nil
	}

	// Convert to map for easier filtering
	if dataMap, ok := (data).(map[string]interface{}); ok {

		for _, m := range maskingData {
			if _, exists := dataMap[m]; exists {
				dataMap[m] = "***" // Mask data
			}
		}

		return dataMap
	}

	return data // Return original if not a map
}
