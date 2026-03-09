package middlewares

import (
	"homework04/utils"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			log.Printf("认证失败: 缺少Authorization头")
			utils.Error(c, utils.ErrMissingToken.Code, utils.ErrMissingToken.Message)
			c.Abort()
			return
		}

		// 提取 Token（Bearer <token>）
		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("认证失败: 无效的Authorization格式")
			utils.Error(c, utils.ErrInvalidToken.Code, "无效的授权头格式")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 Token
		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			log.Printf("认证失败: 无效的令牌, Error=%v", err)
			utils.Error(c, utils.ErrInvalidToken.Code, utils.ErrInvalidToken.Message)
			c.Abort()
			return
		}

		// 将用户信息存储到 Context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		log.Printf("认证成功: UserID=%d, Username=%s", claims.UserID, claims.Username)
		c.Next()
	}
}
