package interactive

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
)

// GetOAuthClientID prompts the user for OAuth client ID interactively.
func GetOAuthClientID(ctx context.Context) (string, error) {
	clientIDPrompt := &promptui.Prompt{
		Label: "OAuth Client ID",
	}

	clientID, err := clientIDPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get client ID: %w", err)
	}

	return clientID, nil
}

// GetOAuthClientSecret prompts the user for OAuth client secret interactively.
func GetOAuthClientSecret(ctx context.Context) (string, error) {
	clientSecretPrompt := &promptui.Prompt{
		Label: "OAuth Client Secret",
		Mask:  '*',
	}

	clientSecret, err := clientSecretPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get client secret: %w", err)
	}

	return clientSecret, nil
}
