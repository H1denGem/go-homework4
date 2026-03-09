package handlers

import (
	"homework04/models"
	"homework04/services"
	"homework04/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req models.CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	if err := h.postService.CreatePost(req.Title, req.Content, int64(userID)); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"message": "文章创建成功",
	})
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.postService.GetPosts()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, posts)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的文章ID")
		return
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的文章ID")
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	if err := h.postService.UpdatePost(id, req.Title, req.Content, int64(userID)); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"message": "文章更新成功",
	})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的文章ID")
		return
	}

	userID := c.GetUint64("user_id")
	if err := h.postService.DeletePost(id, int64(userID)); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"message": "文章删除成功",
	})
}
