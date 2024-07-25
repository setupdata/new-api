package vertex_openai

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	relaycommon "one-api/relay/common"
	"strings"
	"sync"
	"time"
)

var accessTokenMap sync.Map

func GetAccessToken(json string) (string, error) {
	data, ok := accessTokenMap.Load(json)
	if ok {
		token := data.(oauth2.Token)
		timeUntilExpiry := time.Until(token.Expiry)
		if timeUntilExpiry >= 10*time.Minute {
			return token.AccessToken, nil
		}
	}
	creds, err := google.CredentialsFromJSON(context.Background(), []byte(json), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return "", err
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", err
	}
	accessTokenMap.Store(json, *token)
	return token.AccessToken, nil
}

func GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	// LOCATION europe-west1 or us-east5
	LOCATION := "us-central1"
	parts := strings.SplitN(info.ApiKey, "|", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid api key: %s", info.ApiKey)
	}
	projectId := strings.TrimSpace(parts[0])
	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1beta1/projects/%s/locations/%s/endpoints/openapi/chat/completions", LOCATION, projectId, LOCATION)
	return url, nil
}

func GetRedirectModel(requestModel string) (string, error) {
	if model, ok := ModelIdMap[requestModel]; ok {
		return model, nil
	}
	return "", errors.Errorf("model %s not found", requestModel)
}

func GetModelList() []string {
	var models []string
	for n := range ModelIdMap {
		models = append(models, n)
	}
	return models
}
