package user

import (
	service "articles-api/services/user"
	"articles-api/utils"
	token_model "articles-api/models/token"
	model "articles-api/models/user"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserService service.IUserService
}

type IUserController interface {
	UpdateUserById(c echo.Context) error
	DeleteUserById(c echo.Context) error
	GetUserById(c echo.Context) error
	GetMe(c echo.Context) error
	GetAllUsers(c echo.Context) error
	UpdateUserProfileImgById(c echo.Context) error
	GetMyBlockList(c echo.Context) error
	GetUserFollowingListById(c echo.Context) error
	GetUserFollowersListById(c echo.Context) error
	FollowUserById(c echo.Context) error
	BlockUserById(c echo.Context) error
}

func NewUserController(userService service.IUserService) IUserController {
	return &UserController{
		UserService: userService,
	}
}

func (controller *UserController) GetMyBlockList(c echo.Context) error {
	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "user not found",
		})
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

	result, err := controller.UserService.GetMyBlockList(user.User.Id, page)

	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"result": result,
	})
}

func (controller *UserController) GetUserFollowingListById(c echo.Context) error {
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

	id := c.Param("id")

	parsedUUID, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var current_userId *uuid.UUID

	user, ok := c.Get("user").(*token_model.Claim)

	if ok {
		current_userId = &user.User.Id
	}

	result, err := controller.UserService.GetUserFollowingListById(current_userId, parsedUUID, page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"result": result,
		})
	}

	return c.JSON(200, result)
}

func (controller *UserController) GetUserFollowersListById(c echo.Context) error {
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

	id := c.Param("id")

	parsedUUID, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var current_userId *uuid.UUID

	user, ok := c.Get("user").(*token_model.Claim)

	if ok {
		current_userId = &user.User.Id
	}

	result, err := controller.UserService.GetUserFollowersListById(current_userId, parsedUUID, page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"result": result,
		})
	}

	return c.JSON(200, result)
}

func (controller *UserController) FollowUserById(c echo.Context) error {
	id := c.Param("id")

	parsedUUID, err := uuid.Parse(id)

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

	command, err := controller.UserService.FollowUserById(user.User.Id, parsedUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.String(200, command)
}

func (controller *UserController) BlockUserById(c echo.Context) error {

	id := c.Param("id")

	parsedUUID, err := uuid.Parse(id)

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

	command, err := controller.UserService.BlockUserById(user.User.Id, parsedUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.String(200, command)
}

func (controller *UserController) GetAllUsers(c echo.Context) error {

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

	users, err := controller.UserService.GetAllUsers(current_userId, page)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"data": users,
	})
}

func (controller *UserController) UpdateUserProfileImgById(c echo.Context) error {

	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "unexpected error",
		})
	}

	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Image not found in request",
			"error":   err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Failed to open file",
			"error":   err.Error(),
		})
	}
	defer src.Close()

	uploadResult, err := cld.Upload.Upload(c.Request().Context(), src, uploader.UploadParams{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to upload image",
			"error":   err.Error(),
		})
	}

	if err := controller.UserService.UpdateUserById(user.User.Id, map[string]interface{}{
		"profile_img": uploadResult.SecureURL,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to upload image",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Image uploaded successfully",
		"url":     uploadResult.SecureURL,
	})
}

func (controller *UserController) GetMe(c echo.Context) error {
	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "user not found",
		})
	}

	return c.JSON(200, echo.Map{
		"user": user,
	})
}

func (controller *UserController) GetUserById(c echo.Context) error {

	id := c.Param("id")

	parsedUUID, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var current_userId *uuid.UUID

	user, ok := c.Get("user").(*token_model.Claim)

	if ok {
		current_userId = &user.User.Id
	}

	result, err := controller.UserService.GetUserById(current_userId, parsedUUID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"data": result,
	})
}

func (controller *UserController) UpdateUserById(c echo.Context) error {
	var reqBody model.UpdateUserModel

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

	if reqBody.Password != "" {

		hashedPassword, err := utils.HashPassword(reqBody.Password)

		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}

		reqBody.Password = hashedPassword

	}

	user, ok := c.Get("user").(*token_model.Claim)

	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "unexpected error",
		})
	}

	if err := controller.UserService.UpdateUserById(user.User.Id, utils.StructToMap(reqBody)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "update succesful",
	})

}

func (controller *UserController) DeleteUserById(c echo.Context) error {
	id := c.Param("id")

	parsedUUID, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := controller.UserService.DeleteMe(parsedUUID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "delete ok",
	})
}
