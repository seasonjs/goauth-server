package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"oauthServer/pkg/ginServer"
	"oauthServer/pkg/rbac"
	"oauthServer/repository"
	"oauthServer/utils"
	"os"
	"strconv"
	"time"
)

var (
	userRepository       = repository.NewRepositories().UserRepository
	r                    = rbac.NewRBAC()
	gServer              = ginServer.Default()
	roleRepository       = repository.NewRepositories().RoleRepository
	permissionRepository = repository.NewRepositories().PermissionRepository
	//userCacheRepository = autowired.New().Repository.UserCacheRepository
)

// AccessCheckHandler godoc
// @Summary 访问校验
// @Description apisix登陆与权限校验，必须实现
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Router /wolf/rbac/access_check [get]
func AccessCheckHandler(c *gin.Context) {
	utils.FmtRequest(os.Stdout, "AccessCheckHandler", c.Request) // Ignore the error
	log.Println("resName：", c.Query("resName"), c.Query("appID"), c.Query("action"), c.Query("clientIP"))
	ctx := c.Request.Context()
	token := c.GetHeader("X-Rbac-Token")
	tokenInfo, err := gServer.Manager.LoadAccessToken(ctx, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取token信息失败",
		})
		return
	}
	log.Println("tokenInfo", tokenInfo.GetUserID())

	userId, err := strconv.Atoi(tokenInfo.GetUserID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "非法userID",
		})
		return
	}
	user, err := userRepository.FindUserById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "未正常获得用户信息",
		})
		return
	}
	roles, err := roleRepository.FindRoleByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "未正常获得角色信息",
		})
		return
	}
	data := map[string]interface{}{
		"expires_in": int64(tokenInfo.GetAccessCreateAt().Add(tokenInfo.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  tokenInfo.GetClientID(),
		"user_id":    tokenInfo.GetUserID(),
		"scope":      tokenInfo.GetScope(),
		"user_name":  user.Username,
		"userInfo":   roles,
	}
	for _, role := range roles {
		permission, err := permissionRepository.FindPermissionByPermissionResourceAndRoleType(c.Query("resName"), role.RoleType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "未正常获得角色信息",
			})
			return
		}
		roleId := strconv.Itoa(role.RoleType)
		permissionId := strconv.Itoa(permission.ID)
		acc := r.CheckRoleAccess(roleId, permissionId)
		if acc == true {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
				"data":    data,
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "无权限，禁止访问",
				"data":    data,
			})
			return
		}
	}

}
