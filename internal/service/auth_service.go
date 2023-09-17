package service

// AuthAdminService - admin APIs for auth data
type AuthAdminService interface {
	// OrganizationService base interface
	OrganizationService

	// PrincipalService base interface
	PrincipalService

	// ResourceService base interface
	ResourceService

	// PermissionService base interface
	PermissionService

	// 	RoleService base interface
	RoleService

	// 	GroupService base interface
	GroupService

	// 	RelationshipService base interface
	RelationshipService
}
