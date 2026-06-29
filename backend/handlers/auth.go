package handlers

import (
	"net/http"

	"backend/backend/services"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
	
}

func (h * AuthHandler) Register(*gin.Context){
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); 
	err!=nil{
		c.JSON(http.StatusBadRequest , gin.H{"error": err.Error()})
		return
	}
	response,err := h.authService.Register(input)
	if err != nil{
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return 
	}
	c.JSON(http.StatusCreated,response)
}    
func (h *AuthHandler) login(c *gin.Context){
	var input services.LoginInput 
	if err := c.shouldBindJSON(&input);
	err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return 
	}
	response , err := h.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return 
	}
	c.JSON(http.StatusOk,response)

	
}

func (h *AuthHandler) Me(c *gin.Context) {
    
    user, _ := c.Get("user")
    c.JSON(http.StatusOK, user)
}