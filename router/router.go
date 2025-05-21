package router

import (
	"payso-internal-api/controller"
	"payso-internal-api/handler"
	"payso-internal-api/repository"
	"payso-internal-api/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func SetupRoutes(app *fiber.App) {

	merchantController := controller.NewMerchantController(service.NewMerchantService(handler.NewMerchantHandler()))
	customerCtl := controller.NewCustomerController(service.NewCustomerService(handler.NewCustomerHandler()))
	productCtl := controller.NewProductController(service.NewProductService(handler.NewProductHandler()))

	api := app.Group("/", func(c *fiber.Ctx) error {
		if !strings.Contains(c.Request().URI().String(), "/ping") {
			log.Infof("all : %v", c.Request().URI().String())
		}
		return c.Next()
	})

	product := api.Group("/api/product")
	product.Get("/products", productCtl.GetProducts)
	product.Get("/:id", productCtl.GetProduct)
	product.Post("/", productCtl.CreateProduct)
	product.Put("/:id", productCtl.UpdateProduct)
	product.Delete("/:id", productCtl.DeleteProduct)

	cust := api.Group("/api/customer")

	cust.Get("/customers", customerCtl.GetCustomers)            // GET  /api/customer/customers?page=&row=
	cust.Get("/customer", customerCtl.GetCustomer)              // GET  /api/customer/customer?id=11
	cust.Post("/create-customer", customerCtl.CreateCustomer)   // POST /api/customer/create-customer
	cust.Put("/update-customer", customerCtl.UpdateCustomer)    // PUT  /api/customer/update-customer?id=11
	cust.Delete("/delete-customer", customerCtl.DeleteCustomer) // DELETE /api/customer/delete-customer?id=11

	merchant := api.Group("/api/merchant")
	merchant.Get("/merchant", merchantController.GetMerchant)
	merchant.Post("/create-merchant", merchantController.CreateMerchant)
	merchant.Delete("/delete-merchant", merchantController.DeleteMerchant)

	api.Get("/pingdb", func(c *fiber.Ctx) error {
		db := repository.ConnectDB()
		if err := db.Ping(); err != nil {
			log.Error("DB Ping Failed: ", err)
			return c.Status(500).SendString("❌ DB not connected: " + err.Error())
		}
		return c.SendString("✅ DB connected!")
	})

}
