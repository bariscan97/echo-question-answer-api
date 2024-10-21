package post

import (
	posts_model "articles-api/models/post"
	token_model "articles-api/models/token"
	service "articles-api/services/post"
	"articles-api/utils"
	"net/http"
	"strconv"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PostController struct {
	PostService service.IPostService
}

type IPostController interface {
	CreatePost(c echo.Context) error
	CreateComment(c echo.Context) error
	UpdatePostById(c echo.Context) error
	GetPostCommentsById(c echo.Context) error
	DeletePostById(c echo.Context) error
	GetSinglePostById(c echo.Context) error
	GetAllPost(c echo.Context) error
	GetUserPostsById(c echo.Context) error
	LikePostById(c echo.Context) error
}

func NewAuthController(postService service.IPostService) IPostController {
	return &PostController{
		PostService: postService,
	}
}

func (controller *PostController) GetCategories(c echo.Context) error {
	page := c.QueryParam("page")

	if page == "" {
		page = "0"
	}

	_, err := strconv.Atoi(page)

	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	result, err := controller.PostService.GetCategories(page)

	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"result": result,
	})

}

func (controller *PostController) LikePostById(c echo.Context) error {

	id := c.Param("id")

	post_id, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "unexpected error",
		})
	}

	command, err := controller.PostService.LikePostById(user.User.Id, post_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "unexpected error",
		})

	}

	return c.String(200, command)
}

func (controller *PostController) GetUserPostsById(c echo.Context) error {

	page := c.QueryParam("page")

	if page == "" {
		page = "0"
	}

	_, err := strconv.Atoi(page)

	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var current_userId *uuid.UUID

	id := c.Param("userid")

	parsedUUID, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	currentUser, ok := c.Get("user").(*token_model.Claim)

	if ok {
		current_userId = &currentUser.User.Id
	}

	posts, err := controller.PostService.GetUserPostsById(current_userId, parsedUUID, page)

	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": posts,
	})
}

func (controller *PostController) CreatePost(c echo.Context) error {
	var reqBody posts_model.CreatePostModel

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	validate := validator.New()

	if err := validate.Struct(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "unexpected error",
		})
	}

	post, err := controller.PostService.CreatePost(user.User.Id, &reqBody)

	if err != nil {
		return c.JSON(http.StatusNotImplemented, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": post,
	})
}

func (controller *PostController) GetPostCommentsById(c echo.Context) error {
	id := c.Param("id")

	page := c.QueryParam("page")

	if page == "" {
		page = "0"
	}

	_, err := strconv.Atoi(page)

	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var current_userId *uuid.UUID

	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		current_userId = nil
	}

	current_userId = &user.User.Id

	parsedUUID, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	posts, err := controller.PostService.GetPostCommentsById(current_userId, parsedUUID, page)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"data": posts,
	})

}

func (controller *PostController) GetSinglePostById(c echo.Context) error {

	id := c.Param("id")

	parsedUUID, err := uuid.Parse(id)

	if err != nil {

		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var current_userId *uuid.UUID

	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		current_userId = nil
	}

	current_userId = &user.User.Id

	post, err := controller.PostService.GetPostById(current_userId, parsedUUID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"data": post,
	})

}

func (controller *PostController) UpdatePostById(c echo.Context) error {

	var reqBody posts_model.UpdatePostModel

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	validate := validator.New()

	if err := validate.Struct(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	post_id := c.Param("id")

	parsedUUID, err := uuid.Parse(post_id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "unexpected error",
		})
	}

	if err := controller.PostService.UpdatePostById(parsedUUID, user.User.Id, utils.StructToMap(reqBody)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "update succesful",
	})

}

func (controller *PostController) DeletePostById(c echo.Context) error {

	post_id := c.Param("id")

	parsedUUID, err := uuid.Parse(post_id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "unexpected error",
		})
	}

	if err := controller.PostService.DeletePostById(user.User.Id, parsedUUID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "succesful",
	})

}

func (controller *PostController) CreateComment(c echo.Context) error {
	var reqBody posts_model.CreatePostModel

	parent_id := c.Param("parentid")

	parsedUUID, err := uuid.Parse(parent_id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	validate := validator.New()

	if err := validate.Struct(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "unexpected error",
		})
	}

	post, err := controller.PostService.CreateComment(user.User.Id, parsedUUID, &reqBody)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"data": post,
	})
}

func (controller *PostController) GetAllPost(c echo.Context) error {

	var cat *string

	s := c.Param("category")

	flag := true

	for _, r := range s {
		if !unicode.IsLetter(r) {
			flag = false
			break
		}
	}

	if flag && s != "" {
		cat = &s
	}

	page := c.QueryParam("page")

	if page == "" {
		page = "0"
	}

	_, err := strconv.Atoi(page)

	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var current_userId *uuid.UUID

	user, ok := c.Get("user").(*token_model.Claim)

	if ok {
		current_userId = &user.User.Id
	}

	posts, err := controller.PostService.GetAllPost(current_userId, page, cat)

	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": posts,
	})
}
