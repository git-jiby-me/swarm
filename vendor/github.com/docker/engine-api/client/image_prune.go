package client

import (
	"encoding/json"
	"fmt"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ImagesPrune requests the daemon to delete unused data
func (cli *Client) ImagesPrune(ctx context.Context, cfg types.ImagesPruneConfig) (types.ImagesPruneReport, error) {
	var report types.ImagesPruneReport

	if err := cli.NewVersionError("1.25", "image prune"); err != nil {
		return report, err
	}

	serverResp, err := cli.post(ctx, "/images/prune", nil, cfg, nil)
	if err != nil {
		return report, err
	}
	defer ensureReaderClosed(serverResp)

	if err := json.NewDecoder(serverResp.body).Decode(&report); err != nil {
		return report, fmt.Errorf("Error retrieving disk usage: %v", err)
	}

	return report, nil
}
