package http_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test_auth/application/user"
)

type AuthHandler struct {
	RegisterUC *user.RegisterUser
	LoginUC    *user.LoginUser
}

// Methods for AuthHandler to register routes
// - Register: handles user registration
// - Login: handles user login and JWT issuance
// - Profile: retrieves user profile information (not implemented here)
// - Logout: handles user logout (not implemented here)

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct{ Username, Password string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	theUser, err := h.RegisterUC.Execute(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": theUser.ID})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct{ Username, Password string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.LoginUC.Execute(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
