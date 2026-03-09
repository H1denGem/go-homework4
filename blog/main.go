package main

import (
	"homework04/config"
	"homework04/handlers"
	"homework04/middlewares"
	"homework04/models"
	"homework04/services"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

import "gorm.io/driver/sqlite"
import "gorm.io/gorm"

func main() {

	// viper 配置文件
	config.Init()
	cfg := config.GetConfig()

	// 初始化日志文件
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("打开日志文件失败: %v", err)
		panic(err)
	}
	defer logFile.Close()

	// 同时输出到文件和标准输出
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
		panic(err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
		panic(err)
	}

	// 初始化服务
	userService := services.NewUserService(db)
	postService := services.NewPostService(db)
	commentService := services.NewCommentService(db)

	userHandler := handlers.NewUserHandler(userService, []byte(cfg.JWT.Secret))
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)

	// 创建 gin 实例
	r := gin.Default()

	// 全局中间件
	r.Use(middlewares.CORS())
	r.Use(middlewares.Logger())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/register", userHandler.Register)                           // 用户注册
		public.POST("/login", userHandler.Login)                                 // 用户登录
		public.GET("/posts", postHandler.GetPosts)                               // 获取文章列表
		public.GET("/posts/:id", postHandler.GetPost)                            // 获取单篇文章
		public.GET("/posts/:id/comments", commentHandler.GetComments)            // 获取文章评论列表
		public.GET("/posts/:id/comments/:comment_id", commentHandler.GetComment) // 获取单条评论
	}

	// 需要认证的路由
	auth := r.Group("/api/v1")
	auth.Use(middlewares.Auth([]byte(cfg.JWT.Secret)))
	{
		auth.POST("/posts", postHandler.CreatePost)                                  // 创建文章
		auth.PUT("/posts/:id", postHandler.UpdatePost)                               // 更新文章（仅作者）
		auth.DELETE("/posts/:id", postHandler.DeletePost)                            // 删除文章（仅作者）
		auth.POST("/posts/:id/comments", commentHandler.CreateComment)               // 创建评论
		auth.PUT("/posts/:id/comments/:comment_id", commentHandler.UpdateComment)    // 更新评论（仅作者）
		auth.DELETE("/posts/:id/comments/:comment_id", commentHandler.DeleteComment) // 删除评论（仅作者）
	}

	// 启动服务
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("服务启动成功: %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
		panic(err)
	}
}
