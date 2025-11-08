package service

import (
	"fmt"
	"ipl-be-svc/internal/models"
	"ipl-be-svc/internal/repository"
	"ipl-be-svc/pkg/logger"
)

// UserService interface defines user service methods
type UserService interface {
	GetUserDetailByProfileID(profileID uint) (*models.UserDetail, error)
}

// userService implements UserService interface
type userService struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, logger *logger.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetUserDetailByProfileID gets user detail by profile ID
func (s *userService) GetUserDetailByProfileID(profileID uint) (*models.UserDetail, error) {
	if profileID == 0 {
		s.logger.WithField("profile_id", profileID).Error("Invalid profile ID")
		return nil, fmt.Errorf("invalid profile ID")
	}

	userDetail, err := s.userRepo.GetUserDetailByProfileID(profileID)
	if err != nil {
		s.logger.WithError(err).WithField("profile_id", profileID).Error("Failed to get user detail")
		return nil, err
	}

	s.logger.WithFields(map[string]interface{}{
		"profile_id": profileID,
		"user_id":    userDetail.UserID,
		"email":      userDetail.Email,
	}).Info("User detail retrieved successfully")

	return userDetail, nil
}
