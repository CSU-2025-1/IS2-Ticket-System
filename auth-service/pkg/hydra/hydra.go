package hydra

import (
	"context"
	"fmt"
	oryclient "github.com/ory/hydra-client-go"
	"log/slog"
)

type Client struct {
	hydra *oryclient.APIClient
}

func New(ctx context.Context, config *Config) *Client {
	configuration := oryclient.NewConfiguration()
	configuration.Servers = []oryclient.ServerConfiguration{
		{
			URL: config.AdminURL,
		},
	}
	configuration.Debug = config.Debug

	apiClient := oryclient.NewAPIClient(configuration)

	client := &Client{
		hydra: apiClient,
	}

	return client
}

func (c *Client) IntrospectOAuth2Token(ctx context.Context, token string) (*oryclient.OAuth2TokenIntrospection, error) {
	tokenInfo, resp, err := c.hydra.AdminApi.
		IntrospectOAuth2Token(ctx).
		Token(token).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("hydra admin: IntrospectOAuth2Token: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("hydra admin: IntrospectOAuth2Token returned status code %d", resp.StatusCode)
	}

	if tokenInfo == nil || tokenInfo.Sub == nil {
		return nil, ErrTokenExpired
	}

	slog.Debug("Token info", slog.String("info", fmt.Sprintf("%+v", tokenInfo)))

	return tokenInfo, nil
}

func (c *Client) AcceptLoginRequest(ctx context.Context, challenge, sub string) (*oryclient.CompletedRequest, error) {
	oryAcceptLoginResp, oryResp, err := c.hydra.AdminApi.AcceptLoginRequest(ctx).
		AcceptLoginRequest(oryclient.AcceptLoginRequest{
			Subject: sub,
		}).
		LoginChallenge(challenge).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("hydra admin: accept login request: %w", err)
	}

	slog.Debug("oryAcceptLoginReq",
		slog.String("oryAcceptLoginResp", fmt.Sprintf("%+v", oryAcceptLoginResp)),
		slog.String("oryResp", fmt.Sprintf("%+v", oryResp)),
	)

	return oryAcceptLoginResp, nil
}

func (c *Client) AcceptConsentRequest(ctx context.Context, challenge string, scopes []string) (*oryclient.CompletedRequest, error) {
	acr := oryclient.NewAcceptConsentRequest()
	acr.GrantScope = scopes

	oryAcceptConsentResp, oryResp, err := c.hydra.AdminApi.AcceptConsentRequest(ctx).
		ConsentChallenge(challenge).
		AcceptConsentRequest(*acr).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("hydra admin: accept consent request: %w", err)
	}

	slog.Debug("oryAcceptConsentResp",
		slog.String("oryAcceptConsentResp", fmt.Sprintf("%+v", oryAcceptConsentResp)),
		slog.String("oryResp", fmt.Sprintf("%+v", oryResp)),
	)

	return oryAcceptConsentResp, nil
}

func (c *Client) GetClientIDByLogin(ctx context.Context, challenge string) (*oryclient.OAuth2Client, error) {
	loginResp, resp, err := c.hydra.AdminApi.
		GetLoginRequest(ctx).
		LoginChallenge(challenge).
		Execute()
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("hydra admin: get client id: %w, status code: %d body: %s", err, resp.StatusCode, resp.Body)
		}
		return nil, fmt.Errorf("hydra admin: get client id: %w", err)
	}

	return &loginResp.Client, nil
}

func (c *Client) GetClientIDByConsent(ctx context.Context, challenge string) (*oryclient.OAuth2Client, error) {
	loginResp, resp, err := c.hydra.AdminApi.
		GetConsentRequest(ctx).
		ConsentChallenge(challenge).
		Execute()
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("hydra admin: get client id: %w, status code: %d body: %s", err, resp.StatusCode, resp.Body)
		}
		return nil, fmt.Errorf("hydra admin: get client id: %w", err)
	}

	return loginResp.Client, nil
}

func (c *Client) CreateOrUpdateOauthClient(ctx context.Context, clients ...oryclient.OAuth2Client) error {
	for _, client := range clients {
		_, _, err := c.hydra.AdminApi.GetOAuth2Client(ctx, *client.ClientId).Execute()
		if err != nil {
			_, _, err = c.hydra.AdminApi.CreateOAuth2Client(ctx).
				OAuth2Client(client).
				Execute()
			if err != nil {
				return fmt.Errorf("hydra admin: create oauth2 client: %w", err)
			}
		}

		_, _, err = c.hydra.AdminApi.UpdateOAuth2Client(ctx, *client.ClientId).
			OAuth2Client(client).
			Execute()
		if err != nil {
			return fmt.Errorf("hydra admin: update oauth2 client: %w", err)
		}
	}

	return nil
}
