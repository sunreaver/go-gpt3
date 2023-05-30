package gpt3

import (
	"net/http"
	"time"
)

// ClientOption are options that can be passed when creating a new client
type ClientOption func(*client) error

// WithOrg is a client option that allows you to override the organization ID
func WithOrg(id string) ClientOption {
	return func(c *client) error {
		c.idOrg = id
		return nil
	}
}

// WithDefaultEngine is a client option that allows you to override the default engine of the client
func WithDefaultEngine(engine EngineType) ClientOption {
	return func(c *client) error {
		c.gpt3.defaultEngine = engine
		return nil
	}
}

// WithUserAgent is a client option that allows you to override the default user agent of the client
func WithUserAgent(userAgent string) ClientOption {
	return func(c *client) error {
		c.userAgent = userAgent
		return nil
	}
}

// WithBaseURL is a client option that allows you to override the default base url of the client.
// The default base url is "https://api.openai.com/v1"
func WithBaseURL(baseURL string) ClientOption {
	return func(c *client) error {
		if baseURL == "" {
			return nil
		}
		c.baseURL = baseURL
		return nil
	}
}

// WithHTTPClient allows you to override the internal http.Client used
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout is a client option that allows you to override the default timeout duration of requests
// for the client. The default is 30 seconds. If you are overriding the http client as well, just include
// the timeout there.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}

// Stop
func WithStop(stop []string) ClientOption {
	return func(c *client) error {
		c.gpt3.stop = stop
		return nil
	}
}

// maxtokens
func WithMaxtokens(maxtokens int) ClientOption {
	return func(c *client) error {
		c.gpt3.maxtokens = maxtokens
		return nil
	}
}

// 注入系统提示
func WithSystemPrompt(prompt string) ClientOption {
	return func(c *client) error {
		c.gpt3.systemprompt = prompt
		return nil
	}
}

// 注入请求query参数，字符串格式: key1=v1&key2=v2
func WithQuery(query string) ClientOption {
	return func(c *client) error {
		c.gpt3.query = query
		return nil
	}
}

// 注入apikey；将会放在header中: api-key: YOUR_API_KEY
// 用户Azure的权限。
func WithApiKey(apikey string) ClientOption {
	return func(c *client) error {
		c.gpt3.apikey = apikey
		return nil
	}
}

// 注入Authtoken；将会放在header中: Authorization: Bearer YOUR_AUTH_TOKEN.
// 用户Openai的权限。
func WithAuthtoken(authtoken string) ClientOption {
	return func(c *client) error {
		c.gpt3.authtoken = authtoken
		return nil
	}
}

// WithMaxRetry 注入重试次数。
func WithMaxRetry(try int) ClientOption {
	if try < 1 {
		try = DefaultRetry
	}

	return func(c *client) error {
		c.gpt3.maxretry = try
		return nil
	}
}
