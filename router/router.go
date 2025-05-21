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
	categoryController := controller.NewCategoryController(service.NewCategoryService(handler.NewcategoryHandler()))

	api := app.Group("/", func(c *fiber.Ctx) error {
		if !strings.Contains(c.Request().URI().String(), "/ping") {
			log.Infof("all : %v", c.Request().URI().String())
		}
		return c.Next()
	})

	productCtl := controller.NewProductController(
		service.NewProductService(handler.NewProductHandler()),
	)
	product := api.Group("/api/product")
	product.Get("/products", productCtl.GetProducts)
	product.Get("/product", productCtl.GetProduct)
	product.Post("/create-product", productCtl.CreateProduct)
	product.Put("/update-product", productCtl.UpdateProduct)
	product.Delete("/delete-product", productCtl.DeleteProduct)

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

<<<<<<< HEAD
	category := api.Group("/api/category")
	category.Get("/category", categoryController.GetCategory)
	category.Post("/create-category", categoryController.CreateCategory)
	category.Put("/update-category", categoryController.UpdateCategory)
	category.Delete("/delete-category", categoryController.DeleteCategory)
=======
	permissionRoutes := api.Group("/api/permission")

	permissionRoutes.Get("/", permissionController.GetPermissions)            // GET /api/permission?id=0&page=1&row=10
	permissionRoutes.Get("/detail", permissionController.GetPermissionByID)   // GET /api/permission/detail?id=123
	permissionRoutes.Post("/create", permissionController.CreatePermission)   // POST /api/permission/create
	permissionRoutes.Put("/update", permissionController.UpdatePermission)    // PUT /api/permission/update
	permissionRoutes.Delete("/delete", permissionController.DeletePermission) // DELETE /api/permission/delete?id=123
>>>>>>> feature/permission1

}
