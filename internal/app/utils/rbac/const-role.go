package rbac

var RoleViewer = role{
	permission: CarView + CarList,
	scope:      Company,
	name:       "Viewer",
}

var RoleLabeler = role{
	permission: CarView + CarList + CarAdd + CarEditSubmit,
	scope:      Company,
	name:       "Labeler",
}

var RoleChecker = role{
	permission: CarView + CarList + CarAdd + CarEditSubmit + CarApprove + CarCancel + CarDisapprove + CarRevoke,
	scope:      Company,
	name:       "Checker",
}

var RoleAdmin = role{
	permission: CarView + CarList + CarAdd + CarEditSubmit + CarApprove + CarCancel + CarDisapprove + CarRevoke + CompanyView +
		UserView + UserList + /*AdminUserEdit +*/ AdminUserDelete + UserEdit + UserDelete + UserAdd,
	scope: Company,
	name:  "Admin",
}

var RoleKBankAdmin = role{
	permission: CompanyView + CompanyList + CompanyEdit + CompanyAdd + UserView +
		UserList + AdminUserView + AdminUserList + AdminUserEdit + AdminUserDelete + AdminUserAdd + CarView + CarList,
	scope: Global,
	name:  "KBankAdmin",
}

var Roles = map[string]role{
	"Viewer":     RoleViewer,
	"Labeler":    RoleLabeler,
	"Checker":    RoleChecker,
	"Admin":      RoleAdmin,
	"KBankAdmin": RoleKBankAdmin,
}
