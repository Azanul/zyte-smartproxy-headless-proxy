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

type Config struct {
	Debug                   bool
	DoNotVerifyCrawleraCert bool `toml:"dont_verify_crawlera_cert"`
	BindIP                  string
	BindPort                int
	APIKey                  string
	CrawleraHost            string `toml:"crawlera_host"`
	CrawleraPort            int    `toml:"crawlera_port"`
	XHeaders                map[string]string
}

func (c *Config) Bind() string {
	return net.JoinHostPort(c.BindIP, strconv.Itoa(c.BindPort))
}

func (c *Config) CrawleraURL() string {
	return fmt.Sprintf("http://%s:@%s",
		c.APIKey,
		net.JoinHostPort(c.CrawleraHost, strconv.Itoa(c.CrawleraPort)))
}

func (c *Config) MaybeSetDebug(value bool) {
	c.Debug = c.Debug || value
}

func (c *Config) MaybeDoNotVerifyCrawleraCert(value bool) {
	c.DoNotVerifyCrawleraCert = c.DoNotVerifyCrawleraCert || value
}

func (c *Config) MaybeSetBindIP(value net.IP) {
	if value != nil {
		c.BindIP = value.String()
	}
}

func (c *Config) MaybeSetBindPort(value int) {
	if value > 0 {
		c.BindPort = value
	}
}

func (c *Config) MaybeSetAPIKey(value string) {
	if value != "" {
		c.APIKey = value
	}
}

func (c *Config) MaybeSetCrawleraHost(value string) {
	if value != "" {
		c.CrawleraHost = value
	}
}

func (c *Config) MaybeSetCrawleraPort(value int) {
	if value > 0 {
		c.CrawleraPort = value
	}
}

func (c *Config) SetXHeader(key, value string) {
	key = strings.ToLower(key)
	if strings.HasPrefix(key, "x-crawlera-") {
		key = key[len("x-crawlera-"):]
	}
	c.XHeaders[fmt.Sprintf("X-Crawlera-%s", strings.Title(key))] = value
}

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

func NewConfig() *Config {
	return &Config{
		BindIP:       "127.0.0.1",
		BindPort:     3128,
		CrawleraHost: "proxy.crawlera.com",
		CrawleraPort: 8010,
		XHeaders:     map[string]string{},
	}
}
