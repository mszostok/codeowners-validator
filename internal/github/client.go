package github

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"

	"go.szostok.io/codeowners-validator/pkg/url"

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

// Validate validates if provided client options are valid.
func (c *ClientConfig) Validate() error {
	if c.AccessToken == "" && c.AppID == 0 {
		return errors.New("GitHub authorization is required, provide ACCESS_TOKEN or APP_ID")
	}

	if c.AccessToken != "" && c.AppID != 0 {
		return errors.New("GitHub ACCESS_TOKEN cannot be provided when APP_ID is specified")
	}

	if c.AppID != 0 {
		if c.AppInstallationID == 0 {
			return errors.New("GitHub APP_INSTALLATION_ID is required with APP_ID")
		}
		if c.AppPrivateKey == "" {
			return errors.New("GitHub APP_PRIVATE_KEY is required with APP_ID")
		}
	}

	return nil
}

func NewClient(ctx context.Context, cfg *ClientConfig) (ghClient *github.Client, isApp bool, err error) {
	if err := cfg.Validate(); err != nil {
		return nil, false, err
	}

	httpClient := http.DefaultClient

	if cfg.AccessToken != "" {
		httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.AccessToken},
		))
	} else if cfg.AppID != 0 {
		httpClient, err = createAppInstallationHTTPClient(cfg)
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

func createAppInstallationHTTPClient(cfg *ClientConfig) (client *http.Client, err error) {
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, cfg.AppID, cfg.AppInstallationID, []byte(cfg.AppPrivateKey))
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: itr}, nil
}
