package args

import (
	"errors"
	"flag"
)

type Args struct {
	AccessToken string
}

func GetArgs() (*Args, error) {
	var accessToken string
	flag.StringVar(&accessToken, "access-token", "", "the access token to be used to communication with YNAB")

	flag.Parse()

	if accessToken == "" {
		return nil, errors.New("an access token must be provided")
	}

	return &Args{
		AccessToken: accessToken,
	}, nil
}
