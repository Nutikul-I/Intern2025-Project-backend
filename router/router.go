package router

import (
	"payso-internal-api/controller"
	"payso-internal-api/handler"
	"payso-internal-api/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func SetupRoutes(app *fiber.App) {

	permissionService := service.NewPermissionService()
	permissionController := controller.NewPermissionController(permissionService)
	merchantController := controller.NewMerchantController(service.NewMerchantService(handler.NewMerchantHandler()))
	categoryController := controller.NewCategoryController(service.NewCategoryService(handler.NewCategoryHandler()))
	customerCtl := controller.NewCustomerController(service.NewCustomerService(handler.NewCustomerHandler()))
	productCtl := controller.NewProductController(service.NewProductService(handler.NewProductHandler()))
	discountController := controller.NewDiscountController(service.NewDiscountService(handler.NewDiscountHandler()))

	api := app.Group("/", func(c *fiber.Ctx) error {
		if !strings.Contains(c.Request().URI().String(), "/ping") {
			log.Infof("all : %v", c.Request().URI().String())
		}
		return c.Next()
	})

	product := api.Group("/api/product")
	product.Get("/products", productCtl.GetProducts)
	product.Get("/product", productCtl.GetProduct)
	product.Post("/create-product", productCtl.CreateProduct)
	product.Put("/update-product", productCtl.UpdateProduct)
	product.Delete("/delete-product", productCtl.DeleteProduct)

	cust := api.Group("/api/customer")

	cust.Get("/customers", customerCtl.GetCustomers)
	cust.Get("/customer", customerCtl.GetCustomer)
	cust.Post("/create-customer", customerCtl.CreateCustomer)
	cust.Put("/update-customer", customerCtl.UpdateCustomer)
	cust.Delete("/delete-customer", customerCtl.DeleteCustomer)

	merchant := api.Group("/api/merchant")
	merchant.Get("/merchant", merchantController.GetMerchant)
	merchant.Post("/create-merchant", merchantController.CreateMerchant)
	merchant.Delete("/delete-merchant", merchantController.DeleteMerchant)

	category := api.Group("/api/category")
	category.Get("/category", categoryController.GetCategory)
	category.Post("/create-category", categoryController.CreateCategory)
	category.Put("/update-category", categoryController.UpdateCategory)
	category.Delete("/delete-category", categoryController.DeleteCategory)

	permissionRoutes := api.Group("/api/permission")
	permissionRoutes.Get("/permission", permissionController.GetPermissions)             // GET /api/permission?id=0&page=1&row=10
	permissionRoutes.Get("/detail-permission", permissionController.GetPermissionByID)   // GET /api/permission/detail?id=123
	permissionRoutes.Post("/create-permission", permissionController.CreatePermission)   // POST /api/permission/create
	permissionRoutes.Put("/update-permission", permissionController.UpdatePermission)    // PUT /api/permission/update
	permissionRoutes.Delete("/delete-permission", permissionController.DeletePermission) // DELETE /api/permission/delete?id=123

	discount := api.Group("/api/discount")
	discount.Get("/discount", discountController.GetDiscount)
	discount.Post("/create-discount", discountController.CreateDiscount)
	discount.Put("/update-discount", discountController.UpdateDiscount)
	discount.Delete("/delete-discount", discountController.DeleteDiscount)
}
