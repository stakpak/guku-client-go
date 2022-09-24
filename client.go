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

type Transport struct {
	underlyingTransport http.RoundTripper
	token               string
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("authorization", t.token)
	return t.underlyingTransport.RoundTrip(req)
}

func (c *Client) addGraphqlClient(url string) {
	httpClient := &http.Client{
		Transport: &Transport{
			underlyingTransport: http.DefaultTransport,
			token:               *c.idToken,
		},
	}

	client := graphql.NewClient(url, httpClient)
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

func (c *Client) GetPlatform(platformID string, platformVersion string) (*Platform, error) {
	getResponse, err := getPlatform(c.ctx, *c.graphqlClient, platformID, platformVersion)
	if err != nil {
		return nil, err
	}
	return getResponse.GetGetPlatform(), nil
}

func (c *Client) GetPlatformBinding(clusterID string, platformBindingID string) (*PlatformBindingGet, error) {
	getResponse, err := getPlatformBinding(c.ctx, *c.graphqlClient, clusterID, platformBindingID)
	if err != nil {
		return nil, err
	}
	return getResponse.GetGetPlatformBinding(), nil
}

func (c *Client) CreatePlatformBinding(clusterID string, platformID string, platformVersion string, platformConfigID string) (*PlatformBindingCreate, error) {
	getResponse, err := createPlatformBinding(c.ctx, *c.graphqlClient, platformVersion, platformID, platformConfigID, clusterID)
	if err != nil {
		return nil, err
	}
	return (*PlatformBindingCreate)(getResponse.GetCreatePlatformBinding()), nil
}

func (c *Client) UpdatePlatformBinding(clusterID string, platformBindingID string, platformConfigID *string, platformVersion *string) (*PlatformBindingUpdate, error) {
	getResponse, err := updatePlatformBinding(c.ctx, *c.graphqlClient, clusterID, platformBindingID, platformConfigID, platformVersion)
	if err != nil {
		return nil, err
	}
	return getResponse.GetUpdatePlatformBinding(), nil
}

func (c *Client) DeletePlatformBinding(clusterID string, platformBindingID string) (*PlatformBindingDelete, error) {
	getResponse, err := deletePlatformBinding(c.ctx, *c.graphqlClient, clusterID, platformBindingID)
	if err != nil {
		return nil, err
	}
	return getResponse.GetDeletePlatformBinding(), nil
}
