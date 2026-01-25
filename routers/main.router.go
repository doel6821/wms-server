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
		// AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge: 12 * time.Hour,
	}))

	r.Gin.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

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
		v1.POST("/forgot", r.Controller.V1().ForgotPassword)
		v1.POST("/change-password", r.Controller.V1().ChangePassword)
		v1.POST("/register-user", r.Controller.V1().CekToken, r.Controller.V1().RegisterUser)
		v1.POST("/update-user/:id", r.Controller.V1().CekToken, r.Controller.V1().UpdateUser)
		v1.POST("/register-tenant", r.Controller.V1().RegisterTenant)
		v1.GET("/user/all", r.Controller.V1().CekToken, r.Controller.V1().ListUser)

		dashboard := v1.Group("/dashboard")	
		dashboard.GET("/", r.Controller.V1().CekToken, r.Controller.V1().GetDataDashboard)
		

		customer := v1.Group("/customer")
		customer.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListCustomer)
		customer.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetCustomerById)
		customer.POST("/", r.Controller.V1().CekToken, r.Controller.V1().SaveCustomer)
		customer.DELETE("/:id", r.Controller.V1().CekToken, r.Controller.V1().DeleteCustomerById)

		supplier := v1.Group("/supplier")
		supplier.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListSupplier)
		supplier.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetSupplierById)
		supplier.POST("/", r.Controller.V1().CekToken, r.Controller.V1().SaveSupplier)
		supplier.DELETE("/:id", r.Controller.V1().CekToken, r.Controller.V1().DeleteSupplierById)

		product := v1.Group("/product")
		product.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListProduct)
		product.GET("/total", r.Controller.V1().CekToken, r.Controller.V1().ProductTotal)
		product.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetProductById)
		product.POST("/", r.Controller.V1().CekToken, r.Controller.V1().SaveProduct)
		product.DELETE("/:id", r.Controller.V1().CekToken, r.Controller.V1().DeleteProductById)

		location := v1.Group("/location")
		location.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListLocation)
		// location.GET("/:code", r.Controller.V1().CekToken, r.Controller.V1().GetLocationByCode)
		location.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetLocationById)
		location.POST("/", r.Controller.V1().CekToken, r.Controller.V1().SaveLocation)
		location.DELETE("/:id", r.Controller.V1().CekToken, r.Controller.V1().DeleteLocationById)

		config := v1.Group("/config")
		config.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListConfiguration)
		// config.GET("/:code", r.Controller.V1().CekToken, r.Controller.V1().GetConfigurationByName)
		config.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetConfigurationById)
		config.POST("/", r.Controller.V1().CekToken, r.Controller.V1().SaveConfiguration)
		config.DELETE("/:id", r.Controller.V1().CekToken, r.Controller.V1().DeleteConfigurationById)

		sales := v1.Group("/sales-order")
		sales.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListSalesOrder)
		sales.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetSalesOrderById)
		sales.POST("", r.Controller.V1().CekToken, r.Controller.V1().SaveSalesOrder)

		salesOrderItem := v1.Group("/sales-order-item")
		salesOrderItem.GET("", r.Controller.V1().CekToken, r.Controller.V1().GetSalesOrderItemByProductId)

		packing := v1.Group("/packing-order")
		packing.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListPackingOrder)
		packing.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetPackingOrderDetailById)
		packing.POST("/", r.Controller.V1().CekToken, r.Controller.V1().SavePackingOrder)

		packingOrderItem := v1.Group("/packing-order-item")
		packingOrderItem.GET("/", r.Controller.V1().CekToken, r.Controller.V1().GetPackingOrderItemByProductId)

		invoice := v1.Group("/invoice")
		invoice.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListInvoice)
		invoice.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetInvoiceDetailById)
		invoice.POST("/", r.Controller.V1().CekToken, r.Controller.V1().SaveInvoice)
		invoice.PUT("/", r.Controller.V1().CekToken, r.Controller.V1().UpdateInvoiceOrder)

		finance := v1.Group("/finance")
		finance.GET("/payment/all", r.Controller.V1().CekToken, r.Controller.V1().ListAccountPayable)
		finance.GET("/receive/all", r.Controller.V1().CekToken, r.Controller.V1().ListAccountReceivable)
		finance.POST("/payment", r.Controller.V1().CekToken, r.Controller.V1().SaveAccountPayable)
		finance.POST("/receive", r.Controller.V1().CekToken, r.Controller.V1().SaveAccountReceivable)

		purchase := v1.Group("/purchase")
		purchase.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListPurchaseOrder)
		purchase.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetPurchaseOrderById)
		purchase.GET("/recomendation/:id", r.Controller.V1().CekToken, r.Controller.V1().PurchaseOrderRecomendation)
		purchase.POST("/", r.Controller.V1().CekToken, r.Controller.V1().SavePurchaseOrder)

		receive := v1.Group("/receive-order")
		receive.GET("/available", r.Controller.V1().CekToken, r.Controller.V1().GetAvailableReceive)
		receive.GET("/all", r.Controller.V1().CekToken, r.Controller.V1().ListReceiveOrder)
		receive.GET("/:id", r.Controller.V1().CekToken, r.Controller.V1().GetReceiveOrderDetailById)
		receive.GET("/product/:productId", r.Controller.V1().CekToken, r.Controller.V1().GetReceiveOrderDetailByProductId)
		receive.POST("/", r.Controller.V1().CekToken, r.Controller.V1().SaveReceiveOrder)
		receive.POST("/stocked", r.Controller.V1().CekToken, r.Controller.V1().SaveStocked)
		receive.PUT("/", r.Controller.V1().CekToken, r.Controller.V1().UpdateReceiveOrder)
		


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
