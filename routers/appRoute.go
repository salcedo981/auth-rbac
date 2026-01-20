package routers

import (
	"go_template_v3/pkg/middleware"
	ctrAuth "go_template_v3/pkg/services/auth/controller"
	svcHealthcheck "go_template_v3/pkg/services/healthcheck"
	officesController "go_template_v3/pkg/services/offices/controller"
	ctrRbac "go_template_v3/pkg/services/rbac/controller"

	"github.com/gofiber/fiber/v3"
)

func APIRoute(app *fiber.App) {
	publicV1 := app.Group("/api/public/v1")
	privateV1 := app.Group("/api/private/v1")

	// HealthCheck
	publicV1.Get("/", svcHealthcheck.HealthCheck)
	privateV1.Get("/", svcHealthcheck.HealthCheck)

	// // Sample
	// sampleEndpoint := publicV1.Group("/sample")
	// sampleEndpoint.Get("/", ctrFeatureOne.GetSampleData)

	auth := publicV1.Group("/auth")
	auth.Post("/register", ctrAuth.RegisterUser)
	auth.Post("/login", ctrAuth.LoginUser)
	auth.Post("/logout", ctrAuth.LogoutUser)
	auth.Post("/change-temp-password", ctrAuth.ChangeTempPassword)
	auth.Post("/delete-user", ctrAuth.DeleteUser)
	auth.Post("/update-user/:username", ctrAuth.UpdateUser)
	auth.Post("/forgot-password", ctrAuth.ForgotPassword)
	auth.Post("/verify-reset-token", ctrAuth.VerifyResetToken)

	// ----------------------------
	// 🔐 RBAC Endpoints
	// ----------------------------
	rbac := publicV1.Group("/rbac", middleware.AuthMiddleware)
	// rbac.Get("/getmenubyrole", ctrRbac.GetUserMenus)
	rbac.Get("/roles", ctrRbac.FetchAllUserRoles, middleware.RequirePermission("view:role"))
	rbac.Put("/users/:staffId/roles/:roleId", ctrRbac.AssignUserRole, middleware.RequirePermission("update:role"))

	//CRUD Actions
	rbac.Post("/actions", ctrRbac.CreateAction, middleware.RequirePermission("create:action"))
	rbac.Get("/actions", ctrRbac.GetActions, middleware.RequirePermission("view:action"))
	rbac.Put("/actions/:id", ctrRbac.UpdateAction, middleware.RequirePermission("update:action"))
	rbac.Delete("/actions/:id", ctrRbac.DeleteAction, middleware.RequirePermission("delete:action"))

	// CRUD Resources
	rbac.Post("/resources", ctrRbac.CreateResource, middleware.RequirePermission("create:action"))
	rbac.Get("/resources", ctrRbac.GetResources, middleware.RequirePermission("view:action"))
	rbac.Put("/resources/:id", ctrRbac.UpdateResource, middleware.RequirePermission("update:action"))
	rbac.Delete("/resources/:id", ctrRbac.DeleteResource, middleware.RequirePermission("delete:action"))

	// // ROle permissions Assignment
	rbac.Post("/roles/:roleId/permissions", ctrRbac.AssignRolePermission, middleware.RequirePermission("create:permission"))
	rbac.Get("/roles/permissions", ctrRbac.GetAllRolesPermissions, middleware.RequirePermission("view:permission"))
	rbac.Get("/roles/:roleId/permissions", ctrRbac.GetRolePermissionsbyRole, middleware.RequirePermission("view:permission"))
	rbac.Delete("/roles/:roleId/permissions", ctrRbac.RemoveRolePermission, middleware.RequirePermission("delete:permission"))

	// ----------------------------
	//  OFFICES Endpoints
	// ----------------------------

	offices := publicV1.Group("/offices", middleware.AuthMiddleware)
	offices.Get("/branches", officesController.GetBranches)
	offices.Get("/units", officesController.GetUnits)

}
