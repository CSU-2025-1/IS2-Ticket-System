package migrator

import (
	"auth-service/pkg/hydra"
	"context"
	"encoding/json"
	"fmt"
	oryclient "github.com/ory/hydra-client-go"
	"io"
	"os"
	"strings"
)

type OAuthClient struct {
	ID                      string   `json:"id"`
	Name                    string   `json:"name"`
	Secret                  string   `json:"secret"`
	Scopes                  []string `json:"scopes"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	RedirectURIs            []string `json:"redirect_uris"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
}

func MigrateHydra(ctx context.Context, cfg hydra.Config) error {
	file, err := os.Open("./hydra/clients.json")
	if err != nil {
		return fmt.Errorf("open clients.json: %w", err)
	}
	defer file.Close()

	in, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read clients.json: %w", err)
	}

	var oauthClients []OAuthClient
	err = json.Unmarshal(in, &oauthClients)
	if err != nil {
		return fmt.Errorf("unmarshal clients.json: %w", err)
	}

	oryClients := make([]oryclient.OAuth2Client, len(oauthClients))
	for i, oauthClient := range oauthClients {
		scopes := strings.Join(oauthClient.Scopes, " ")
		oryClients[i] = oryclient.OAuth2Client{
			ClientId:                &oauthClient.ID,
			ClientName:              &oauthClient.Name,
			ClientSecret:            &oauthClient.Secret,
			GrantTypes:              oauthClient.GrantTypes,
			RedirectUris:            oauthClient.RedirectURIs,
			ResponseTypes:           oauthClient.ResponseTypes,
			Scope:                   &scopes,
			TokenEndpointAuthMethod: &oauthClient.TokenEndpointAuthMethod,
		}
	}

	client := hydra.New(ctx, &cfg)

	err = client.CreateOrUpdateOauthClient(ctx, oryClients...)
	if err != nil {
		return fmt.Errorf("migrate clients.json: %w", err)
	}

	return nil
}
