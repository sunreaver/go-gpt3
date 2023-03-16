package gpt3

import (
	"context"
	"errors"
	"strings"
)

type GPT3client struct {
	client Client

	stop          []string
	maxtokens     int
	systemprompt  string
	defaultEngine EngineType
}

func MakeGPT3Client(apikey string, options ...ClientOption) *GPT3client {
	c := &GPT3client{
		defaultEngine: DefaultEngine,
		maxtokens:     256,
		stop:          nil,
	}
	c.client = NewClient(
		apikey,
		c,
		options...)

	return c
}

func (c *GPT3client) DoStream(ctx context.Context, say []ChatCompletionMessage, fn func(cr CompletionResponseInterface)) error {
	if len(say) == 0 {
		return errors.New("您得说些什么。")
	}
	tmp := append([]ChatCompletionMessage{
		{
			Role:    "system",
			Content: c.systemprompt,
		},
	}, say...)
	if c.defaultEngine == Gpt35TurboEngine ||
		c.defaultEngine == Gpt35Turbo0301Engine {
		request := ChatCompletionRequest{
			Model:     c.defaultEngine,
			Messages:  tmp,
			MaxTokens: &c.maxtokens,
		}
		return c.client.ChatCompletionStream(ctx, request, fn)
	}
	return c.client.CompletionStreamWithEngine(ctx, c.defaultEngine, c.makeCompletionRequest(tmp), fn)
}

func (c *GPT3client) DoOnce(ctx context.Context, say []ChatCompletionMessage) (CompletionResponseInterface, error) {
	if len(say) == 0 {
		return nil, errors.New("您得说些什么。")
	}
	tmp := append([]ChatCompletionMessage{
		{
			Role:    "system",
			Content: c.systemprompt,
		},
	}, say...)
	if c.defaultEngine == Gpt35TurboEngine ||
		c.defaultEngine == Gpt35Turbo0301Engine {
		request := ChatCompletionRequest{
			Model:     Gpt35TurboEngine,
			Messages:  tmp,
			Stop:      c.stop,
			MaxTokens: &c.maxtokens,
		}
		return c.client.ChatCompletion(ctx, request)
	}
	return c.client.CompletionWithEngine(ctx, c.defaultEngine, c.makeCompletionRequest(tmp))
}

func (c *GPT3client) makeCompletionRequest(say []ChatCompletionMessage) CompletionRequest {
	// 组装 内容
	text := strings.Builder{}
	system := ""
	for _, v := range say {
		if v.Role == "system" {
			system = v.Content
		} else {
			text.WriteString(v.Content)
		}
	}
	syslen := len(system)
	tstr := text.String()
	// max 4096 限制
	if maxlen, l := 4096-syslen, len(tstr); l > maxlen {
		tstr = string([]rune(tstr[l-maxlen:])[2:])
	}
	return CompletionRequest{
		Prompt:    []string{system + tstr},
		MaxTokens: &c.maxtokens,
	}
}

func (c *GPT3client) CreateImage(ctx context.Context, say CreateImageReq) (*CreateImageResp, error) {
	if len(say.Prompt) == 0 {
		return nil, errors.New("您得说些什么。")
	} else if say.N > 10 {
		return nil, errors.New("不要太贪心，先试试取一幅图~")
	} else if say.N <= 0 {
		say.N = 1
	}
	return c.client.CreateImage(ctx, say)
}
