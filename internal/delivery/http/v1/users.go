package v1

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhosyaaa/RoommateTap/internal/domain"
	"github.com/zhosyaaa/RoommateTap/internal/service"
	"io/ioutil"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)

		users.GET("/google_login", h.GoogleLogin)
		users.GET("/google_callback", h.GoogleCallback)

		users.POST("/auth/refresh", h.userRefresh)
		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.GET("/verify/:code", h.userVerify)
			//
			//	schools := authenticated.Group("/schools/")
			//	{
			//		schools.POST("", h.userCreateSchool)
			//		schools.GET("", h.userGetSchools)
			//		schools.GET("/:id", h.userGetSchoolById)
			//		schools.PUT("/:id", h.userUpdateSchool)
			//	}
		}
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

	ip := c.ClientIP()
	fingerprint := c.GetHeader("X-Fingerprint")
	tokens, sessionID, err := h.services.Users.SignUp(c.Request.Context(), service.UserSignUpInput{
		Username: inp.Name,
		Email:    inp.Email,
		Phone:    inp.Phone,
		Password: inp.Password,
	}, fingerprint, ip)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}
	c.SetCookie("sessionID", sessionID, 3600, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{"token": tokens})
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
	ip := c.ClientIP()
	fingerprint := c.GetHeader("X-Fingerprint")
	res, sessionID, err := h.services.Users.SignIn(c.Request.Context(), service.UserSignInInput{
		Email:    inp.Email,
		Password: inp.Password,
	}, fingerprint, ip)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}
	c.SetCookie("sessionID", sessionID, 3600, "/", "", false, true)
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
	sessionID, err := c.Cookie("sessionID")
	//fingerprint := c.GetHeader("X-Fingerprint")
	fingerprint := "testttt"
	res, sessionID, err := h.services.Users.RefreshTokens(sessionID, inp.Token, fingerprint)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.SetCookie("sessionID", sessionID, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) userVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newResponse(c, http.StatusBadRequest, "code is empty")

		return
	}

	id, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Users.Verify(c.Request.Context(), id, code); err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{"success"})
}
func (h *Handler) GoogleLogin(c *gin.Context) {
	url := h.cfg.AuthCodeURL("randomstate")

	c.Redirect(http.StatusSeeOther, url)
}
func (h *Handler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "randomstate" {
		c.String(http.StatusBadRequest, "States don't Match!!")
		return
	}

	code := c.Query("code")

	googlecon := h.cfg

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Code-Token Exchange Failed")
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "User Data Fetch Failed")
		return
	}

	defer resp.Body.Close()
	userData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "JSON Parsing Failed")
		return
	}

	var googleUser service.GoogleUser
	if err := json.Unmarshal(userData, &googleUser); err != nil {
		c.String(http.StatusInternalServerError, "User Data Unmarshal Failed")
		return
	}
	_, _, err = h.services.Users.OAuthSignIn(c, googleUser, "", "")
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	c.String(http.StatusOK, string(userData))
}
