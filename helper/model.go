package helper

import (
	"retail-management/model/domain"
	"retail-management/model/web"
)

func ToUserRegisterResponse(user domain.User) web.UserRegisterResponse {
	return web.UserRegisterResponse{
		UserID:   user.UserID,
		Username: user.Username,
	}
}

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		UserID:   user.UserID,
		Username: user.Username,
	}
}

func ToUserResponses(users []domain.User) []web.UserResponse {
	userResponses := make([]web.UserResponse, 0)

	for _, user := range users {
		userResponses = append(userResponses, ToUserResponse(user))
	}
	return userResponses
}
