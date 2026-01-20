package middleware

import (
	"fmt"
	mdlAuth "go_template_v3/pkg/services/auth/model"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	"github.com/gofiber/fiber/v3"
)

func RequirePermission(permission string) fiber.Handler {
	return func(c fiber.Ctx) error {
		fmt.Println("CHECKING PERMISSIONS...")
		// 1️⃣ Check if user context exists
		rawUser := c.Locals("user")

		if rawUser == nil {
			// "Token missing."
			return v1.JSONResponse(c, respcode.ERR_CODE_111, respcode.ERR_CODE_111_MSG, fiber.StatusUnauthorized)
		}

		// 2️⃣ Validate type assertion
		user, ok := rawUser.(*mdlAuth.UserWithPermissions)
		if !ok || user == nil {
			// Internal server error
			return v1.JSONResponse(c, respcode.ERR_CODE_300, "Invalid user context.", fiber.StatusInternalServerError)
		}

		// 3️⃣ Automatically allow super admin
		if user.RoleName == "super_admin" {
			return c.Next()
		}

		// 4️⃣ Check required permission
		for _, p := range user.Permissions {
			if p == permission {
				return c.Next()
			}
		}
		// 5️⃣ Permission denied
		// "Access denied."
		return v1.JSONResponse(c, respcode.ERR_CODE_105_CD, respcode.ERR_CODE_105_CD_MSG, fiber.StatusForbidden)
	}
}
