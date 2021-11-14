package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4/middleware"

	"log"
	"mini_pos/config"
	"net/url"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"


	_authHttpDelivery "mini_pos/api/v1/auth/delivery/http"
	_authRepo "mini_pos/api/v1/auth/repository/mysql"
	_authUcase "mini_pos/api/v1/auth/usecase"

	_productHttpDelivery "mini_pos/api/v1/product/delivery/http"
	_productRepo "mini_pos/api/v1/product/repository/mysql"
	_productUcase "mini_pos/api/v1/product/usecase"

	_categoryRepo "mini_pos/api/v1/category/repository/mysql"
	_outletRepo "mini_pos/api/v1/outlet/repository/mysql"

	_supplierHttpDelivery "mini_pos/api/v1/supplier/delivery/http"
	_supplierRepo "mini_pos/api/v1/supplier/repository/mysql"
	_supplierUcase "mini_pos/api/v1/supplier/usecase"

	_customerHttpDelivery "mini_pos/api/v1/customer/delivery/http"
	_customerRepo "mini_pos/api/v1/customer/repository/mysql"
	_customerUcase "mini_pos/api/v1/customer/usecase"

	_purchaseHttpDelivery "mini_pos/api/v1/purchase/delivery/http"
	_purchaseRepo "mini_pos/api/v1/purchase/repository/mysql"
	_purchaseUcase "mini_pos/api/v1/purchase/usecase"

	_product_outletHttpDelivery "mini_pos/api/v1/product_outlets/delivery/http"
	_product_outletRepo "mini_pos/api/v1/product_outlets/repository/mysql"
	_product_outletUcase "mini_pos/api/v1/product_outlets/usecase"

	_orderRepo "mini_pos/api/v1/order/repository/mysql"
	_orderItemHttpDelivery "mini_pos/api/v1/order_item/delivery/http"
	_orderItemRepo "mini_pos/api/v1/order_item/repository/mysql"
	_orderItemUcase "mini_pos/api/v1/order_item/usecase"

	_paymentHttpDelivery "mini_pos/api/v1/payment/delivery/http"
	_paymentRepo "mini_pos/api/v1/payment/repository/mysql"
	_paymentUcase "mini_pos/api/v1/payment/usecase"

	_mediaHttpDelivery "mini_pos/api/v1/media/delivery/http"
	_mediaUcase "mini_pos/api/v1/media/usecase"
)

func dbConnection(cfg *config.DBConfig) *sql.DB{
	dbHost := cfg.Host
	dbPort := cfg.Port
	dbUser := cfg.Username
	dbPass := cfg.Password
	dbName := cfg.Name
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}else{
		//fmt.Println("connect", dsn)
	}


	return dbConn
}

func main() {
	cfg := config.InitConfig()

	connection := dbConnection(&cfg.DB)
	err := connection.Ping()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{ "http://localhost:9090"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	timeoutContext := time.Duration(cfg.Server.WriteTimeout) * time.Second

	//Auth
	authRepo := _authRepo.NewMysqlAuthRepository(connection)
	auth := _authUcase.NewAuthUseCase(authRepo, timeoutContext)
	_authHttpDelivery.NewAuthHandler(e, auth, cfg)

	supplierRepo := _supplierRepo.NewMysqlSupplierRepository(connection)
	supplier := _supplierUcase.NewSupplierUsecase(supplierRepo,  timeoutContext)
	_supplierHttpDelivery.NewSupplierHandler(e, supplier, cfg)

	categoryRepo := _categoryRepo.NewMysqlCategoryRepository(connection)
	outletRepo := _outletRepo.NewMysqlOutletRepository(connection)
	productRepo := _productRepo.NewMysqlProductRepository(connection)
	product := _productUcase.NewproductUsecase(productRepo, categoryRepo, timeoutContext)
	_productHttpDelivery.NewProductHandler(e, product, cfg)

	customerRepo := _customerRepo.NewMysqlCustomerRepository(connection)
	customer := _customerUcase.NewCustomerUsecase(customerRepo,  timeoutContext)
	_customerHttpDelivery.NewCustomerHandler(e, customer, cfg)

	purchaseRepo := _purchaseRepo.NewMysqlPurchaseRepository(connection)
	purchase := _purchaseUcase.NewPurchaseUsecase(productRepo, purchaseRepo, categoryRepo, outletRepo, supplierRepo, timeoutContext)
	_purchaseHttpDelivery.NewPurchaseHandler(e, purchase, cfg)

	productOutletRepo := _product_outletRepo.NewMysqlProductOutletsRepository(connection)
	productOutlet := _product_outletUcase.NewProductOutletsUsecase(productOutletRepo, productRepo, outletRepo, supplierRepo, timeoutContext)
	_product_outletHttpDelivery.NewProductOutletsHandler(e, productOutlet, cfg)

	orderRepo := _orderRepo.NewMysqlOrderRepository(connection)
	orderItemRepo := _orderItemRepo.NewMysqlOrderItemRepository(connection)
	orderItem := _orderItemUcase.NewOrderItemUsecase(orderItemRepo, orderRepo, productRepo, timeoutContext)
	_orderItemHttpDelivery.NewOrderItemHandler(e, orderItem, cfg)

	paymentRepo := _paymentRepo.NewMysqlPaymentRepository(connection)
	payment := _paymentUcase.NewPaymentUsecase(orderRepo, paymentRepo,  timeoutContext)
	_paymentHttpDelivery.NewPaymentHandler(e, payment, cfg)

	media := _mediaUcase.NewCustomerUsecase(cfg,  timeoutContext)
	_mediaHttpDelivery.NewCustomerHandler(e, media, cfg)


	// run server
	go func() {
		address := fmt.Sprintf("%s", cfg.Server.Addr)

		e.Static("/","public")
		if err := e.Start(address); err != nil {
			log.Fatal(e.Start(address))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}
