package gpt3

import (
	"context"
	"errors"
)

type GPT3client struct {
	client Client
}

func MakeGPT3Client(apikey string, options ...ClientOption) *GPT3client {
	return &GPT3client{
		client: NewClient(
			apikey,
			options...),
	}
}

func (c *GPT3client) DoStream(ctx context.Context, say []ChatCompletionMessage, fn func(cr CompletionResponseInterface)) error {
	if len(say) == 0 {
		return errors.New("您得说些什么。")
	}
	if c.client.DefaultEngine() == Gpt35TurboEngine ||
		c.client.DefaultEngine() == Gpt35Turbo0301Engine {
		request := ChatCompletionRequest{
			Model:     c.client.DefaultEngine(),
			Messages:  say,
			MaxTokens: c.client.Maxtokens(),
		}
		return c.client.ChatCompletionStream(ctx, request, fn)
	}
	// 组装 内容
	text := make([]string, len(say))
	for idx, v := range say {
		text[idx] = v.Content
	}
	// text := strings.Builder{}
	// for _, v := range say {
	// 	text.WriteString(v.Content)
	// }
	// tstr := text.String()
	// // max 4096 限制
	// if l := len(tstr); l > 4096 {
	// 	tstr = string([]rune(tstr[l-4096:])[2:])
	// }
	request := CompletionRequest{
		Prompt:    text,
		MaxTokens: c.client.Maxtokens(),
	}
	return c.client.CompletionStream(ctx, request, fn)
}

func (c *GPT3client) DoOnce(ctx context.Context, say []ChatCompletionMessage) (CompletionResponseInterface, error) {
	if len(say) == 0 {
		return nil, errors.New("您得说些什么。")
	}
	if c.client.DefaultEngine() == Gpt35TurboEngine ||
		c.client.DefaultEngine() == Gpt35Turbo0301Engine {
		request := ChatCompletionRequest{
			Model:     Gpt35TurboEngine,
			Messages:  say,
			Stop:      c.client.Stop(),
			MaxTokens: c.client.Maxtokens(),
		}
		return c.client.ChatCompletion(ctx, request)
	}
	request := CompletionRequest{
		Prompt:    []string{say[0].Content},
		MaxTokens: c.client.Maxtokens(),
	}
	return c.client.Completion(ctx, request)
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
