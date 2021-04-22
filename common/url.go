package common

import (
	"net"
	"net/url"
	"strings"
	"sync"
)

import (
	perrors "github.com/pkg/errors"
)

type baseUrl struct {
	Protocol string
	Location string // ip+port
	Ip       string
	Port     string

	PrimitiveURL string
}
type URL struct {
	params url.Values
	baseUrl
	//url.Values is not safe map, add to avoid concurrent map read and map write error
	paramsLock sync.RWMutex
	Path       string // like  /com.ikurento.dubbo.UserProvider
	Username   string
	Password   string
	Methods    []string
	// special for registry
	SubURL *URL
}
type Option func(*URL)

// NewURL will create a new url
// the urlString should not be empty
func NewURL(urlString string, opts ...Option) (*URL, error) {
	s := URL{baseUrl: baseUrl{}}
	if urlString == "" {
		return &s, nil
	}

	rawUrlString, err := url.QueryUnescape(urlString)
	if err != nil {
		return &s, perrors.Errorf("url.QueryUnescape(%s),  error{%v}", urlString, err)
	}

	// rawUrlString = "//" + rawUrlString
	if !strings.Contains(rawUrlString, "//") {
		t := URL{baseUrl: baseUrl{}}
		for _, opt := range opts {
			opt(&t)
		}
		rawUrlString = t.Protocol + "://" + rawUrlString
	}

	serviceUrl, urlParseErr := url.Parse(rawUrlString)
	if urlParseErr != nil {
		return &s, perrors.Errorf("url.Parse(url string{%s}),  error{%v}", rawUrlString, err)
	}

	s.params, err = url.ParseQuery(serviceUrl.RawQuery)
	if err != nil {
		return &s, perrors.Errorf("url.ParseQuery(raw url string{%s}),  error{%v}", serviceUrl.RawQuery, err)
	}

	s.PrimitiveURL = urlString
	s.Protocol = serviceUrl.Scheme
	s.Username = serviceUrl.User.Username()
	s.Password, _ = serviceUrl.User.Password()
	s.Location = serviceUrl.Host
	s.Path = serviceUrl.Path
	if strings.Contains(s.Location, ":") {
		s.Ip, s.Port, err = net.SplitHostPort(s.Location)
		if err != nil {
			return &s, perrors.Errorf("net.SplitHostPort(url.Host{%s}), error{%v}", s.Location, err)
		}
	}
	for _, opt := range opts {
		opt(&s)
	}
	return &s, nil
}

// WithProtocol sets protocol for url
func WithProtocol(proto string) Option {
	return func(url *URL) {
		url.Protocol = proto
	}
}

// WithUsername sets username for url
func WithUsername(username string) Option {
	return func(url *URL) {
		url.Username = username
	}
}

// WithPassword sets password for url
func WithPassword(pwd string) Option {
	return func(url *URL) {
		url.Password = pwd
	}
}

// WithLocation sets location for url
func WithLocation(location string) Option {
	return func(url *URL) {
		url.Location = location
	}
}

// SetParams will put all key-value pair into url.
// 1. if there already has same key, the value will be override
// 2. it's not thread safe
// 3. think twice when you want to invoke this method
func (c *URL) SetParams(m url.Values) {
	for k := range m {
		c.SetParam(k, m.Get(k))
	}
}

// SetParam will put the key-value pair into url
// usually it should only be invoked when you want to initialized an url
func (c *URL) SetParam(key string, value string) {
	c.paramsLock.Lock()
	defer c.paramsLock.Unlock()
	if c.params == nil {
		c.params = url.Values{}
	}
	c.params.Set(key, value)
}

// WithParams sets params for url
func WithParams(params url.Values) Option {
	return func(url *URL) {
		url.params = params
	}
}
