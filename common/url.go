package common

import (
	"bytes"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/logger"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

import (
	perrors "github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// role constant
const (
	// CONSUMER is consumer role
	CONSUMER = iota
	// CONFIGURATOR is configurator role
	CONFIGURATOR
	// ROUTER is router role
	ROUTER
	// PROVIDER is provider role
	PROVIDER
	PROTOCOL = "protocol"
)

var (
	// DubboNodes Dubbo service node
	DubboNodes = [...]string{"consumers", "configurators", "routers", "providers"}
	// DubboRole Dubbo service role
	DubboRole = [...]string{"consumer", "", "routers", "provider"}
)

// nolint
type RoleType int

func (t RoleType) String() string {
	return DubboNodes[t]
}

// WithParamsValue sets params field for url
func WithParamsValue(key, val string) Option {
	return func(url *URL) {
		url.SetParam(key, val)
	}
}

// Role returns role by @RoleType
func (t RoleType) Role() string {
	return DubboRole[t]
}

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

func ServiceKey(intf string, group string, version string) string {
	if intf == "" {
		return ""
	}
	buf := &bytes.Buffer{}
	if group != "" {
		buf.WriteString(group)
		buf.WriteString("/")
	}

	buf.WriteString(intf)

	if version != "" && version != "0.0.0" {
		buf.WriteString(":")
		buf.WriteString(version)
	}

	return buf.String()
}

// NewURLWithOptions will create a new url with options
func NewURLWithOptions(opts ...Option) *URL {
	newURL := &URL{}
	for _, opt := range opts {
		opt(newURL)
	}
	newURL.Location = newURL.Ip + ":" + newURL.Port
	return newURL
}

// WithPath sets path for url
func WithPath(path string) Option {
	return func(url *URL) {
		url.Path = "/" + strings.TrimPrefix(path, "/")
	}
}

// WithIp sets ip for url
func WithIp(ip string) Option {
	return func(url *URL) {
		url.Ip = ip
	}
}

// WithPort sets port for url
func WithPort(port string) Option {
	return func(url *URL) {
		url.Port = port
	}
}

// WithToken sets token for url
func WithToken(token string) Option {
	return func(url *URL) {
		if len(token) > 0 {
			value := token
			if strings.ToLower(token) == "true" || strings.ToLower(token) == "default" {
				u, err := uuid.NewV4()
				if err != nil {
					logger.Errorf("could not generator UUID: %v", err)
					return
				}
				value = u.String()
			}
			url.SetParam(constant.TOKEN_KEY, value)
		}
	}
}

// WithMethods sets methods for url
func WithMethods(methods []string) Option {
	return func(url *URL) {
		url.Methods = methods
	}
}

// AddParam will add the key-value pair
func (c *URL) AddParam(key string, value string) {
	c.paramsLock.Lock()
	defer c.paramsLock.Unlock()
	if c.params == nil {
		c.params = url.Values{}
	}
	c.params.Add(key, value)
}

// GetParamBool judge whether @key exists or not
func (c *URL) GetParamBool(key string, d bool) bool {
	r, err := strconv.ParseBool(c.GetParam(key, ""))
	if err != nil {
		return d
	}
	return r
} // GetParam gets value by key
func (c *URL) GetParam(s string, d string) string {
	c.paramsLock.RLock()
	defer c.paramsLock.RUnlock()

	var r string
	if len(c.params) > 0 {
		r = c.params.Get(s)
	}
	if len(r) == 0 {
		r = d
	}

	return r
}
