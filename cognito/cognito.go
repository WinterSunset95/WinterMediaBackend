package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var CognitoClient *cognitoidentityprovider.Client

func InitCognito() {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Error loading cognito config: " + err.Error())
	}
	CognitoClient = cognitoidentityprovider.NewFromConfig(sdkConfig)
}
