package rbac

const ( // Permission Table
	CarView = 1 << iota
	CarList
	CarAdd
	CarEditSubmit
	CarDisapprove
	CarApprove
	CarCancel
	CarRevoke
	CompanyView
	CompanyList
	CompanyEdit
	CompanyAdd
	UserView
	UserList
	UserEdit
	UserDelete
	UserAdd
	AdminUserView
	AdminUserList
	AdminUserEdit
	AdminUserDelete
	AdminUserAdd
)

const (
	Global = iota
	Company
)
