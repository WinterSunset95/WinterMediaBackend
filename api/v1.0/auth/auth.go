package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/WinterSunset95/WinterMediaBackend/cognito"
	"github.com/WinterSunset95/WinterMediaBackend/database"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type PostData struct {
	Email string
	Name string
	Pass string
	UserId string
}

type SigninData struct {
	Name string
	Pass string
}

type ConfirmData struct {
	Name string
	Code string
}

type TokenData struct {
	AccessToken string
	IdToken string
}

func getAuthState(jwkSet jwk.Set, ctx *gin.Context) (jwt.Token, error) {
	var tokenData TokenData
	err := ctx.ShouldBindJSON(&tokenData)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse([]byte(tokenData.IdToken), jwt.WithKeySet(jwkSet))
	if err != nil {
		return nil, err
	}
	return token, nil
}


func ApplyRoutes(r *gin.RouterGroup) {
	jwkSet, jwkError := jwk.Fetch(context.Background(), "https://cognito-idp.ap-south-1.amazonaws.com/ap-south-1_lLhgYR1Ao/.well-known/jwks.json")
	client := cognito.CognitoClient
	_ = client
	clientId := os.Getenv("COGNITO_CLIENT_ID")
	clientSec := os.Getenv("COGNITO_CLIENT_SECRET")
	auth := r.Group("/auth")
	{
		auth.POST("/signup", func(ctx *gin.Context) {
			var postData PostData
			err := ctx.ShouldBindJSON(&postData)
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": "Error parsing POST data",
				})
				return
			}
			fmt.Println(postData)
			if postData.Name == "" || postData.Email == "" {
				ctx.JSON(200, gin.H{
					"error": "name and email should not be empty",
				})
				return
			}
			mac := hmac.New(sha256.New, []byte(clientSec))
			mac.Write([]byte(postData.Name + clientId))
			secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
			res, err := client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
				ClientId: aws.String(clientId),
				SecretHash: aws.String(secretHash),
				Password: aws.String(postData.Pass),
				Username: aws.String(postData.Name),
				UserAttributes: []types.AttributeType{
					{Name: aws.String("email"), Value: aws.String(postData.Email)},
					{Name: aws.String("nickname"), Value: aws.String(postData.Name)},
				},
			})
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": err.Error(),
				})
				return
			}
			// Add user to the Users table
			query := "insert into Users (user_id, name, email) values ('" + postData.Name + "', '" + postData.Name + "', '" + postData.Email + "')"
			db := database.DB
			_, err = db.Exec(query)
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": "Error adding user to database: " + err.Error(),
				})
				return
			}
			ctx.JSON(200, gin.H{
				"success": res,
			})
		})

		auth.POST("/confirm", func(ctx *gin.Context) {
			var codeData ConfirmData
			err := ctx.ShouldBindJSON(&codeData)
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": "Error getting post data",
				})
				return
			}
			name := codeData.Name
			code := codeData.Code
			if code == "" {
				ctx.JSON(200, gin.H{
					"error": "Code is empty",
				})
				return
			}
			mac := hmac.New(sha256.New, []byte(clientSec))
			mac.Write([]byte(name + clientId))
			secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
			res, err := client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
				ClientId: aws.String(clientId),
				SecretHash: aws.String(secretHash),
				Username: aws.String(name),
				ConfirmationCode: aws.String(code),
			})
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(200, gin.H{
				"success": res,
			})
		})

		auth.POST("/signin", func(ctx *gin.Context) {
			var signinData SigninData
			err := ctx.ShouldBindJSON(&signinData)
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": "Could not get POST data: " + err.Error(),
				})
				return
			}
			name := signinData.Name
			pass := signinData.Pass
			mac := hmac.New(sha256.New, []byte(clientSec))
			mac.Write([]byte(name + clientId))
			secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
			res, err := client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
				AuthFlow: types.AuthFlowTypeUserPasswordAuth,
				AuthParameters: map[string]string{
					"USERNAME": *aws.String(name),
					"PASSWORD": *aws.String(pass),
					"SECRET_HASH": *aws.String(secretHash),
				},
				ClientId: aws.String(clientId),
			})

			if err != nil {
				ctx.JSON(200, gin.H{
					"error": err.Error(),
				})
				return
			}

			ctx.JSON(200, gin.H{
				"success": res,
			})
		})

		auth.POST("/get", func(ctx *gin.Context) {
			if jwkError != nil {
				ctx.JSON(200, gin.H{
					"error": "Error fetching jwk: " + jwkError.Error(),
				})
				return
			}
			token, err := getAuthState(jwkSet, ctx)
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": "Error parsing token: " + err.Error(),
				})
			}
			ctx.JSON(200, gin.H{
				"success": token,
			})
		})
	}
}
