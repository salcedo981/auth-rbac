package ctrRbac

import (
	"errors"
	"fmt"
	errRbac "go_template_v3/pkg/services/rbac/error"
	mdlRbac "go_template_v3/pkg/services/rbac/model"
	scpRbac "go_template_v3/pkg/services/rbac/script"

	"net/http"
	"strconv"
	"strings"

	// v1 "iprovidence/pkg/middleware/v1"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	"github.com/gofiber/fiber/v3"
)

// ----------------------------
// Role Permissions
// ----------------------------

func AssignRolePermission(c fiber.Ctx) error {
	roleIDStr := c.Params("roleId") // FIXED
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	var req mdlRbac.PermissionToRoleReq

	// FIX: Use BodyParser
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_301, "Invalid request body.", err, http.StatusBadRequest,
		)
	}

	// Validation
	if roleID <= 0 || req.ActionName == "" || req.ResourceName == "" {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Missing required fields.", http.StatusBadRequest,
		)
	}

	// Call service
	assignPermToRole, err := scpRbac.AssignRolePermission(roleID, req.ActionName, req.ResourceName)
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to assign permission to role.", err, http.StatusInternalServerError)
	}

	if !assignPermToRole.Success {
		return v1.JSONResponse(c, respcode.ERR_CODE_400, assignPermToRole.Message, http.StatusBadRequest)
	}

	return v1.JSONResponse(
		c, respcode.SUC_CODE_200, assignPermToRole.Message, http.StatusOK,
	)
}
func GetAllRolesPermissions(c fiber.Ctx) error {

	allPerms, err := scpRbac.GetAllRolePermissionsGrouped()
	if err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, "Failed to fetch roles and permissions.", err, http.StatusInternalServerError,
		)
	}

	return v1.JSONResponseWithData(
		c, respcode.SUC_CODE_200, "Successfully fetched all role permissions!", allPerms, http.StatusOK,
	)
}

func GetRolePermissionsbyRole(c fiber.Ctx) error {
	// 1️⃣ Get role ID from params
	roleIDStr := c.Params("roleId")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil || roleID <= 0 {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Invalid role ID", http.StatusBadRequest,
		)
	}

	// 2️⃣ Fetch role permissions
	rolePerms, err := scpRbac.GetRolePermissionsByRole(roleID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_404, "Role not found.", http.StatusNotFound,
			)
		}

		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, "Failed to fetch role permissions.", err, http.StatusInternalServerError,
		)
	}

	// 4️⃣ Success
	return v1.JSONResponseWithData(
		c,
		respcode.SUC_CODE_200,
		"Fetching role permissions successful!",
		rolePerms,
		http.StatusOK,
	)
}

// DeleteRolePermission removes a permission from a role
func RemoveRolePermission(c fiber.Ctx) error {

	roleIDStr := c.Params("roleId") // FIXED
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	var req mdlRbac.PermissionToRoleReq

	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_301, "Invalid request body.", err, http.StatusBadRequest,
		)
	}

	// Validation
	if roleID <= 0 || req.ActionName == "" || req.ResourceName == "" {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Missing required fields.", http.StatusBadRequest,
		)
	}

	// Call service
	removePermResp, err := scpRbac.RemoveRolePermission(roleID, req.ActionName, req.ResourceName)
	if err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, removePermResp.Message, err, http.StatusInternalServerError,
		)
	}

	if !removePermResp.Success {
		return v1.JSONResponse(c, respcode.ERR_CODE_400, removePermResp.Message, http.StatusBadRequest)
	}
	cleanMessage := strings.ReplaceAll(removePermResp.Message, `\"`, "")

	return v1.JSONResponse(
		c, respcode.SUC_CODE_200, cleanMessage, http.StatusOK,
	)
}

// ----------------------------
// Permissions
// ----------------------------

// CreatePermission adds a new permission to the database
func CreatePermission(c fiber.Ctx) error {
	type CreatePermissionRequest struct {
		Name string `json:"name"`
	}

	var req CreatePermissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_301, "Invalid request body.", err, http.StatusBadRequest)
	}

	if req.Name == "" {
		return v1.JSONResponse(c, respcode.ERR_CODE_401, "Missing permission name.", http.StatusBadRequest)
	}

	err := scpRbac.CreatePermission(req.Name)
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to create permission.", err, http.StatusInternalServerError)
	}

	return v1.JSONResponse(c, respcode.SUC_CODE_200, "Permission created successfully.", http.StatusOK)
}

// FetchPermissions retrieves all permissions from the database
func FetchAllPermissions(c fiber.Ctx) error {
	users, err := scpRbac.FetchPermissions()
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_502, "Failed to fetch permissions.", err, http.StatusInternalServerError)
	}

	if len(users) == 0 {
		return v1.JSONResponse(c, respcode.SUC_CODE_204, "No permission found.", http.StatusOK)
	}

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200, "Permissions fetched successfully.", users, http.StatusOK)
}

