package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

// Config stores global configuration data of the application.
type Config struct {
	Debug                   bool
	DoNotVerifyCrawleraCert bool `toml:"dont_verify_crawlera_cert"`
	BindPort                int
	CrawleraPort            int `toml:"crawlera_port"`
	BindIP                  string
	APIKey                  string
	CrawleraHost            string `toml:"crawlera_host"`
	XHeaders                map[string]string
}

// Bind returns a string for the http.ListenAndServe based on config
// information.
func (c *Config) Bind() string {
	return net.JoinHostPort(c.BindIP, strconv.Itoa(c.BindPort))
}

// CrawleraURL builds and returns URL to crawlera. Basically, this is required
// for http.ProxyURL to have embedded credentials etc.
func (c *Config) CrawleraURL() string {
	return fmt.Sprintf("http://%s:@%s",
		c.APIKey,
		net.JoinHostPort(c.CrawleraHost, strconv.Itoa(c.CrawleraPort)))
}

// MaybeSetDebug enabled debug mode of crawlera-headless-proxy (verbosity
// mostly). If given value is not defined (false) then changes nothing.
func (c *Config) MaybeSetDebug(value bool) {
	c.Debug = c.Debug || value
}

// MaybeDoNotVerifyCrawleraCert defines is it necessary to verify Crawlera
// TLS certificate. If given value is not defined (false) then changes nothing.
func (c *Config) MaybeDoNotVerifyCrawleraCert(value bool) {
	c.DoNotVerifyCrawleraCert = c.DoNotVerifyCrawleraCert || value
}

// MaybeSetBindIP sets an IP crawlera-headless-proxy should listen on.
// If given value is not defined (0) then changes nothing.
//
// If you want to have a global access (which is not recommended) please
// set it to 0.0.0.0.
func (c *Config) MaybeSetBindIP(value net.IP) {
	if value != nil {
		c.BindIP = value.String()
	}
}

// MaybeSetBindPort sets a port crawlera-headless-proxy should listen on.
// If given value is not defined (0) then changes nothing.
func (c *Config) MaybeSetBindPort(value int) {
	if value > 0 {
		c.BindPort = value
	}
}

// MaybeSetAPIKey sets an API key of Crawlera. If given value is not
// defined ("") then changes nothing.
func (c *Config) MaybeSetAPIKey(value string) {
	if value != "" {
		c.APIKey = value
	}
}

// MaybeSetCrawleraHost sets a host of Crawlera (usually it is
// 'proxy.crawlera.com'). If given value is not defined ("") then changes
// nothing.
func (c *Config) MaybeSetCrawleraHost(value string) {
	if value != "" {
		c.CrawleraHost = value
	}
}

// MaybeSetCrawleraPort a port Crawlera is listening to (usually it is 8010).
// If given value is not defined (0) then changes nothing.
func (c *Config) MaybeSetCrawleraPort(value int) {
	if value > 0 {
		c.CrawleraPort = value
	}
}

// SetXHeader sets a header value of Crawlera X-Header. It is actually
// allowed to pass values in both ways: with full name (x-crawlera-profile)
// for example, and in the short form: just 'profile'. This effectively the
// same.
func (c *Config) SetXHeader(key, value string) {
	key = strings.ToLower(key)
	key = strings.TrimPrefix(key, "x-crawlera-")
	key = strings.Title(key)

	c.XHeaders[fmt.Sprintf("X-Crawlera-%s", key)] = value
}

// Parse processes incoming file handler (usually, an instance of *os.File)
// and returns an instance of Config with fields set.
//
// Basically, new Config instance gets its fields in this order:
//   1. Defaults
//   2. Values from the config file.
func Parse(file io.Reader) (*Config, error) {
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Annotate(err, "Cannot read config file")
	}

	conf := NewConfig()
	if _, err := toml.Decode(string(buf), conf); err != nil {
		return nil, errors.Annotate(err, "Cannot parse config file")
	}

	xheaders := conf.XHeaders
	conf.XHeaders = map[string]string{}
	for k, v := range xheaders {
		conf.SetXHeader(k, v)
	}

	return conf, nil
}

// NewConfig returns new instance of configuration data structure with
// fields set to sensible defaults.
func NewConfig() *Config {
	return &Config{
		BindIP:       "127.0.0.1",
		BindPort:     3128,
		CrawleraHost: "proxy.crawlera.com",
		CrawleraPort: 8010,
		XHeaders:     map[string]string{},
	}
}