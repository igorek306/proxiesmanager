package proxiesmanager

import (
	"bufio"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ProxiesManager struct {
	proxies                 []*url.URL
	useIndex                int
	defaultScheme           string
	defaultProxyTryAttempts int
}

type TargetSite struct {
	url     string
	request *http.Request
}

// Returns pointer to ProxiesManager struct.
func NewProxiesManager() *ProxiesManager {
	return &ProxiesManager{
		defaultScheme:           "http",
		defaultProxyTryAttempts: 2,
	}
}

// Sets default scheme which will be used while loading proxies, if scheme is not provided
func (m *ProxiesManager) DefaultScheme(scheme string) string {
	if scheme != "" {
		m.defaultScheme = scheme
	}
	return m.defaultScheme
}

// Sets count of attempts for all loaded proxies for one request, must be > 0.
func (m *ProxiesManager) DefaultTryAttempts(count int) int {
	if count < 1 {
		m.defaultProxyTryAttempts = count
	}
	return m.defaultProxyTryAttempts
}

// Loads proxies from file.
// Each one line of file is parsed by url.Parse function.
func (m *ProxiesManager) LoadFromFile(path string) (int, error) {
	proxyFile, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer proxyFile.Close()
	loaded := m.loadFromBuff(proxyFile)
	return loaded, nil
}

// Loads proxies from web.
// Each one line of response body is parsed by url.Parse function.
// Note: if you pass http client and set use proxies to true, http client will remain with the transport set to the last proxy
func (m *ProxiesManager) LoadFromWeb(ts TargetSite, httpClient *http.Client, useProxies bool) (int, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	if ts.request == nil {
		if ts.url == "" {
			return 0, errors.New("neither url or request was provided")
		}
		var err error
		ts.request, err = http.NewRequest("GET", ts.url, nil)
		if err != nil {
			return 0, nil
		}
	}
	var res *http.Response
	var err error
	if !useProxies {
		res, err = httpClient.Do(ts.request)
		if err != nil {
			return 0, err
		}
	} else {
		res, err = m.proxiedRequest(httpClient, ts.request)
		if err != nil {
			return 0, err
		}
	}
	defer res.Body.Close()
	loaded := m.loadFromBuff(res.Body)
	return loaded, nil
}

func (m *ProxiesManager) CallRequest(httpClient *http.Client, request *http.Request) (*http.Response, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return m.proxiedRequest(httpClient, request)
}

func (m *ProxiesManager) proxiedRequest(httpClient *http.Client, request *http.Request) (*http.Response, error) {
	var res *http.Response
	if m.Count() < 1 {
		return res, errors.New("no proxies loaded")
	}
	var t *http.Transport
	var err error
	for i := 0; i < m.defaultProxyTryAttempts*m.Count(); i++ {
		t = &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return m.Proxy(), nil
			},
		}
		httpClient.Transport = t
		res, err = httpClient.Do(request)
		if err == nil {
			break
		}
		m.Next()
	}
	return res, err

}

func (m *ProxiesManager) loadFromBuff(buff io.Reader) int {
	scanner := bufio.NewScanner(buff)
	var proxyLine string
	var loaded int
	for scanner.Scan() {
		proxyLine = scanner.Text()
		if !strings.Contains(proxyLine, "://") {
			proxyLine = m.defaultScheme + "://" + proxyLine
		}
		host, err := url.Parse(proxyLine)

		if err != nil || host.Host == "" {
			continue
		}
		if host.Scheme == "" {
			host.Scheme = "http"
		}
		m.proxies = append(m.proxies, host)
		loaded++
	}
	return loaded
}

// Loads URL pointers into the proxy array.
func (m *ProxiesManager) InsertProxies(proxies *[]*url.URL) {
	m.proxies = append(m.proxies, *proxies...)
}

// Calls m.Next() and return m.Proxy().
func (m *ProxiesManager) NextProxy() *url.URL {
	m.Next()
	return m.Proxy()
}

// Returns pointer to URL of the current proxy.
func (m *ProxiesManager) Proxy() *url.URL {
	if m.Count() < 1 {
		return nil
	}
	return m.proxies[m.useIndex]
}

// Increases index of the current proxy, which means m.Proxy() will return a different proxy.
func (m *ProxiesManager) Next() {
	m.useIndex++
	if m.useIndex >= len(m.proxies) {
		m.useIndex = 0
	}
}

// Returns the number of loaded proxies.
func (m *ProxiesManager) Count() int {
	return len(m.proxies)
}

// Sets the current proxy index to random, between 0 and count of loaded proxies.
// Should be called after proxies are loaded but before they are used to prevent using the same proxy all the time.
func (m *ProxiesManager) SelectRandom() {
	m.useIndex = rand.Intn(m.Count())
}

// Returns pointer to the array of proxies.
// Should be used as read-only.
func (m *ProxiesManager) Proxies() *[]*url.URL {
	return &m.proxies
}

// Writes each one loaded proxy, each in a new line,
// and returns sum of all bytes written and error if it occurred while writing.
func (m *ProxiesManager) PrintAll(buffer io.Writer) (int, error) {
	var bytesWritten int
	for _, proxy := range m.proxies {
		n, err := buffer.Write([]byte(proxy.String() + "\n"))
		bytesWritten += n
		if err != nil {
			return bytesWritten, nil
		}
	}
	return bytesWritten, nil
}
