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

	merchantController := controller.NewMerchantController(service.NewMerchantService(handler.NewMerchantHandler()))

	api := app.Group("/", func(c *fiber.Ctx) error {
		if !strings.Contains(c.Request().URI().String(), "/ping") {
			log.Infof("all : %v", c.Request().URI().String())
		}
		return c.Next()
	})

	merchant := api.Group("/api/merchant")
	merchant.Get("/merchant", merchantController.GetMerchant)
	merchant.Post("/create-merchant", merchantController.CreateMerchant)
	merchant.Delete("/delete-merchant", merchantController.DeleteMerchant)

	permissionService := service.NewPermissionService()
	permissionController := controller.NewPermissionController(permissionService)

	permissionRoutes := api.Group("/api/permission")

	permissionRoutes.Get("/", permissionController.GetPermissions)            // GET /api/permission?id=0&page=1&row=10
	permissionRoutes.Get("/detail", permissionController.GetPermissionByID)   // GET /api/permission/detail?id=123
	permissionRoutes.Post("/create", permissionController.CreatePermission)   // POST /api/permission/create
	permissionRoutes.Put("/update", permissionController.UpdatePermission)    // PUT /api/permission/update
	permissionRoutes.Delete("/delete", permissionController.DeletePermission) // DELETE /api/permission/delete?id=123

}
