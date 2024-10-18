package user

import (
	"articles-api/database/repository/user"
	model "articles-api/models/user"

	"github.com/google/uuid"
)

type UserService struct {
	UserRepository user.IUserRepository
}

type IUserService interface {
	CreateUser(data *model.RegisterUserModel) (map[string]interface{}, error)
	GetUserByEmail(email string) (*model.FetchUserModel, error)
	DeleteMe(current_userId uuid.UUID) error
    GetAllUsers(current_userId *uuid.UUID, page string) ([]model.FetchUserModel, error)
    GetUserById(current_userId *uuid.UUID, id uuid.UUID) (*model.FetchUserModel, error)
    UpdateUserById(current_userId uuid.UUID, data map[string]interface{}) error
    GetMyBlockList(current_userId uuid.UUID, page string) ([]map[string]interface{}, error)
	GetUserFollowingListById(current_userId *uuid.UUID, id uuid.UUID, page string) ([]map[string]interface{}, error) 
	GetUserFollowersListById(current_userId *uuid.UUID, id uuid.UUID, page string) ([]map[string]interface{}, error)
	FollowUserById(current_userId uuid.UUID, id uuid.UUID) (string, error)
	BlockUserById(current_userId uuid.UUID, id uuid.UUID) (string , error)
}

func NewAuthService(userRepository user.IUserRepository) IUserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (userService *UserService) GetMyBlockList(current_userId uuid.UUID, page string) ([]map[string]interface{}, error)  {
	return userService.UserRepository.GetMyBlockList(current_userId, page)
}

func (userService *UserService) GetUserFollowingListById(current_userId *uuid.UUID, id uuid.UUID, page string) ([]map[string]interface{}, error)  {
	return userService.UserRepository.GetUserFollowingListById(current_userId, id ,page)
}

func (userService *UserService) GetUserFollowersListById(current_userId *uuid.UUID, id uuid.UUID, page string) ([]map[string]interface{}, error) {
	return userService.UserRepository.GetUserFollowersListById(current_userId, id, page)
}

func (userService *UserService) FollowUserById(current_userId uuid.UUID, id uuid.UUID) (string, error) {
	return userService.UserRepository.FollowUserById(current_userId, id)
}

func (userService *UserService) BlockUserById(current_userId uuid.UUID, id uuid.UUID) (string , error) {
	return userService.UserRepository.BlockUserById(current_userId, id)
}

func (userService *UserService) CreateUser(data *model.RegisterUserModel) (map[string]interface{}, error) {
	return userService.UserRepository.CreateUser(data)
}

func (userService *UserService) GetUserByEmail(email string) (*model.FetchUserModel, error) {
	return userService.UserRepository.GetUserByEmail(email)
}

func (userService *UserService) DeleteMe(current_userId uuid.UUID) error {
	return userService.UserRepository.DeleteMe(current_userId)
}

func (userService *UserService) GetAllUsers(current_userId *uuid.UUID, page string) ([]model.FetchUserModel, error) {
	return userService.UserRepository.GetAllUsers(current_userId, page)
}

func (userService *UserService) GetUserById(current_userId *uuid.UUID, id uuid.UUID) (*model.FetchUserModel, error) {
	return userService.UserRepository.GetUserById(current_userId ,id)
}

func (userService *UserService) UpdateUserById(current_userId uuid.UUID, data map[string]interface{}) error {
	return userService.UserRepository.UpdateUserById(current_userId, data)
}
