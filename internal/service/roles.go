package service

// Проектные роли
const (
	ProjectRoleOwner   = "owner"
	ProjectRoleManager = "manager"
	ProjectRoleMember  = "member"
)

// Проверка, имеет ли пользователь права для создания задач/проектов
func hasManagementRights(role string) bool {
	return role == ProjectRoleOwner || role == ProjectRoleManager
}

// Проверка, имеет ли пользователь доступ к проекту
func hasProjectAccess(role string) bool {
	return role == ProjectRoleOwner || role == ProjectRoleManager || role == ProjectRoleMember
}
