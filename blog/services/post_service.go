package services

import (
	"homework04/models"
	"homework04/utils"
	"log"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(title, content string, userID int64) error {
	post := models.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}
	if err := s.db.Create(&post).Error; err != nil {
		log.Printf("创建文章失败: UserID=%d, Error=%v", userID, err)
		return utils.WrapError(err, "创建文章失败")
	}
	log.Printf("文章创建成功: ID=%d, UserID=%d", post.ID, userID)
	return nil
}

func (s *PostService) GetPosts() ([]models.Post, error) {
	var posts []models.Post
	err := s.db.Preload("User").Find(&posts).Error
	if err != nil {
		log.Printf("获取文章列表失败: %v", err)
		return nil, utils.WrapError(err, "获取文章列表失败")
	}
	log.Printf("获取文章列表成功: Count=%d", len(posts))
	return posts, nil
}

func (s *PostService) GetPostByID(id int64) (*models.Post, error) {
	var post models.Post
	err := s.db.Preload("User").Where("id = ?", id).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("文章不存在: ID=%d", id)
			return nil, utils.ErrPostNotFound
		}
		log.Printf("获取文章失败: ID=%d, Error=%v", id, err)
		return nil, utils.WrapError(err, "获取文章失败")
	}
	return &post, nil
}

func (s *PostService) UpdatePost(id int64, title, content string, userID int64) error {
	post := models.Post{}
	result := s.db.Where("id = ?", id).First(&post)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("文章不存在或无权更新: PostID=%d, UserID=%d", id, userID)
			return utils.ErrPermissionDenied
		}
		log.Printf("查询文章失败: PostID=%d, Error=%v", id, result.Error)
		return utils.WrapError(result.Error, "查询文章失败")
	}

	if post.UserID != userID {
		log.Printf("无权更新文章: PostID=%d, PostOwner=%d, UserID=%d", id, post.UserID, userID)
		return utils.ErrPermissionDenied
	}

	if err := s.db.Model(&post).Updates(map[string]interface{}{
		"title":   title,
		"content": content,
	}).Error; err != nil {
		log.Printf("更新文章失败: PostID=%d, Error=%v", id, err)
		return utils.WrapError(err, "更新文章失败")
	}

	log.Printf("文章更新成功: ID=%d, UserID=%d", id, userID)
	return nil
}

func (s *PostService) DeletePost(id int64, userID int64) error {
	post := models.Post{}
	result := s.db.Where("id = ?", id).First(&post)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("文章不存在或无权删除: PostID=%d, UserID=%d", id, userID)
			return utils.ErrPermissionDenied
		}
		log.Printf("查询文章失败: PostID=%d, Error=%v", id, result.Error)
		return utils.WrapError(result.Error, "查询文章失败")
	}

	if post.UserID != userID {
		log.Printf("无权删除文章: PostID=%d, PostOwner=%d, UserID=%d", id, post.UserID, userID)
		return utils.ErrPermissionDenied
	}

	if err := s.db.Delete(&post).Error; err != nil {
		log.Printf("删除文章失败: PostID=%d, Error=%v", id, err)
		return utils.WrapError(err, "删除文章失败")
	}

	log.Printf("文章删除成功: ID=%d, UserID=%d", id, userID)
	return nil
}
