package ntfy

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
)

// Service providing Ntfy as a notification service
type Service struct {
	standard.Standard
	config     *Config
	pkr        format.PropKeyResolver
	httpClient *http.Client
	client     jsonclient.Client
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) (err error) {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)
	if err = service.config.SetURL(configURL); err != nil {
		return
	}
	service.httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: service.config.DisableTLS,
			},
		},
		// Set a reasonable timeout to prevent one bad transfer from block all subsequent ones
		Timeout: 10 * time.Second,
	}
	service.client = jsonclient.NewWithHTTPClient(service.httpClient)
	// If username and password are provided, will add auth header in every request
	if len(service.config.Username) > 0 && len(service.config.Password) > 0 {
		auth := service.authorization()
		service.client.Headers().Set("Authorization", auth)
	}
	return
}

// Authorization returns the corresponding `Authorization` HTTP header value for the Token
func (service *Service) authorization() string {
	userPass := service.config.Username + ":" + service.config.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(userPass))
}

// Send a notification message to Ntfy
func (service *Service) Send(message string, params *types.Params) (err error) {
	cfg := service.config
	if err := service.pkr.UpdateConfigFromParams(cfg, params); err != nil {
		service.Logf("Failed to update params: %v", err)
	}
	postURL := buildPostURL(cfg)
	request := &messageRequest{
		Message:  message,
		Title:    cfg.MessageTitle,
		Priority: cfg.Priority,
		Topic:    cfg.Topic,
		Tags:     cfg.Tags,
	}
	response := &messageResponse{}
	// fmt.Printf("Sending %+v to %s\n", cfg, postURL)
	if err = service.client.Post(postURL, request, response); err != nil {
		fmt.Printf("Err %s\n", err.Error())
		errorRes := &errorResponse{}
		if service.client.ErrorResponse(err, errorRes) {
			return errorRes
		}
		return fmt.Errorf("failed to send notification to Ntfy: %s", err)
	}
	return
}

func buildPostURL(config *Config) string {
	scheme := "https"
	if config.DisableTLS {
		scheme = scheme[:4]
	}
	return fmt.Sprintf("%s://%s", scheme, config.Host)
}
