package handler

import (
	"time"

	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService          service.UserService
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewUserHandler(userService service.UserService, accessTokenDuration time.Duration, refreshTokenDuration time.Duration) *UserHandler {
	return &UserHandler{
		userService:          userService,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func CreateCookie(name string, value string, duration time.Duration) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.HTTPOnly = true
	cookie.MaxAge = int(duration.Seconds())
	cookie.SameSite = "Lax"
	cookie.Path = "/"

	// todo
	// set cookie.SameSite = "Strict" in production
	// set cookie.Domain in production
	// set cookie.Secure = true in production

	return cookie
}

func (u *UserHandler) Register(c fiber.Ctx) error {
	var req request.UserAuthRequest
	err := c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	err = u.userService.Register(c.Context(), req)
	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (u *UserHandler) Login(c fiber.Ctx) error {
	var req request.UserAuthRequest
	err := c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	resp, accessToken, refreshToken, err := u.userService.Login(c.Context(), req)

	if err != nil {
		return HandleServiceError(c, err)
	}

	c.Cookie(CreateCookie("access_token", accessToken, u.accessTokenDuration))
	c.Cookie(CreateCookie("refresh_token", refreshToken, u.refreshTokenDuration))

	return c.JSON(response.Response[*response.UserResponse]{Data: resp})
}

func (u *UserHandler) RefreshToken(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token", "")

	accessToken, refreshToken, err := u.userService.RefreshToken(c.Context(), refreshToken)

	if err != nil {
		return HandleServiceError(c, err)
	}

	c.Cookie(CreateCookie("access_token", accessToken, u.accessTokenDuration))
	c.Cookie(CreateCookie("refresh_token", refreshToken, u.refreshTokenDuration))

	return c.SendStatus(fiber.StatusNoContent)
}

func (u *UserHandler) Logout(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token", "")

	err := u.userService.Logout(c.Context(), refreshToken)
	if err != nil {
		return HandleServiceError(c, err)
	}

	c.Cookie(CreateCookie("access_token", "", -1*time.Second))
	c.Cookie(CreateCookie("refresh_token", "", -1*time.Second))

	return c.SendStatus(fiber.StatusNoContent)
}

func (u *UserHandler) SoftDelete(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	err := u.userService.SoftDelete(c.Context(), userID)
	if err != nil {
		return HandleServiceError(c, err)
	}

	c.Cookie(CreateCookie("access_token", "", -1*time.Second))
	c.Cookie(CreateCookie("refresh_token", "", -1*time.Second))

	return c.SendStatus(fiber.StatusNoContent)
}

func (u *UserHandler) UpdateAvatarURL(c fiber.Ctx) error {
	var req request.UserAvatarRequest
	err := c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	userID := c.Locals("user_id").(int)

	err = u.userService.UpdateAvatarURL(c.Context(), userID, req)
	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (u *UserHandler) UpdateUsername(c fiber.Ctx) error {
	var req request.UserUsernameRequest
	err := c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	userID := c.Locals("user_id").(int)

	err = u.userService.UpdateUsername(c.Context(), userID, req)
	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (u *UserHandler) GetByID(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	resp, err := u.userService.GetByID(c.Context(), userID)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.JSON(response.Response[*response.UserResponse]{Data: resp})
}

func (u *UserHandler) SearchByUsername(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	query := c.Query("q", "")

	resp, err := u.userService.SearchByUsername(c.Context(), userID, query)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.JSON(response.Response[[]response.UserResponse]{Data: resp})
}
