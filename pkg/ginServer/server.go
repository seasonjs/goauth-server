package ginServer

import (
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	oRedis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"oauthServer/repository"
	"strconv"
	"sync"
	"time"
)

var (
	gServer        *server.Server
	once           sync.Once
	userRepository = repository.NewRepositories().UserRepository
)

// InitServer Initialize the service
func InitServer(manager oauth2.Manager) *server.Server {
	gServer = server.NewDefaultServer(manager)
	return gServer
}

// HandleTokenRequest access token登陆方式为返回授权模式
func HandleTokenRequest(c *gin.Context) {
	err := gServer.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	//err.Track
	c.Abort()
}

func HandleAuthorizeRequest(c *gin.Context) {
	err := gServer.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Abort()
}

func ValidationBearerToken(c *gin.Context) {
	token, err := gServer.ValidationBearerToken(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
		"scope":      token.GetScope(),
	}
	c.JSON(http.StatusOK, data)
	return
}

// PasswordAuthorizationService 通过密码进行校验，返回access token
func PasswordAuthorizationService(username string, password string) (userID string, err error) {
	user, err := userRepository.FindUserByNameAndPassword(username, password)
	if err != nil {
		log.Println("错误：", err)
		return "", err
	}
	log.Println("find user", user, err)
	userID = strconv.Itoa(user.ID)
	return userID, nil
}

// UserAuthorizationService 通过Authorization校验跳转返回
func UserAuthorizationService(w http.ResponseWriter, r *http.Request) (userID string, err error) {

	users, err := userRepository.FindUsers()
	//usersCache, cErr := userRepository.FindUsers()

	if err != nil {
		log.Println(err)
	}
	log.Println("find user", users, err)

	userID = "test"
	return userID, nil
}

// Default 默认的权限服务初始化配置
func Default() *server.Server {
	once.Do(func() {
		manager := manage.NewDefaultManager()
		manager.MapAccessGenerate(generates.NewAccessGenerate())
		manager.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
		// token 使用redis
		manager.MapTokenStorage(oRedis.NewRedisStore(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   15,
		}))
		// 模拟数据,已经注册 store，这块只是web应用的话直接注入到内存也是合理的
		clientStore := store.NewClientStore()
		clientStore.Set("000001", &models.Client{
			ID: "000000",
			//Secret: "999999",
			Domain: "http://localhost",
			//UserID: "test",
		})
		manager.MapClientStorage(clientStore)
		//manager.MapTokenStorage()
		// 初始化标准的 oauth2 服务
		InitServer(manager)

		//access_token时调用的方法，需要返回userID
		SetPasswordAuthorizationHandler(PasswordAuthorizationService)
		//Authorization登陆模式的时候需要返回UserID
		SetUserAuthorizationHandler(UserAuthorizationService)

		SetInternalErrorHandler(func(err error) (re *errors.Response) {
			log.Println("Internal Error:", err.Error())
			return
		})

		SetResponseErrorHandler(func(re *errors.Response) {
			log.Println("Response Error:", re.Error.Error())
		})

		SetAllowGetAccessRequest(true)
		SetClientInfoHandler(server.ClientFormHandler)
		SetForcePKCE(false)
	})
	return gServer

}
