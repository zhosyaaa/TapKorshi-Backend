package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhosyaaa/RoommateTap/internal/domain"
	"github.com/zhosyaaa/RoommateTap/internal/service"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/auth/refresh", h.userRefresh)

		//authenticated := users.Group("/", h.userIdentity)
		//{
		//	authenticated.POST("/verify/:code", h.userVerify)
		//
		//	schools := authenticated.Group("/schools/")
		//	{
		//		schools.POST("", h.userCreateSchool)
		//		schools.GET("", h.userGetSchools)
		//		schools.GET("/:id", h.userGetSchoolById)
		//		schools.PUT("/:id", h.userUpdateSchool)
		//	}
		//}
	}
}

type userSignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Phone    string `json:"phone" binding:"required,max=13"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func (h *Handler) userSignUp(c *gin.Context) {
	var inp userSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	if err := h.services.Users.SignUp(c.Request.Context(), service.UserSignUpInput{
		Username: inp.Name,
		Email:    inp.Email,
		Phone:    inp.Phone,
		Password: inp.Password,
	}); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusCreated)
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) userSignIn(c *gin.Context) {
	var inp signInInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Users.SignIn(c.Request.Context(), service.UserSignInInput{
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handler) userRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Users.RefreshTokens(c.Request.Context(), inp.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
