package args

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Args struct {
	PrintBudget           bool
	PrintHiddenCategories bool
	DryRun                bool
	OAuthServerPort       int
	ConfigFilePath        string
}

func GetArgs() (*Args, error) {
	args := &Args{
		OAuthServerPort: 54520,
		ConfigFilePath:  "changes.yaml",
	}

	for _, arg := range os.Args {
		switch {
		case strings.HasPrefix(arg, "--file="):
			args.ConfigFilePath = strings.TrimPrefix(arg, "--file=")
		case strings.HasPrefix(arg, "--oauth-server-port="):
			port, err := strconv.Atoi(strings.TrimPrefix(arg, "--oauth-server-port="))
			if err != nil {
				return nil, fmt.Errorf("invalid OAuth server port specified ('%s'): %w", arg, err)
			}

			args.OAuthServerPort = port
		case strings.HasPrefix(arg, "--dry-run"):
			if arg == "--dry-run" {
				args.DryRun = true
			} else {
				parsedBool, err := strconv.ParseBool(strings.TrimPrefix(arg, "--dry-run="))
				if err != nil {
					return nil, fmt.Errorf("invalid dry-run value specified ('%s'): %w", arg, err)
				}

				args.DryRun = parsedBool
			}
		case strings.HasPrefix(arg, "--print-budget="):
			parsedBool, err := strconv.ParseBool(strings.TrimPrefix(arg, "--print-budget="))
			if err != nil {
				return nil, fmt.Errorf("invalid print-budget value specified ('%s'): %w", arg, err)
			}

			args.PrintBudget = parsedBool
		}
	}

	return args, nil
}
