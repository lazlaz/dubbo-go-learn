package common

import (
	"bytes"
	"fmt"
	gxset "github.com/dubbogo/gost/container/set"
	"github.com/laz/dubbo-go/common/constant"
	"github.com/laz/dubbo-go/common/logger"
	"math"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

import (
	"github.com/jinzhu/copier"
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
	// CompareURLEqualFunc compare two url is equal
	compareURLEqualFunc CompareURLEqualFunc
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

// Clone will copy the url
func (c *URL) Clone() *URL {
	newURL := &URL{}
	if err := copier.Copy(newURL, c); err != nil {
		// this is impossible
		return newURL
	}
	newURL.params = url.Values{}
	c.RangeParams(func(key, value string) bool {
		newURL.SetParam(key, value)
		return true
	})

	return newURL
}

// RangeParams will iterate the params
func (c *URL) RangeParams(f func(key, value string) bool) {
	c.paramsLock.RLock()
	defer c.paramsLock.RUnlock()
	for k, v := range c.params {
		if !f(k, v[0]) {
			break
		}
	}
}

// Service gets service
func (c *URL) Service() string {
	service := c.GetParam(constant.INTERFACE_KEY, strings.TrimPrefix(c.Path, "/"))
	if service != "" {
		return service
	} else if c.SubURL != nil {
		service = c.SubURL.GetParam(constant.INTERFACE_KEY, strings.TrimPrefix(c.Path, "/"))
		if service != "" { // if url.path is "" then return suburl's path, special for registry url
			return service
		}
	}
	return ""
}

// Key gets key
func (c *URL) Key() string {
	buildString := fmt.Sprintf("%s://%s:%s@%s:%s/?interface=%s&group=%s&version=%s",
		c.Protocol, c.Username, c.Password, c.Ip, c.Port, c.Service(), c.GetParam(constant.GROUP_KEY, ""), c.GetParam(constant.VERSION_KEY, ""))
	return buildString
}

// Copy url based on the reserved parameter's keys.
func (c *URL) CloneWithParams(reserveParams []string) *URL {
	params := url.Values{}
	for _, reserveParam := range reserveParams {
		v := c.GetParam(reserveParam, "")
		if len(v) != 0 {
			params.Set(reserveParam, v)
		}
	}

	return NewURLWithOptions(
		WithProtocol(c.Protocol),
		WithUsername(c.Username),
		WithPassword(c.Password),
		WithIp(c.Ip),
		WithPort(c.Port),
		WithPath(c.Path),
		WithMethods(c.Methods),
		WithParams(params),
	)
}

// GetParams gets values
func (c *URL) GetParams() url.Values {
	return c.params
}
func (c *URL) CloneExceptParams(excludeParams *gxset.HashSet) *URL {
	newURL := &URL{}
	if err := copier.Copy(newURL, c); err != nil {
		// this is impossible
		return newURL
	}
	newURL.params = url.Values{}
	c.RangeParams(func(key, value string) bool {
		if !excludeParams.Contains(key) {
			newURL.SetParam(key, value)
		}
		return true
	})
	return newURL
}
func (c *URL) String() string {
	c.paramsLock.Lock()
	defer c.paramsLock.Unlock()
	var buf strings.Builder
	if len(c.Username) == 0 && len(c.Password) == 0 {
		buf.WriteString(fmt.Sprintf("%s://%s:%s%s?", c.Protocol, c.Ip, c.Port, c.Path))
	} else {
		buf.WriteString(fmt.Sprintf("%s://%s:%s@%s:%s%s?", c.Protocol, c.Username, c.Password, c.Ip, c.Port, c.Path))
	}
	buf.WriteString(c.params.Encode())
	return buf.String()
}

// ServiceKey gets a unique key of a service.
func (c *URL) ServiceKey() string {
	return ServiceKey(c.GetParam(constant.INTERFACE_KEY, strings.TrimPrefix(c.Path, "/")),
		c.GetParam(constant.GROUP_KEY, ""), c.GetParam(constant.VERSION_KEY, ""))
}
func MergeUrl(serviceUrl *URL, referenceUrl *URL) *URL {
	// After Clone, it is a new url that there is no thread safe issue.
	mergedUrl := serviceUrl.Clone()
	params := mergedUrl.GetParams()
	// iterator the referenceUrl if serviceUrl not have the key ,merge in
	// referenceUrl usually will not changed. so change RangeParams to GetParams to avoid the string value copy.
	for key, value := range referenceUrl.GetParams() {
		if v := mergedUrl.GetParam(key, ""); len(v) == 0 {
			if len(value) > 0 {
				params[key] = value
			}
		}
	}

	// loadBalance,cluster,retries strategy config
	methodConfigMergeFcn := mergeNormalParam(params, referenceUrl, []string{constant.LOADBALANCE_KEY, constant.CLUSTER_KEY, constant.RETRIES_KEY, constant.TIMEOUT_KEY})

	// remote timestamp
	if v := serviceUrl.GetParam(constant.TIMESTAMP_KEY, ""); len(v) > 0 {
		params[constant.REMOTE_TIMESTAMP_KEY] = []string{v}
		params[constant.TIMESTAMP_KEY] = []string{referenceUrl.GetParam(constant.TIMESTAMP_KEY, "")}
	}

	// finally execute methodConfigMergeFcn
	for _, method := range referenceUrl.Methods {
		for _, fcn := range methodConfigMergeFcn {
			fcn("methods." + method)
		}
	}
	// In this way, we will raise some performance.
	mergedUrl.ReplaceParams(params)
	return mergedUrl
}
func mergeNormalParam(params url.Values, referenceUrl *URL, paramKeys []string) []func(method string) {
	methodConfigMergeFcn := make([]func(method string), 0, len(paramKeys))
	for _, paramKey := range paramKeys {
		if v := referenceUrl.GetParam(paramKey, ""); len(v) > 0 {
			params[paramKey] = []string{v}
		}
		methodConfigMergeFcn = append(methodConfigMergeFcn, func(method string) {
			if v := referenceUrl.GetParam(method+"."+paramKey, ""); len(v) > 0 {
				params[method+"."+paramKey] = []string{v}
			}
		})
	}
	return methodConfigMergeFcn
}

// GetParamInt gets int64 value by @key
func (c *URL) GetParamInt(key string, d int64) int64 {
	r, err := strconv.ParseInt(c.GetParam(key, ""), 10, 64)
	if err != nil {
		return d
	}
	return r
}

// GetMethodParamInt64 gets int64 method param
func (c *URL) GetMethodParamInt64(method string, key string, d int64) int64 {
	r := c.GetMethodParamInt(method, key, math.MinInt64)
	if r == math.MinInt64 {
		return c.GetParamInt(key, d)
	}
	return r
}

type CompareURLEqualFunc func(l *URL, r *URL, excludeParam ...string) bool

func GetCompareURLEqualFunc() CompareURLEqualFunc {
	return compareURLEqualFunc
}

// GetMethodParamInt gets int method param
func (c *URL) GetMethodParamInt(method string, key string, d int64) int64 {
	r, err := strconv.ParseInt(c.GetParam("methods."+method+"."+key, ""), 10, 64)
	if err != nil {
		return d
	}
	return r
}

// GetMethodParamBool judge whether @method param exists or not
func (c *URL) GetMethodParamBool(method string, key string, d bool) bool {
	r := c.GetParamBool("methods."+method+"."+key, d)
	return r
}

// GetMethodParam gets method param
func (c *URL) GetMethodParam(method string, key string, d string) string {
	r := c.GetParam("methods."+method+"."+key, "")
	if r == "" {
		r = d
	}
	return r
}

// ReplaceParams will replace the URL.params
// usually it should only be invoked when you want to modify an url, such as MergeURL
func (c *URL) ReplaceParams(param url.Values) {
	c.paramsLock.Lock()
	defer c.paramsLock.Unlock()
	c.params = param
}
