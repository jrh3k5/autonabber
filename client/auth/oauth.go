package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/int128/oauth2cli"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

// DefaultOAuthServerPort is the default port for the OAuth server to listen for successful authentication
const DefaultOAuthServerPort = 54520

// Logger is the type of logger used during OAuth token retrieval
type Logger func(format string, args ...any)

// TokenOptions contains options for GetOAuthToken
type TokenOptions struct {
	logger     Logger
	serverPort int
}

// TokenOption is an option for GetOAuthToken
type TokenOption func(*TokenOptions)

// WithLogger sets the logger
func WithLogger(logger Logger) TokenOption {
	return func(opts *TokenOptions) {
		opts.logger = logger
	}
}

// WithOAuthServerPort sets the port on which the OAuth server should run
func WithOAuthServerPort(port int) TokenOption {
	return func(opts *TokenOptions) {
		opts.serverPort = port
	}
}

func GetOAuthToken(ctx context.Context, clientID string, clientSecret string, opts ...TokenOption) (*oauth2.Token, error) {
	tokenOptions := &TokenOptions{
		serverPort: DefaultOAuthServerPort,
		logger: func(string, ...any) {
			// deliberately no-op
		},
	}

	for _, opt := range opts {
		opt(tokenOptions)
	}

	localServerURLChan := make(chan string)
	defer close(localServerURLChan)

	cliConfig := oauth2cli.Config{
		OAuth2Config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://app.ynab.com/oauth/authorize",
				TokenURL: "https://api.ynab.com/oauth/token",
			},
		},
		LocalServerBindAddress: []string{fmt.Sprintf("127.0.0.1:%d", tokenOptions.serverPort)},
		LocalServerReadyChan:   localServerURLChan,
		Logf:                   tokenOptions.logger,
	}

	errGroup, errGroupCtx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		select {
		case url := <-localServerURLChan:
			log.Printf("Open %s", url)
			if err := browser.OpenURL(url); err != nil {
				tokenOptions.logger("could not open the browser: %s", err)
			}
			return nil
		case <-errGroupCtx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})

	tokenReturn := make(chan *oauth2.Token, 1)
	errGroup.Go(func() error {
		token, err := oauth2cli.GetToken(errGroupCtx, cliConfig)
		if err != nil {
			return fmt.Errorf("could not get a token: %w", err)
		}

		fmt.Println("Got token")

		tokenReturn <- token

		fmt.Println("Done")

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	close(tokenReturn)

	return <-tokenReturn, nil
}
