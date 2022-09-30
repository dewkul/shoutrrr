package ntfy

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Config for use within the ntfy plugin
type Config struct {
	standard.EnumlessConfig
	Host         string   `desc:"Server hostname (and optionally port)" url:"host,port" default:"ntfy.sh"`
	Topic        string   `desc:"Target topic name (*Required)" url:"path1"`
	Username     string   `desc:"Username of a protected topic" url:"User" default:""`
	Password     string   `desc:"Password of a protected topic" url:"Password" default:""`
	Priority     uint8    `desc:"Message priority with 1=min, 3=default and 5=max" key:"priority" default:"3"`
	MessageTitle string   `desc:"Message title" key:"title" default:""`
	DisableTLS   bool     `desc:"Using http instead of https" key:"disabletls" default:"No"`
	Tags         []string `desc:"List of tags that may or not map to emojis separated by \",\" (comma)" key:"tags,tag" default:""`
}

// GetURL returns a URL representation of it's current field values
func (config *Config) GetURL() *url.URL {
	resolver := format.NewPropKeyResolver(config)
	return config.getURL(&resolver)
}

// SetURL updates a ServiceConfig from a URL representation of it's field values
func (config *Config) SetURL(url *url.URL) error {
	resolver := format.NewPropKeyResolver(config)
	return config.setURL(&resolver, url)
}

func (config *Config) getURL(resolver types.ConfigQueryResolver) *url.URL {
	var userInfo url.Userinfo
	if len(config.Username) > 0 && len(config.Password) > 0 {
		userInfo = *url.UserPassword(config.Username, config.Password)
	}
	return &url.URL{
		Host:       config.Host,
		Path:       config.Topic,
		User:       &userInfo,
		Scheme:     Scheme,
		ForceQuery: false,
		RawQuery:   format.BuildQuery(resolver),
	}
}

func (config *Config) setURL(resolver types.ConfigQueryResolver, url *url.URL) (err error) {
	config.Host = url.Hostname() + url.Port()
	if username := url.User.Username(); len(username) > 0 {
		config.Username = username
	}
	if password, isPasswordSet := url.User.Password(); isPasswordSet {
		config.Password = password
	}
	config.Topic = url.Path[1:]

	for key, vals := range url.Query() {
		if err = resolver.Set(key, vals[0]); err != nil {
			return
		}
	}
	return
}

const (
	// Scheme is the identifying part of this service's configuration URL
	Scheme = "ntfy"
)
