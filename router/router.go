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

	// Setup permission endpoints
	permissionRoutes.Get("/", permissionController.GetPermissions)
	permissionRoutes.Get("/detail", permissionController.GetPermissionByID) // คงไว้สำหรับ query parameter
	permissionRoutes.Get("/:id", permissionController.GetPermissionByID)    // เพิ่มสำหรับ path parameter
	permissionRoutes.Post("/create-permission", permissionController.CreatePermission)
	permissionRoutes.Put("/update-permission", permissionController.UpdatePermission)
	permissionRoutes.Delete("/delete-permission", permissionController.DeletePermission)
	permissionRoutes.Delete("/delete-permission/:id", permissionController.DeletePermission)

}
