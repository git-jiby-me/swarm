package client

import (
	"encoding/json"
	"net/url"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	timetypes "github.com/docker/engine-api/types/time"
)

// Events returns a stream of events in the daemon. It's up to the caller to close the stream
// by cancelling the context. Once the stream has been completely read an io.EOF error will
// be sent over the error channel. If an error is sent all processing will be stopped. It's up
// to the caller to reopen the stream in the event of an error by reinvoking this method.
func (cli *Client) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {

	messages := make(chan events.Message)
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		query, err := buildEventsQueryParams(cli.version, options)
		if err != nil {
			errs <- err
			return
		}

		resp, err := cli.get(ctx, "/events", query, nil)
		if err != nil {
			errs <- err
			return
		}
		defer resp.body.Close()

		decoder := json.NewDecoder(resp.body)

		for {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
				var event events.Message
				if err := decoder.Decode(&event); err != nil {
					errs <- err
					return
				}

				select {
				case messages <- event:
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				}
			}
		}
	}()

	return messages, errs
}

func buildEventsQueryParams(cliVersion string, options types.EventsOptions) (url.Values, error) {
	query := url.Values{}
	ref := time.Now()

	if options.Since != "" {
		ts, err := timetypes.GetTimestamp(options.Since, ref)
		if err != nil {
			return nil, err
		}
		query.Set("since", ts)
	}

	if options.Until != "" {
		ts, err := timetypes.GetTimestamp(options.Until, ref)
		if err != nil {
			return nil, err
		}
		query.Set("until", ts)
	}

	if options.Filters.Len() > 0 {
		filterJSON, err := filters.ToParamWithVersion(cliVersion, options.Filters)
		if err != nil {
			return nil, err
		}
		query.Set("filters", filterJSON)
	}

	return query, nil
}
