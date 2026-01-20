package scpRbac

import (
	"encoding/json"
	"fmt"
	"go_template_v3/pkg/config"
	errRbac "go_template_v3/pkg/services/rbac/error"
	mdlRbac "go_template_v3/pkg/services/rbac/model"
	"log"
	"strings"
)

// ----------------------------
// Role Permissions
// ----------------------------

// AssignPermissionToRole assigns a permission to a role
func AssignRolePermission(roleID int, actionName string, resourceName string) (*mdlRbac.PermissionResult, error) {
	db := config.DBConnList[0] // remove the pointer & (not needed)

	query := `SELECT * FROM assign_role_permission($1, $2, $3);`

	var result mdlRbac.PermissionResult

	err := db.Raw(query, roleID, actionName, resourceName).Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("error executing assign_permission_from_role: %v", err)
	}

	fmt.Printf("AssignRolePermission => success=%v, msg=%s\n", result.Success, result.Message)

	fmt.Printf("✅ Assigned permission %s to role %d\n", actionName, roleID)
	return &result, nil
}

func GetAllRolePermissionsGrouped() ([]mdlRbac.RoleWithPermissions, error) {
	db := config.DBConnList[0]

	query := `SELECT get_all_roles_permissions_json()::text`

	var jsonStr string
	if err := db.Raw(query).Scan(&jsonStr).Error; err != nil {
		return nil, fmt.Errorf("error executing get_all_roles_permissions: %v", err)
	}

	// Direct unmarshal to slice
	var result []mdlRbac.RoleWithPermissions
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}

// AssignPermissionToRole assigns a permission to a role
func GetRolePermissionsByRole(roleId int) (*mdlRbac.RoleWithPermissions, error) {
	db := config.DBConnList[0]

	query := `SELECT get_role_permissions_json(?)::text`

	var jsonStr string
	if err := db.Raw(query, roleId).Scan(&jsonStr).Error; err != nil {
		// Check if error contains "not found" message
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, fmt.Errorf("error not found")
		}
		return nil, fmt.Errorf("error executing get_role_permissions: %v", err)
	}

	// Direct unmarshal
	var result mdlRbac.RoleWithPermissions
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &result, nil
}

// RemovePermissionFromRole removes a permission from a role
func RemoveRolePermission(roleID int, actionName string, resourceName string) (*mdlRbac.PermissionResult, error) {
	db := &config.DBConnList[0]

	query := `
		SELECT success, message, role_name, action_name, resource_name
		FROM remove_role_permission(?, ?, ?);
	`

	var result mdlRbac.PermissionResult

	// Use Raw().Scan() for functions that return rows
	err := db.Raw(query, roleID, actionName, resourceName).Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("error executing remove_permission_from_role: %v", err)
	}

	fmt.Printf("RemoveRolePermission => success=%v, msg=%s\n", result.Success, result.Message)

	return &result, nil
}

// ----------------------------
// Permissions
// ----------------------------

// AddPermission adds a new permission to the database
func CreatePermission(name string) error {
	db := &config.DBConnList[0]

	query := `
		INSERT INTO public.permissions (name, created_at, updated_at)
		VALUES (?, NOW(), NOW());
	`

	if err := db.Exec(query, name).Error; err != nil {
		return fmt.Errorf("failed to add permission: %v", err)
	}

	fmt.Printf("✅ Added permission '%s'\n", name)
	return nil
}

// FetchPermissions retrieves all permissions from the database
func FetchPermissions() ([]mdlRbac.Permissions, error) {
	db := &config.DBConnList[0]
	var permissions []mdlRbac.Permissions

	query := `
		SELECT 
			id, name
		FROM public.permissions;
	`

	if err := db.Raw(query).Scan(&permissions).Error; err != nil {
		log.Printf("❌ Failed to fetch permissions: %v", err)
		return nil, err
	}
	fmt.Println("✅ Fetched permissions:", permissions)
	return permissions, nil
}

// UpdatePermission updates the name of an existing permission
func UpdatePermission(permissionID int, name string) error {
	db := &config.DBConnList[0]

	query := `
		UPDATE public.permissions
		SET name = ?, updated_at = NOW()
		WHERE id = ?;
	`

	if err := db.Exec(query, name, permissionID).Error; err != nil {
		return fmt.Errorf("failed to update permission: %v", err)
	}

	fmt.Printf("✅ Updated permission %d to '%s'\n", permissionID, name)
	return nil
}

// DeletePermission removes a permission by its ID
func DeletePermission(permissionID int) error {
	db := &config.DBConnList[0]

	query := `
		DELETE FROM public.permissions
		WHERE id = ?;
	`

	if err := db.Exec(query, permissionID).Error; err != nil {
		return fmt.Errorf("failed to delete permission: %v", err)
	}

	fmt.Printf("🗑️ Deleted permission %d\n", permissionID)
	return nil
}

// ----------------------------
// USER
// ----------------------------

