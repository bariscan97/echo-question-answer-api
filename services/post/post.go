package post

import (
	repository "articles-api/database/repository/post"
	model "articles-api/models/post"

	"github.com/google/uuid"
)

type PostService struct {
	PostRepository repository.IPostRepository
}

type IPostService interface {
	CreatePost(current_userId uuid.UUID, data *model.CreatePostModel) (*model.FetchPostModel, error)
	CreateComment(current_userId uuid.UUID, parent_id uuid.UUID, data *model.CreatePostModel) (*model.FetchPostModel, error)
	UpdatePostById(current_id uuid.UUID, post_id uuid.UUID, data map[string]interface{}) error
	GetPostCommentsById(current_userId *uuid.UUID, id uuid.UUID, page string) ([]model.FetchPostModel, error)
	DeletePostById(current_id uuid.UUID, post_id uuid.UUID) error
	GetPostById(current_userId *uuid.UUID, id uuid.UUID) (*model.FetchPostModel, error)
	GetAllPost(current_userId *uuid.UUID, page string, category *string) ([]model.FetchPostModel, error)
	GetUserPostsById(current_userId *uuid.UUID, user_id uuid.UUID, page string) ([]model.FetchPostModel, error)
	LikePostById(current_userId uuid.UUID, post_id uuid.UUID) (string, error)
	GetCategories(page string) ([]map[string]interface{}, error)
}

func NewPostService(postRepository repository.IPostRepository) IPostService {
	return &PostService{
		PostRepository: postRepository,
	}
}

func (postService *PostService) GetCategories(page string) ([]map[string]interface{}, error) {
	return postService.PostRepository.GetCategories(page)
}

func (postService *PostService) CreatePost(current_userId uuid.UUID, data *model.CreatePostModel) (*model.FetchPostModel, error) {
	return postService.PostRepository.CreatePost(current_userId, data)
}

func (postService *PostService) CreateComment(current_userId uuid.UUID, parent_id uuid.UUID, data *model.CreatePostModel) (*model.FetchPostModel, error) {
	return postService.PostRepository.CreateComment(current_userId, parent_id, data)
}

func (postService *PostService) UpdatePostById(current_id uuid.UUID, post_id uuid.UUID, data map[string]interface{}) error {
	return postService.PostRepository.UpdatePostById(current_id, post_id, data)
}

func (postService *PostService) GetPostCommentsById(current_id *uuid.UUID, post_id uuid.UUID, page string) ([]model.FetchPostModel, error) {
	return postService.PostRepository.GetPostCommentsById(current_id, post_id, page)
}

func (postService *PostService) DeletePostById(current_id uuid.UUID, post_id uuid.UUID) error {
	return postService.PostRepository.DeletePostById(current_id, post_id)
}

func (postService *PostService) GetPostById(current_userId *uuid.UUID, id uuid.UUID) (*model.FetchPostModel, error) {
	return postService.PostRepository.GetSinglePostById(current_userId, id)
}

func (postService *PostService) GetAllPost(current_userId *uuid.UUID, page string, category *string) ([]model.FetchPostModel, error) {
	return postService.PostRepository.GetAllPost(current_userId, page, category)
}

func (postService *PostService) GetUserPostsById(current_userId *uuid.UUID, user_id uuid.UUID, page string) ([]model.FetchPostModel, error) {
	return postService.PostRepository.GetUserPostsById(current_userId, user_id, page)
}

func (postService *PostService) LikePostById(current_userId uuid.UUID, post_id uuid.UUID) (string, error) {
	return postService.PostRepository.LikePostById(current_userId, post_id)
}
