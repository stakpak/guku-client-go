package guku

import (
	"context"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

type Client struct {
	ctx           context.Context
	graphqlClient *graphql.Client
	accountId     *string
	accessToken   *string
	idToken       *string
	refreshToken  *string
}

func (c *Client) addGraphqlClient(url string) {
	client := graphql.NewClient(url, http.DefaultClient)
	c.graphqlClient = &client
}

func NewClient(ctx context.Context, url, username string, password string) (*Client, error) {
	client := Client{ctx: ctx}
	err := client.login(COGNITO_REGION, COGNITO_POOL_ID, COGNITO_CLIENT_ID, username, password)
	if err != nil {
		return nil, err
	}

	client.addGraphqlClient(url)
	return &client, nil
}

func (c *Client) GetCluster(id string) (*Cluster, error) {
	getResponse, err := getCluster(c.ctx, *c.graphqlClient, id)
	if err != nil {
		return nil, err
	}
	return getResponse.GetGetCluster(), nil
}

func (c *Client) CreateCluster(name string, server string, ca string, token string, apiVersion string, context string) (*ClusterCreate, error) {
	getResponse, err := createCluster(c.ctx, *c.graphqlClient, name, server, ca, token, apiVersion, context)
	if err != nil {
		return nil, err
	}
	return getResponse.GetCreateCluster(), nil
}

func (c *Client) CreatePrivateCluster(name string, token string, privateTunnelToken string, apiVersion string, context string) (*ClusterCreate, error) {
	getResponse, err := createPrivateCluster(c.ctx, *c.graphqlClient, name, token, privateTunnelToken, apiVersion, context)
	if err != nil {
		return nil, err
	}
	return (*ClusterCreate)(getResponse.GetCreatePrivateCluster()), nil
}

func (c *Client) UpdateCluster(id string, name *string, server *string, ca *string, token *string, apiVersion *string, context *string) (*ClusterUpdate, error) {
	getResponse, err := updateCluster(c.ctx, *c.graphqlClient, id, name, server, ca, token, apiVersion, context)
	if err != nil {
		return nil, err
	}
	return getResponse.GetUpdateCluster(), nil
}

func (c *Client) UpdatePrivateCluster(id string, name *string, token *string, privateTunnelToken *string, apiVersion *string, context *string) (*ClusterUpdate, error) {
	getResponse, err := updatePrivateCluster(c.ctx, *c.graphqlClient, id, name, token, privateTunnelToken, apiVersion, context)
	if err != nil {
		return nil, err
	}
	return (*ClusterUpdate)(getResponse.GetUpdatePrivateCluster()), nil
}

func (c *Client) DeleteCluster(id string) (*ClusterDelete, error) {
	getResponse, err := deleteCluster(c.ctx, *c.graphqlClient, id)
	if err != nil {
		return nil, err
	}
	return getResponse.GetDeleteCluster(), nil
}

func (c *Client) GetPlatfom(platformID string, platformVersion string) (*Platform, error) {
	getResponse, err := getPlatform(c.ctx, *c.graphqlClient, platformID, platformVersion)
	if err != nil {
		return nil, err
	}
	return getResponse.GetGetPlatform(), nil
}
