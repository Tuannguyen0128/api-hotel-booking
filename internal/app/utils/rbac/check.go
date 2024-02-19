package rbac

import (
	"fmt"

	"api-hotel-booking/internal/app"
)

func IsHavePermission(token Token, requirePermission int, requireCompany string) error {
	role, found := Roles[token.GetRole()]
	if !found {
		return app.ErrorMap.PermissionDenied.AddDebug(fmt.Sprintf("unknown role name %s", token.GetRole()))
	}
	if role.permission&requirePermission != requirePermission {
		return app.ErrorMap.PermissionDenied.AddDebug(fmt.Sprintf("user id %s role %s have no permission %d", token.GetUserId(), token.GetRole(), requirePermission))
	}

	if role.Scope() != Global && token.GetCompanyId() != requireCompany {
		return app.ErrorMap.PermissionDenied.AddDebug(fmt.Sprintf("user id %s role %s have no permission %d on scope %s", token.GetUserId(), token.GetRole(), requirePermission, requireCompany))
	}

	return nil
}
