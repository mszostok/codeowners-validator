package github

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/mszostok/codeowners-validator/pkg/url"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

type ClientConfig struct {
	AccessToken string `envconfig:"optional"`

	AppID             int64  `envconfig:"optional"`
	AppPrivateKey     string `envconfig:"optional"`
	AppInstallationID int64  `envconfig:"optional"`

	BaseURL            string        `envconfig:"optional"`
	UploadURL          string        `envconfig:"optional"`
	HTTPRequestTimeout time.Duration `envconfig:"default=30s"`
}

func NewClient(ctx context.Context, cfg ClientConfig) (ghClient *github.Client, isApp bool, err error) {
	httpClient := http.DefaultClient

	if cfg.AccessToken != "" {
		httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.AccessToken},
		))
	} else if cfg.AppID != 0 {
		httpClient, err = createAppInstallationHttpClient(ctx, cfg)
		isApp = true
		if err != nil {
			return
		}
	}
	httpClient.Timeout = cfg.HTTPRequestTimeout

	baseURL, uploadURL := cfg.BaseURL, cfg.UploadURL

	if baseURL == "" {
		ghClient = github.NewClient(httpClient)
		return
	}

	if uploadURL == "" { // often the baseURL is same as the uploadURL, so we do not require to provide both of them
		uploadURL = baseURL
	}

	bURL, uURL := url.CanonicalPath(baseURL), url.CanonicalPath(uploadURL)
	ghClient, err = github.NewEnterpriseClient(bURL, uURL, httpClient)
	return
}

func createAppInstallationHttpClient(ctx context.Context, cfg ClientConfig) (client *http.Client, err error) {
	if cfg.AppInstallationID == 0 {
		return nil, errors.New("Github AppInstallationID is required with AppID")
	}
	if cfg.AppPrivateKey == "" {
		return nil, errors.New("Github AppPrivateKey is required with AppID")
	}

	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, cfg.AppID, cfg.AppInstallationID, []byte(cfg.AppPrivateKey))
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: itr}, nil
}
