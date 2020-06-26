package generic

import (
	"fmt"
	gohttp "net/http"
	"net/url"

	"github.com/cyberark/secretless-broker/pkg/secretless/log"
	"github.com/cyberark/secretless-broker/pkg/secretless/plugin/connector"
)

// Connector injects an HTTP requests with authorization headers.
type Connector struct {
	logger log.Logger
	config *config
}

// Connect implements the http.Connector func signature.
func (c *Connector) Connect(
	r *gohttp.Request,
	credentialsByID connector.CredentialValuesByID,
) error {
	// Validate credential values match expected patterns
	if err := c.config.validate(credentialsByID); err != nil {
		return err
	}

	// Fulfill SSL requests
	if c.config.ForceSSL {
		r.URL.Scheme = "https"
	}

	// Add configured headers to request
	headers, err := render(c.config.Headers, credentialsByID)
	if err != nil {
		return fmt.Errorf("failed to render headers: %s", err)
	}
	for headerName, headerVal := range headers {
		r.Header.Set(headerName, headerVal)
	}

	// Add configured params to request
	params, err := render(c.config.QueryParams, credentialsByID)
	if err != nil {
		return fmt.Errorf("failed to render query params: %s", err)
	}
	var query url.Values
	if len(r.URL.RawQuery) > 0 {
		query = r.URL.Query()
	} else {
		query = url.Values{}
	}
	for key, value := range params {
		query.Add(key, value)
	}
	r.URL.RawQuery = query.Encode()

	return nil
}