// AssignRole assigns a role to a user by staff ID
func AssignUserRole(staffID string, roleID int) error {

	db := &config.DBConnList[0]

	query := `
		UPDATE public.users SET role_id = ? WHERE staff_id = ?
	`
	// exist, err := helperScript.CheckUserExists(staffID)

	// if err != nil {
	// 	return fmt.Errorf("failed to check user existence: %v", err)
	// }

	// if !exist {
	// 	return fmt.Errorf("user with staff ID %v does not exist", staffID)
	// }

	result := db.Exec(query, roleID, staffID)

	if result.Error != nil {
		return fmt.Errorf("failed to update resource: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errRbac.ErrResourceNotFound
	}

	fmt.Printf("✅ Assigned role %d to staff ID %v\n", roleID, staffID)
	return nil
}

func FetchAllUserRoles() ([]mdlRbac.RBACItemResponse, error) {

	db := &config.DBConnList[0]

	var userRoles []mdlRbac.RBACItemResponse

	query := `
		SELECT * FROM roles
	`
	if err := db.Raw(query).Scan(&userRoles).Error; err != nil {
		log.Printf("❌ Failed to fetch user roled: %v", err)
		return nil, err
	}

	return userRoles, nil
}

// ----------------------------
// ACTION
// ----------------------------

// CreateAction - Create a new action in database
func CreateAction(name, description string) error {
	db := &config.DBConnList[0]

	query := `INSERT INTO actions (name, description) VALUES (?, ?)`

	if err := db.Debug().Exec(query, name, description).Error; err != nil {
		if strings.Contains(err.Error(), `unique constraint "actions_name_key"`) {
			return errRbac.ErrResourceNameTaken
		}
		return fmt.Errorf("failed to create action: %v", err)
	}

	return nil
}

// GetActions - Get all actions from database
func GetActions() ([]mdlRbac.RBACItemResponse, error) {
	db := &config.DBConnList[0]
	var actions []mdlRbac.RBACItemResponse

	query := `SELECT id, name, description, created_at, updated_at FROM actions ORDER BY id ASC`

	if err := db.Raw(query).Scan(&actions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch actions: %v", err)
	}

	return actions, nil
}

// GetActionByID - Get single action by ID
func GetActionByID(id int) (*mdlRbac.RBACItemResponse, error) {
	db := &config.DBConnList[0]
	var action mdlRbac.RBACItemResponse

	query := `SELECT id, name, description, created_at, updated_at FROM actions WHERE id = ?`

	if err := db.Raw(query, id).Scan(&action).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch action: %v", err)
	}

	if action.ID == 0 {
		return nil, errRbac.ErrResourceNotFound
	}

	return &action, nil
}

// UpdateAction - Update an action in database
func UpdateAction(id int, name, description string) error {
	db := &config.DBConnList[0]

	query := `UPDATE actions SET name = ?, description = ?, updated_at = NOW() WHERE id = ?`
	result := db.Exec(query, name, description, id)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), `unique constraint "actions_name_key"`) {
			return errRbac.ErrResourceNameTaken
		}
		return fmt.Errorf("failed to update action: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errRbac.ErrResourceNotFound
	}

	return nil
}

// DeleteAction - Delete an action from database
func DeleteAction(id int) error {
	db := &config.DBConnList[0]

	query := `DELETE FROM actions WHERE id = ?`
	result := db.Exec(query, id)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "foreign key") {
			return errRbac.ErrResourceInUse
		}
		return fmt.Errorf("failed to delete action: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errRbac.ErrResourceNotFound
	}

	return nil
}

// ----------------------------
// Resources
// ----------------------------

// CreateResource - Create a new resource in database
func CreateResource(name, description string) error {
	db := &config.DBConnList[0]

	query := `INSERT INTO resources (name, description) VALUES (?, ?)`

	if err := db.Exec(query, name, description).Error; err != nil {
		if strings.Contains(err.Error(), `unique constraint "resources_name_key"`) {
			return errRbac.ErrResourceNameTaken
		}
		return fmt.Errorf("failed to create resource: %v", err)
	}

	return nil
}

// GetResources - Get all resources from database
func GetResources() ([]mdlRbac.RBACItemResponse, error) {
	db := &config.DBConnList[0]
	var resources []mdlRbac.RBACItemResponse

	query := `SELECT id, name, description, created_at, updated_at FROM resources ORDER BY id ASC`

	if err := db.Raw(query).Scan(&resources).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch resources: %v", err)
	}

	return resources, nil
}

// GetResourceByID - Get single resource by ID
func GetResourceByID(id int) (*mdlRbac.RBACItemResponse, error) {
	db := &config.DBConnList[0]
	var resource mdlRbac.RBACItemResponse

	query := `SELECT id, name, description, created_at, updated_at FROM resources WHERE id = ?`

	if err := db.Raw(query, id).Scan(&resource).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch resource: %v", err)
	}

	// Check if resource exists
	if resource.ID == 0 {
		return nil, fmt.Errorf("resource with ID %d not found", id)
	}

	return &resource, nil
}

// UpdateResource - Update a resource in database
func UpdateResource(id int, name, description string) error {
	db := &config.DBConnList[0]

	query := `UPDATE resources SET name = ?, description = ?, updated_at = NOW() WHERE id = ?`
	result := db.Exec(query, name, description, id)

	if result.Error != nil {
		return fmt.Errorf("failed to update resource: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errRbac.ErrResourceNotFound
	}

	return nil
}

func DeleteResource(id int) error {
	db := &config.DBConnList[0]

	query := `DELETE FROM resources WHERE id = ?`
	result := db.Exec(query, id)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "foreign key") {
			return errRbac.ErrResourceInUse
		}
		return fmt.Errorf("failed to delete resource: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errRbac.ErrResourceNotFound
	}

	return nil
}
