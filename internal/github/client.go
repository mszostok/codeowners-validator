package github

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"

	"go.szostok.io/codeowners/internal/config"
	"go.szostok.io/codeowners/pkg/url"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

var httpRequestTimeout int = 30

// Validate validates if provided client options are valid.
func Validate(cfg *config.Config) error {
	if cfg.GithubAccessToken == "" && cfg.GithubAppID == 0 {
		return errors.New("GitHub authorization is required, provide ACCESS_TOKEN or APP_ID")
	}

	if cfg.GithubAccessToken != "" && cfg.GithubAppID != 0 {
		return errors.New("GitHub ACCESS_TOKEN cannot be provided when APP_ID is specified")
	}

	if cfg.GithubAppID != 0 {
		if cfg.GithubAppInstallationID == 0 {
			return errors.New("GitHub APP_INSTALLATION_ID is required with APP_ID")
		}
		if cfg.GithubAppPrivateKey == "" {
			return errors.New("GitHub APP_PRIVATE_KEY is required with APP_ID")
		}
	}

	return nil
}

func NewClient(ctx context.Context, cfg *config.Config) (ghClient *github.Client, isApp bool, err error) {
	if err := Validate(cfg); err != nil {
		return nil, false, err
	}

	httpClient := &http.Client{
		Transport:     http.DefaultClient.Transport,
		CheckRedirect: http.DefaultClient.CheckRedirect,
		Jar:           http.DefaultClient.Jar,
		Timeout:       time.Duration(httpRequestTimeout) * time.Second,
	}

	if cfg.GithubAccessToken != "" {
		httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.GithubAccessToken},
		))
	} else if cfg.GithubAppID != 0 {
		httpClient, err = createAppInstallationHTTPClient(cfg)
		isApp = true
		if err != nil {
			return
		}
	}

	baseURL, uploadURL := cfg.GithubBaseURL, cfg.GithubUploadURL

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

func createAppInstallationHTTPClient(cfg *config.Config) (client *http.Client, err error) {
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, cfg.GithubAppID, cfg.GithubAppInstallationID, []byte(cfg.GithubAppPrivateKey))
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: itr}, nil
}
