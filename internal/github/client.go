package github

import (
	"context"
	"net/http"
	"time"

	"github.com/mszostok/codeowners-validator/pkg/url"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

type ClientConfig struct {
	AccessToken        string        `envconfig:"optional"`
	BaseURL            string        `envconfig:"optional"`
	UploadURL          string        `envconfig:"optional"`
	HTTPRequestTimeout time.Duration `envconfig:"default=30s"`
}

func NewClient(ctx context.Context, cfg ClientConfig) (ghClient *github.Client, err error) {
	httpClient := http.DefaultClient
	if cfg.AccessToken != "" {
		httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.AccessToken},
		))
	}
	httpClient.Timeout = cfg.HTTPRequestTimeout

	baseURL, uploadURL := cfg.BaseURL, cfg.UploadURL

	if baseURL == "" {
		return github.NewClient(httpClient), nil
	}

	if uploadURL == "" { // often the baseURL are same as the uploadURL, so we do not require to provide both of them
		uploadURL = baseURL
	}

	bURL, uURL := url.CanonicalPath(baseURL), url.CanonicalPath(uploadURL)
	return github.NewEnterpriseClient(bURL, uURL, httpClient)
}