// UpdatePermission updates an existing permission in the database
func UpdatePermission(c fiber.Ctx) error {
	permissionIDStr := c.Params("permission_id")
	permissionID, err := strconv.Atoi(permissionIDStr)
	if err != nil || permissionID == 0 {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_401, "Invalid or missing permission_id parameter", http.StatusBadRequest,
		)
	}

	type UpdatePermissionRequest struct {
		PermissionName string `json:"permission_name"`
	}

	var req UpdatePermissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_301, "Invalid request body.", err, http.StatusBadRequest)
	}

	if req.PermissionName == "" {
		return v1.JSONResponse(c, respcode.ERR_CODE_401, "Missing permission name.", http.StatusBadRequest)
	}

	err = scpRbac.UpdatePermission(permissionID, req.PermissionName)
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to update permission.", err, http.StatusInternalServerError)
	}

	return v1.JSONResponse(c, respcode.SUC_CODE_200, "Permission updated successfully.", http.StatusOK)
}

// DeletePermission removes a permission by its ID
func DeletePermission(c fiber.Ctx) error {
	permissionIDStr := c.Params("permission_id")
	permissionID, err := strconv.Atoi(permissionIDStr)
	if err != nil || permissionID == 0 {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_401, "Invalid or missing permission_id parameter.", http.StatusBadRequest,
		)
	}

	// Run delete script
	err = scpRbac.DeletePermission(permissionID)
	if err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, "Failed to delete permission.", err, http.StatusInternalServerError,
		)
	}

	return v1.JSONResponse(
		c, respcode.SUC_CODE_200, "Permission deleted successfully.", http.StatusOK,
	)
}

// ----------------------------
//  USER
// ----------------------------

// AssignRoleToUser assigns a role to a user
func AssignUserRole(c fiber.Ctx) error {
	staffID := c.Params("staffId")
	roleIDStr := c.Params("roleId")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	if staffID == "" || roleID == 0 {
		return v1.JSONResponse(c, respcode.ERR_CODE_401, "Missing staff_id or role_id.", http.StatusBadRequest)
	}

	if err := scpRbac.AssignUserRole(staffID, roleID); err != nil {
		if errors.Is(err, errRbac.ErrResourceNotFound) {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_404, "User not found.", http.StatusNotFound,
			)
		}
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to assign role to user.", err, http.StatusInternalServerError)
	}

	return v1.JSONResponse(c, respcode.SUC_CODE_200, "Role assigned to user successfully.", http.StatusOK)
}

func FetchAllUserRoles(c fiber.Ctx) error {

	userRoles, err := scpRbac.FetchAllUserRoles()

	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_502, "Failed to fetch user roles.", err, http.StatusInternalServerError)
	}

	// Only check for nil pointer here since userRoles is a pointer type.
	if userRoles == nil {
		return v1.JSONResponse(c, respcode.SUC_CODE_204, "No user roles found!", http.StatusOK)
	}

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200, "User roles fetched successfully", userRoles, http.StatusOK)
}

// ----------------------------
//
//	ACTION
//
// ----------------------------
// CreateAction - Create a new action
func CreateAction(c fiber.Ctx) error {
	var req mdlRbac.RbacItemRequest

	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Invalid request body.", http.StatusBadRequest,
		)
	}

	if req.Name == "" {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Action name is required.", http.StatusBadRequest,
		)
	}

	err := scpRbac.CreateAction(req.Name, req.Description)
	if err != nil {
		if errors.Is(err, errRbac.ErrResourceNameTaken) {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_409, "Action name already exists.", http.StatusConflict,
			)
		}
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, "Failed to create action.", err, http.StatusInternalServerError,
		)
	}

	return v1.JSONResponse(
		c, respcode.SUC_CODE_201, "Action created successfully!", http.StatusOK,
	)
}

// GetActions - Get all actions
func GetActions(c fiber.Ctx) error {
	fmt.Println("GetActions Called")
	actions, err := scpRbac.GetActions()
	if err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, "Failed to fetch actions.", err, http.StatusInternalServerError,
		)
	}

	return v1.JSONResponseWithData(
		c, respcode.SUC_CODE_200, "Actions fetched successfully!",
		map[string]interface{}{"actions": actions}, http.StatusOK,
	)
}

// UpdateAction - Update an action
func UpdateAction(c fiber.Ctx) error {
	var req mdlRbac.RbacItemRequest
	actionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Invalid action ID.", http.StatusBadRequest,
		)
	}

	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Invalid request body.", http.StatusBadRequest,
		)
	}

	if req.Name == "" {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Action name is required.", http.StatusBadRequest,
		)
	}

	err = scpRbac.UpdateAction(actionID, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, errRbac.ErrResourceNotFound) {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_404, "Action not found.", http.StatusNotFound,
			)
		}
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, "Failed to update action.", err, http.StatusInternalServerError,
		)
	}

	return v1.JSONResponse(
		c, respcode.SUC_CODE_200, "Action updated successfully!", http.StatusOK,
	)
}

