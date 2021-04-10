package presentme

import (
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	dc "github.com/stanistan/present-me/internal/cache"
)

type Config struct {
	DiskCache dc.CacheOpts `embed:"" prefix:"disk-cache-"`
	Github    GHOpts       `embed:"" prefix:"gh-"`
}

func (c *Config) Configure() {
	log.Infof("config %+v", c)
	configureCache(c.DiskCache)
}

func LoadConfig(path string) (Config, error) {
	var c Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return c, errors.Wrapf(err, "could not read config at path %s", path)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, errors.Wrapf(err, "error parsing the config at path %s", path)
	}

	return c, nil
}

func MustConfig(path string) Config {
	c, err := LoadConfig(path)
	if err != nil {
		panic(err)
	}

	c.Configure()
	return c
}

func configureCache(opts dc.CacheOpts) {
	cache = dc.NewCache(opts)
}

var (
	cache    *dc.Cache = dc.NewCache(dc.CacheOpts{Enabled: false})
	cacheTTL           = 10 * time.Minute
)
