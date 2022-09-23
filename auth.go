package guku

import (
	"context"
	"errors"
	"fmt"
	"time"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"
	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (c *Client) login(region string, poolId string, clientId string, username string, password string) error {
	// configure cognito srp
	csrp, err := cognitosrp.NewCognitoSRP(username, password, poolId, clientId, nil)
	if err != nil {
		panic(err)
	}

	// configure cognito identity provider
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
	)
	if err != nil {
		panic(err)
	}
	svc := cip.NewFromConfig(cfg)

	// initiate auth
	resp, err := svc.InitiateAuth(context.Background(), &cip.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserSrpAuth,
		ClientId:       aws.String(csrp.GetClientId()),
		AuthParameters: csrp.GetAuthParams(),
	})
	if err != nil {
		panic(err)
	}

	// respond to password verifier challenge
	if resp.ChallengeName == types.ChallengeNameTypePasswordVerifier {
		challengeResponses, _ := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())

		resp, err := svc.RespondToAuthChallenge(context.Background(), &cip.RespondToAuthChallengeInput{
			ChallengeName:      types.ChallengeNameTypePasswordVerifier,
			ChallengeResponses: challengeResponses,
			ClientId:           aws.String(csrp.GetClientId()),
		})
		if err != nil {
			return err
		}

		result := *resp.AuthenticationResult

		claims := jwt.MapClaims{}
		_, _, err = new(jwt.Parser).ParseUnverified(*result.IdToken, claims)
		if err != nil {
			return err
		}

		groupsClaim, ok := claims["cognito:groups"]
		if !ok {
			return errors.New("No groups claim found")
		}

		groupsArray, ok := groupsClaim.([]string)
		if !ok || len(groupsArray) == 0 {
			return errors.New("Empty groups claim")
		}

		c.accountId = &groupsArray[0]
		c.accessToken = result.AccessToken
		c.idToken = result.IdToken
		c.refreshToken = result.RefreshToken
		return nil
	}

	return errors.New(fmt.Sprintf("Unexpected challenge type %s", resp.ChallengeName))
}
