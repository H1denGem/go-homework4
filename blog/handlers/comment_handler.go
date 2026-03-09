package handlers

import (
	"homework04/models"
	"homework04/services"
	"homework04/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的文章ID")
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	if err := h.commentService.CreateComment(req.Content, postID, int64(userID)); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"message": "评论创建成功",
	})
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的文章ID")
		return
	}

	comments, err := h.commentService.GetCommentsByPostID(postID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, comments)
}

func (h *CommentHandler) GetComment(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的文章ID")
		return
	}

	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的评论ID")
		return
	}

	comment, err := h.commentService.GetCommentByID(commentID, postID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, comment)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的文章ID")
		return
	}

	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的评论ID")
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	if err := h.commentService.UpdateComment(commentID, postID, int64(userID), req.Content); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"message": "评论更新成功",
	})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的文章ID")
		return
	}

	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrBadRequest.Code, "无效的评论ID")
		return
	}

	userID := c.GetUint64("user_id")
	if err := h.commentService.DeleteComment(commentID, postID, int64(userID)); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, gin.H{
		"message": "评论删除成功",
	})
}
