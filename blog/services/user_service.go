package services

import (
	"homework04/models"
	"homework04/utils"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(username, password, email string) error {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		log.Printf("用户注册失败: 用户名 %s 已存在", username)
		return utils.ErrUserExists
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		log.Printf("用户注册失败: 邮箱 %s 已存在", email)
		return utils.ErrUserExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		return utils.WrapError(err, "密码加密失败")
	}

	// 创建用户
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	if err := s.db.Create(&user).Error; err != nil {
		log.Printf("创建用户失败: %v", err)
		return utils.WrapError(err, "创建用户失败")
	}

	log.Printf("用户创建成功: ID=%d, Username=%s", user.ID, username)
	return nil
}

func (s *UserService) Authenticate(username, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("用户登录失败: 用户名 %s 不存在", username)
			return nil, utils.ErrUserNotFound
		}
		log.Printf("查询用户失败: %v", err)
		return nil, utils.WrapError(err, "查询用户失败")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("用户登录失败: 用户名 %s 密码错误", username)
		return nil, utils.ErrInvalidPassword
	}

	log.Printf("用户登录成功: ID=%d, Username=%s", user.ID, username)
	return &user, nil
}
