package auth

import (
	service "articles-api/services/user"
	"articles-api/utils"
	model "articles-api/models/user"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	UserService service.IUserService
}

type IAuthController interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

func NewAuthController(userService service.IUserService) IAuthController {
	return &AuthController{
		UserService: userService,
	}
}

func (controller *AuthController) Register(c echo.Context) error {
	var reqbody model.RegisterUserModel

	if err := c.Bind(&reqbody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	validate := validator.New()

	if err := validate.Struct(reqbody); err != nil {

		return c.JSON(http.StatusAccepted, echo.Map{
			"error": err.Error(),
		})
	}

	hashedPassword, err := utils.HashPassword(reqbody.Password)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	reqbody.Password = hashedPassword

	result, err := controller.UserService.CreateUser(&reqbody)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "succesful",
		"data":    result,
	})
}

func (controller *AuthController) Login(c echo.Context) error {
	var reqbody model.LoginUserModel

	if err := c.Bind(&reqbody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	validate := validator.New()

	if err := validate.Struct(reqbody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	user, err := controller.UserService.GetUserByEmail(reqbody.Email)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	if err := utils.CheckPassword(reqbody.Password, user.Password); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	token, err := utils.GenerateJwtToken(user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	cookie := new(http.Cookie)

	cookie.Name = "acces_token"
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(72 * time.Hour)

	c.SetCookie(cookie)

	return c.JSON(200, echo.Map{
		"username":    user.Username,
		"email":       user.Email,
		"acces_token": token,
	})
}
