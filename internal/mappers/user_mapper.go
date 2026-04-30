package mappers

import (
	"majestic-gondola/internal/dto"
	"majestic-gondola/internal/models"
	"time"
)

func ToUserResponse(u *models.User) dto.UserResponse {
	return dto.UserResponse{
		Id:        u.Id,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.Format(time.DateTime),
	}
}

func ToUserResponseList(us []models.User) []dto.UserResponse {
	responses := make([]dto.UserResponse, 0, len(us))
	for i := range us {
		responses = append(responses, ToUserResponse(&us[i]))
	}
	return responses
}

func CreateToUser(ur dto.CreateUserRequest) *models.User {
	return &models.User{
		Name: ur.Name,
	}
}

func CreateToUserList(urs []dto.CreateUserRequest) []*models.User {
	users := make([]*models.User, 0, len(urs))
	for i := range urs {
		user := CreateToUser(urs[i])
		users = append(users, user)
	}
	return users
}

func UpdateToUser(id int, ur dto.UpdateUserRequest) *models.User {
	return &models.User{
		Id:   id,
		Name: ur.Name,
	}
}
