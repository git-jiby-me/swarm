package client

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/registry"
	"golang.org/x/net/context"
)

// RegistryLogin authenticates the docker server with a given docker registry.
// It returns UnauthorizerError when the authentication fails.
func (cli *Client) RegistryLogin(ctx context.Context, auth types.AuthConfig) (registry.AuthenticateOKBody, error) {
	resp, err := cli.post(ctx, "/auth", url.Values{}, auth, nil)

	if resp.statusCode == http.StatusUnauthorized {
		return registry.AuthenticateOKBody{}, unauthorizedError{err}
	}
	if err != nil {
		return registry.AuthenticateOKBody{}, err
	}

	var response registry.AuthenticateOKBody
	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
}
