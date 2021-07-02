package rbac

import (
	"github.com/mikespook/gorbac"
	"log"
	"oauthServer/repository"
	"strconv"
	"sync"
)

var (
	once                 sync.Once
	roleRepository       = repository.NewRepositories().RoleRepository
	permissionRepository = repository.NewRepositories().PermissionRepository
	r                    *RBAC
)

type RBAC struct {
	gRbac *gorbac.RBAC
}

// NewGoRBAC 角色框架初始化
func newGoRBAC() *RBAC {
	rbac := gorbac.New()

	return &RBAC{
		rbac,
	}
}

func NewRBAC() *RBAC {
	once.Do(func() {
		r = newGoRBAC()
		err := r.LoadAllRolesAndPermission()
		if err != nil {
			panic("初始化角色失败")
		}
	})
	return r
}

func (rbac *RBAC) LoadAllRolesAndPermission() error {
	roles, err := roleRepository.FindGroupByRoleType()
	if err != nil {
		return err
	}

	for _, role := range roles {
		roleId := strconv.Itoa(role.RoleType)
		permissions, err := permissionRepository.FindPermissionByRoleType(role.RoleType)
		if err != nil {
			return err
		}
		stdRole := gorbac.NewStdRole(roleId)
		for _, permission := range permissions {
			log.Println("permission list:", permission.PermissionName, permission.PermissionResource, permission.RoleType)
			permissionId := strconv.Itoa(permission.ID)
			stdPer := gorbac.NewStdPermission(permissionId)
			err := stdRole.Assign(stdPer)
			if err != nil {
				return err
			}
		}
		er := rbac.gRbac.Add(stdRole)
		if er != nil {
			return er
		}
	}
	return nil
}

// CheckRoleAccess 校验角色是否具有权限
func (rbac *RBAC) CheckRoleAccess(roleId string, permissionId string) bool {
	result := rbac.gRbac.IsGranted(roleId, gorbac.NewStdPermission(permissionId), nil)
	log.Println(result)
	return result
}
