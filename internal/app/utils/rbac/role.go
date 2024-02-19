package rbac

type role struct {
	permission int
	scope      int
	name       string
}

func (r role) Permission() int {
	return r.permission
}

func (r role) Scope() int {
	return r.scope
}

func (r role) Name() string {
	return r.name
}

func (r role) IsRole(role interface{ GetRole() string }) bool {
	return r.name == role.GetRole()
}

func IsViewerCheckerLabeler(role string) bool {
	return role == RoleViewer.name || role == RoleChecker.name || role == RoleLabeler.name
}

func IsCompanyAdmin(role string) bool {
	return role == RoleAdmin.name
}

func IsKBankAdmin(role string) bool {
	return role == RoleKBankAdmin.name
}
