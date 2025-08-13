package webserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/M2rk13/Otus-327619/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type loginCredentials struct {
	Login    string `json:"login" example:"admin"`
	Password string `json:"password" example:"secret"`
}

type jwtClaims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

// @Summary      User login
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      loginCredentials  true  "User login and password"
// @Success 	 200 		  {object}  object{token=string}
// @Failure 	 400 		  {object}  object{error=string}
// @Failure 	 401 		  {object}  object{error=string}
// @Failure 	 500 		  {object}  object{error=string}
// @Router       /auth/login [post]
func loginHandler(c *gin.Context) {
	var creds loginCredentials

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})

		return
	}

	if creds.Login != config.AdminCfg.Login || creds.Password != config.AdminCfg.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})

		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &jwtClaims{
		Login: creds.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AdminCfg.JwtKey))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})

			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be 'Bearer <token>'"})

			return
		}

		tokenString := parts[1]
		claims := &jwtClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.AdminCfg.JwtKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})

			return
		}

		c.Next()
	}
}
