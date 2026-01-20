package mdlRbac

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Menu struct {
	MenuID    int     `json:"menu_id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	SystemTag *string `json:"system_tag,omitempty"`
}

type RoleMenu struct {
	ID     int `json:"id"`
	RoleID int `json:"role_id"`
	MenuID int `json:"menu_id"`
}

type Permissions struct {
	PermissionsID int    `json:"permissions_id"`
	Name          string `json:"name"`
}

type AssignRoleRequest struct {
	StaffID string `json:"staff_id"`
	RoleID  int    `json:"role_id"`
}

type PermissionResult struct {
	Success             bool   `json:"success,omitempty"`
	Message             string `json:"message,omitempty"`
	Role                string `json:"role,omitempty"`
	Action              string `json:"action,omitempty"`
	Resource            string `json:"resource,omitempty"`
	FormattedPermission string `json:"formatted_permission,omitempty"`
}

type RoleWithPermissions struct {
	Role        string           `json:"role"`
	Permissions []PermissionItem `json:"permissions"`
}

type PermissionItem struct {
	Resource  string `json:"resource"`
	Action    string `json:"action"`
	Formatted string `json:"formatted"`
}

type PermissionToRoleReq struct {
	ActionName   string `json:"action"`
	ResourceName string `json:"resource"`
}

type RBACItemResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type RbacItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
