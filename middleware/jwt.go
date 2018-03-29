package middleware

import (
	"time"

	"github.com/mohuishou/scuplus-go/config"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
)

func jwtMiddle(ctx iris.Context) {

	// 登录页面无需验证
	if skipJWT(ctx.Path()) {
		ctx.Next()
		return
	}

	// token 验证
	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Get().JwtSecret), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	// jwt 验证
	if err := jwtHandler.CheckJWT(ctx); err != nil {
		ctx.StopExecution()
		return
	}

	// token信息验证
	token := ctx.Values().Get("jwt").(*jwt.Token)
	userID, ok := token.Claims.(jwt.MapClaims)["user_id"]
	if !ok {
		ctx.JSON(map[string]interface{}{
			"status": 401,
			"msg":    "用户尚未登录，获取用户信息失败",
		})
		ctx.StopExecution()
		return
	}

	// token 时效验证
	// end, ok := token.Claims.(jwt.MapClaims)["end"].(float64)
	// if !ok || time.Now().Unix() > int64(end) {
	// 	log.Println("[Error]: 登录信息已失效", end)
	// 	ctx.JSON(map[string]interface{}{
	// 		"status": 401,
	// 		"msg":    "用户尚未登录，获取用户信息失败",
	// 	})
	// 	ctx.StopExecution()
	// 	return
	// }

	// 设置用户id
	ctx.Values().Set("user_id", userID)

	ctx.Next()
}

// 跳过jwt的链接
func skipJWT(path string) bool {
	urls := []string{
		"/login",
		"/notices",
		"/webhook",
		"/spider/webhook",
		"/spider/jwc/cookies",
	}
	for _, v := range urls {
		if v == path {
			return true
		}
	}
	return false
}

// GetUserID 获取用户的id
func GetUserID(ctx iris.Context) uint {
	return uint(ctx.Values().Get("user_id").(float64))
}

// CreateToken 新建一个Token
func CreateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"end":     time.Now().Unix() + 3600*24*15,
		"start":   time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(config.Get().JwtSecret))
}
