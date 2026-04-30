package repository

import (
	"errors"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/models"

	"github.com/go-pg/pg/v10"
)

type UserRepository interface {
	FindById(id int) (*models.User, error)
	GetAll() ([]models.User, error)
	BulkCreate(users []*models.User) error
	Update(user *models.User) error
}

type userRepository struct {
	db  *pg.DB
	log *slog.Logger
}

func NewUserRepository(db *pg.DB, logger *slog.Logger) UserRepository {
	return &userRepository{db: db, log: logger.With("component", "user_repository")}
}

// BulkCreate implements [UserRepository].
func (u *userRepository) BulkCreate(users []*models.User) error {
	_, err := u.db.Model(&users).Insert()

	if err != nil {
		return err
	}

	u.log.Info("Finished creating users")
	return nil
}

// FindById implements [UserRepository].
func (u *userRepository) FindById(id int) (*models.User, error) {
	user := new(models.User)
	err := u.db.Model(user).Where("user_id = ?", id).Select()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetAll implements [UserRepository].
func (u *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := u.db.Model(&users).Select()

	if err != nil {
		return nil, err
	}

	return users, nil
}

// Update implements [UserRepository].
func (u *userRepository) Update(user *models.User) error {
	res, err := u.db.Model(user).ExcludeColumn("created_at").WherePK().Update()

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	u.log.Info("Finished updating user")
	return nil
}