// DeleteAction - Delete an action
func DeleteAction(c fiber.Ctx) error {
	actionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Invalid action ID.", http.StatusBadRequest,
		)
	}

	err = scpRbac.DeleteAction(actionID)
	if err != nil {
		if errors.Is(err, errRbac.ErrResourceNotFound) {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_404, "Action not found.", http.StatusNotFound,
			)
		}
		if errors.Is(err, errRbac.ErrResourceInUse) {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_409, "Action is in use, can't be deleted.", http.StatusConflict,
			)
		}
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, "Failed to delete action.", err, http.StatusInternalServerError,
		)
	}

	return v1.JSONResponse(
		c, respcode.SUC_CODE_200, "Action deleted successfully!", http.StatusOK,
	)
}

// ----------------------------
//
//	RESOURCES
//
// ----------------------------
func CreateResource(c fiber.Ctx) error {
	var req mdlRbac.RbacItemRequest

	// Parse request body
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Invalid request body.", http.StatusBadRequest,
		)
	}

	// Validate fields
	if req.Name == "" {
		return v1.JSONResponse(
			c, respcode.ERR_CODE_400, "Resource name is required.", http.StatusBadRequest,
		)
	}

	// Call service
	err := scpRbac.CreateResource(req.Name, req.Description)
	if err != nil {

		if errors.Is(err, errRbac.ErrResourceNameTaken) {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_409, "Resource name taken.", http.StatusConflict,
			)
		}
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, "Failed to create resource.", err, http.StatusInternalServerError,
		)
	}

	// Success
	return v1.JSONResponse(
		c, respcode.SUC_CODE_201, "Resource created successfully!", http.StatusOK,
	)
}

// GetResources - Get all resources
func GetResources(c fiber.Ctx) error {
	// Call service
	resources, err := scpRbac.GetResources()
	if err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_500, "Failed to fetch resources.", err, http.StatusInternalServerError,
		)
	}

	// Success
	return v1.JSONResponseWithData(
		c,
		respcode.SUC_CODE_200,
		"Resources fetched successfully!",
		map[string]interface{}{
			"resources": resources,
		},
		http.StatusOK,
	)
}

// UpdateResource - Update a resource
func UpdateResource(c fiber.Ctx) error {
	var req mdlRbac.RbacItemRequest
	// Get resource ID from URL parameter
	resource := c.Params("id")
	resourceID, err := strconv.Atoi(resource)
	if err != nil {
		return v1.JSONResponse(
			c,
			respcode.ERR_CODE_400,
			"Invalid resource ID.",
			http.StatusBadRequest,
		)
	}

	// Parse request body
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponse(
			c,
			respcode.ERR_CODE_400,
			"Invalid request body.",
			http.StatusBadRequest,
		)
	}

	// req.ID = resourceID

	// Validate fields
	if resourceID == 0 {
		return v1.JSONResponse(
			c,
			respcode.ERR_CODE_400,
			"Resource ID is required.",
			http.StatusBadRequest,
		)
	}

	if req.Name == "" {
		return v1.JSONResponse(
			c,
			respcode.ERR_CODE_400,
			"Resource name is required.",
			http.StatusBadRequest,
		)
	}

	// Call service
	err = scpRbac.UpdateResource(resourceID, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, errRbac.ErrResourceNotFound) {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_404, "Resource not found.", http.StatusNotFound,
			)
		}
		return v1.JSONResponseWithError(
			c,
			respcode.ERR_CODE_500,
			"Failed to update resource.",
			err,
			http.StatusInternalServerError,
		)
	}

	// Success
	return v1.JSONResponse(
		c,
		respcode.SUC_CODE_200,
		"Resource updated successfully!",
		http.StatusOK,
	)
}

// DeleteResource - Delete a resource
func DeleteResource(c fiber.Ctx) error {
	// Get resource ID from URL parameter
	resource := c.Params("id")
	resourceID, err := strconv.Atoi(resource)
	if err != nil {
		return v1.JSONResponse(
			c,
			respcode.ERR_CODE_400,
			"Invalid resource ID.",
			http.StatusBadRequest,
		)
	}

	// Call service
	err = scpRbac.DeleteResource(resourceID)
	if err != nil {
		if errors.Is(err, errRbac.ErrResourceNotFound) {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_404, "Resource not found.", http.StatusNotFound,
			)
		}
		if errors.Is(err, errRbac.ErrResourceInUse) {
			return v1.JSONResponse(
				c, respcode.ERR_CODE_409, "Resource is in use, can't be deleted.", http.StatusConflict,
			)
		}
		return v1.JSONResponseWithError(
			c,
			respcode.ERR_CODE_500,
			"Internal server error while deleting resource",
			err,
			http.StatusInternalServerError,
		)
	}

	// Success
	return v1.JSONResponse(
		c,
		respcode.SUC_CODE_200,
		"Resource deleted successfully!",
		http.StatusOK,
	)
}
