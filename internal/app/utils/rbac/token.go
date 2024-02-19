package rbac

type Token interface {
	GetRole() string
	GetUserId() string
	GetCompanyId() string
}
