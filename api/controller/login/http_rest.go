package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrokunz/canoe_reddit_integration/api/util"
	"github.com/pedrokunz/canoe_reddit_integration/internal/repository/login"
)

type Input struct {
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func Execute(repository login.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		var input Input
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		output, err := repository.UserRead(c, login.LoginRepositoryUserReadInput{
			Email:      input.Email,
			CustomerID: input.CustomerID,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		user := output.User
		if input.Password != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := util.GenerateToken(user.Customer.ID, user.Email, user.Roles, user.Attributes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
