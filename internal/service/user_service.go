package service

import (
	"errors"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/models"
	"majestic-gondola/internal/repository"
)

type UserService interface {
	Get(id int) (*models.User, error)
	GetAll() ([]models.User, error)
	BulkCreate(users []*models.User) error
	Update(user *models.User) error
}

type userService struct {
	userRepository repository.UserRepository
	log            *slog.Logger
}

func NewUserService(userRepository repository.UserRepository, logger *slog.Logger) UserService {
	return &userService{log: logger.With("component", "user_service"), userRepository: userRepository}
}

// BulkCreate implements [UserService].
func (u *userService) BulkCreate(users []*models.User) error {
	err := u.userRepository.BulkCreate(users)

	if err != nil {
		u.log.Error("Failed to bulk create users", slog.Any("error", err), slog.Int("parsed_users", len(users)))
		return apperr.Internal(err)
	}

	u.log.Info("Bulk created users", slog.Int("count", len(users)))
	return nil
}

// Get implements [UserService].
func (u *userService) Get(id int) (*models.User, error) {
	user, err := u.userRepository.FindById(id)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.NotFound("User not found", err)
		}
		u.log.Error("Failed to fetch a user by id", slog.Any("error", err), slog.Int("id", id))
		return nil, apperr.Internal(err)
	}

	return user, nil
}

// GetAll implements [UserService].
func (u *userService) GetAll() ([]models.User, error) {
	users, err := u.userRepository.GetAll()

	if err != nil {
		u.log.Error("Failed to fetch users", slog.Any("error", err))
		return nil, apperr.Internal(err)
	}

	return users, nil
}

// Update implements [UserService].
func (u *userService) Update(user *models.User) error {
	err := u.userRepository.Update(user)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.NotFound("User not found", err)
		}
		u.log.Error("Failed to update the user", slog.Any("error", err), slog.Any("parsed_user", user))
		return apperr.Internal(err)
	}

	u.log.Info("Updated u user")
	return nil
}
