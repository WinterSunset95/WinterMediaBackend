package user

import (
	"context"
	"fmt"

	"github.com/WinterSunset95/WinterMediaBackend/common"
	"github.com/WinterSunset95/WinterMediaBackend/database"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type UserPost struct {
	Name string
	AccessToken string
	IdToken string
}

func ApplyRoutes(r *gin.RouterGroup) {
	jwkSet, jwkError := jwk.Fetch(context.Background(), "https://cognito-idp.ap-south-1.amazonaws.com/ap-south-1_lLhgYR1Ao/.well-known/jwks.json")
	users := r.Group("/user")
	{
		users.POST("/get", func(ctx *gin.Context) {
			if jwkError != nil {
				ctx.JSON(200, gin.H{
					"error": "Could not fetch JWK",
				})
				return
			}
			var userData UserPost
			err := ctx.ShouldBindJSON(&userData)
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": "Error parsing POST data",
				})
				return
			}
			_, err = jwt.Parse([]byte(userData.IdToken), jwt.WithKeySet(jwkSet))
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": "Could not parse token",
				})
				return
			}
			// User token is valid
			// Fetch user data from database
			var user common.User
			var db = database.DB
			fmt.Println(userData.Name)
			query := "select user_id, name, email, phone, created_at, picture from Users where user_id = '" + userData.Name + "'"
			rows, err := db.Query(query)
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": "Could not fetch user data: " + err.Error(),
				})
				return
			}
			for rows.Next() {
				err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.CreatedAt, &user.Picture)
				if err != nil {
					ctx.JSON(200, gin.H{
						"error": "Could not fetch rows: " + err.Error(),
					})
					return
				}
			}
			ctx.JSON(200, user)
		})
	}
}
