package interactive

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
)

func GetConfigFilePath(ctx context.Context) (string, error) {
	prompt := &promptui.Prompt{
		Label:   "Config File Path",
		Default: "changes.yaml",
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get config file path: %w", err)
	}

	return result, nil
}
