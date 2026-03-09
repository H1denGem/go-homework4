package services

import (
	"homework04/models"
	"homework04/utils"
	"log"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) CreateComment(content string, postID int64, userID int64) error {
	// 检查文章是否存在
	var post models.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("创建评论失败: 文章不存在 PostID=%d", postID)
			return utils.ErrPostNotFound
		}
		log.Printf("查询文章失败: PostID=%d, Error=%v", postID, err)
		return utils.WrapError(err, "查询文章失败")
	}

	comment := models.Comment{
		Content: content,
		PostID:  postID,
		UserID:  userID,
	}

	if err := s.db.Create(&comment).Error; err != nil {
		log.Printf("创建评论失败: PostID=%d, UserID=%d, Error=%v", postID, userID, err)
		return utils.WrapError(err, "创建评论失败")
	}

	log.Printf("评论创建成功: CommentID=%d, PostID=%d, UserID=%d", comment.ID, postID, userID)
	return nil
}

func (s *CommentService) GetCommentsByPostID(postID int64) ([]models.Comment, error) {
	var comments []models.Comment
	err := s.db.Preload("User").Where("post_id = ?", postID).Find(&comments).Error
	if err != nil {
		log.Printf("获取评论列表失败: PostID=%d, Error=%v", postID, err)
		return nil, utils.WrapError(err, "获取评论列表失败")
	}
	log.Printf("获取评论列表成功: PostID=%d, Count=%d", postID, len(comments))
	return comments, nil
}

func (s *CommentService) GetCommentByID(id int64, postID int64) (*models.Comment, error) {
	var comment models.Comment
	err := s.db.Preload("User").Where("id = ? AND post_id = ?", id, postID).First(&comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("评论不存在: CommentID=%d, PostID=%d", id, postID)
			return nil, utils.ErrCommentNotFound
		}
		log.Printf("获取评论失败: CommentID=%d, Error=%v", id, err)
		return nil, utils.WrapError(err, "获取评论失败")
	}
	return &comment, nil
}

func (s *CommentService) UpdateComment(id int64, postID int64, userID int64, content string) error {
	comment := models.Comment{}
	result := s.db.Where("id = ? AND post_id = ?", id, postID).First(&comment)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("评论不存在或无权更新: CommentID=%d, PostID=%d, UserID=%d", id, postID, userID)
			return utils.ErrPermissionDenied
		}
		log.Printf("查询评论失败: CommentID=%d, Error=%v", id, result.Error)
		return utils.WrapError(result.Error, "查询评论失败")
	}

	if comment.UserID != userID {
		log.Printf("无权更新评论: CommentID=%d, Owner=%d, UserID=%d", id, comment.UserID, userID)
		return utils.ErrPermissionDenied
	}

	if err := s.db.Model(&comment).Update("content", content).Error; err != nil {
		log.Printf("更新评论失败: CommentID=%d, Error=%v", id, err)
		return utils.WrapError(err, "更新评论失败")
	}

	log.Printf("评论更新成功: CommentID=%d, UserID=%d", id, userID)
	return nil
}

func (s *CommentService) DeleteComment(id int64, postID int64, userID int64) error {
	comment := models.Comment{}
	result := s.db.Where("id = ? AND post_id = ?", id, postID).First(&comment)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("评论不存在或无权删除: CommentID=%d, PostID=%d, UserID=%d", id, postID, userID)
			return utils.ErrPermissionDenied
		}
		log.Printf("查询评论失败: CommentID=%d, Error=%v", id, result.Error)
		return utils.WrapError(result.Error, "查询评论失败")
	}

	if comment.UserID != userID {
		log.Printf("无权删除评论: CommentID=%d, Owner=%d, UserID=%d", id, comment.UserID, userID)
		return utils.ErrPermissionDenied
	}

	if err := s.db.Delete(&comment).Error; err != nil {
		log.Printf("删除评论失败: CommentID=%d, Error=%v", id, err)
		return utils.WrapError(err, "删除评论失败")
	}

	log.Printf("评论删除成功: CommentID=%d, UserID=%d", id, userID)
	return nil
}
