package router

import (
	"payso-internal-api/controller"
	"payso-internal-api/service"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

// Setup routes for the application
func SetupRoutes(app *fiber.App) {
	// Setup API group
	api := app.Group("/api")

	// Setup v1 routes
	SetupV1Routes(api)

	log.Info("Routes setup completed")
}

// SetupV1Routes configures all v1 API routes
func SetupV1Routes(api fiber.Router) {
	v1 := api.Group("/v1")

	// Setup merchant routes
	SetupMerchantRoutes(v1)

	// Setup permission routes
	SetupPermissionRoutes(v1)
}

// SetupMerchantRoutes configures all merchant related routes
func SetupMerchantRoutes(router fiber.Router) {
	// Create services and controllers
	merchantService := service.NewMerchantService(nil) // หรือส่ง handler.NewMerchantHandler() เข้าไปหากต้องการ
	merchantController := controller.NewMerchantController(merchantService)

	// Create merchant route group
	merchantRoutes := router.Group("/merchant")

	// Setup merchant endpoints
	merchantRoutes.Get("/", merchantController.GetMerchant)
	merchantRoutes.Post("/", merchantController.CreateMerchant)
	merchantRoutes.Delete("/", merchantController.DeleteMerchant)
}

// SetupPermissionRoutes configures all permission related routes
func SetupPermissionRoutes(router fiber.Router) {
	// Create services and controllers
	permissionService := service.NewPermissionService()
	permissionController := controller.NewPermissionController(permissionService)

	// Create permission route group
	permissionRoutes := router.Group("/permission")

	// Setup permission endpoints
	permissionRoutes.Get("/", permissionController.GetPermissions)
	permissionRoutes.Get("/detail", permissionController.GetPermissionByID)
	permissionRoutes.Post("/", permissionController.CreatePermission)
	permissionRoutes.Put("/", permissionController.UpdatePermission)
	permissionRoutes.Delete("/", permissionController.DeletePermission)
}
