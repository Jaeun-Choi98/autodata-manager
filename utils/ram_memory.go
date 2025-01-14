package utils

import (
	"cju/entity/auth"
	"fmt"
	"sync"
)

var roleIdStore sync.Map

func UploadRoleId(roles []*auth.Role) {
	for _, role := range roles {
		roleIdStore.Store(role.RoleName, role.ID)
	}
}

func GetRoleId(roleName string) (uint, error) {
	val, ok := roleIdStore.Load(roleName)
	if !ok {
		return 0, fmt.Errorf("'%s' is not exsits", roleName)
	}
	return val.(uint), nil
}
