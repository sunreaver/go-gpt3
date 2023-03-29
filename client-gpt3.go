package gpt3

import (
	"context"
	"strings"

	"github.com/pkg/errors"
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
		request, err := c.makeChatCompletionRequest(tmp[0], say...)
		if err != nil {
			return err
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
		request, err := c.makeChatCompletionRequest(tmp[0], say...)
		if err != nil {
			return nil, err
		}
		return c.client.ChatCompletion(ctx, request)
	}
	return c.client.CompletionWithEngine(ctx, c.defaultEngine, c.makeCompletionRequest(tmp))
}

func (c *GPT3client) makeChatCompletionRequest(system ChatCompletionMessage, say ...ChatCompletionMessage) (ChatCompletionRequest, error) {
	// 组装 内容
	maxlen := 3896 - len(system.Content) // 不用 4096 是防止刚好贴边容易出问题；预留200空闲
	tmpCount := 0
CLIP:
	for i := len(say) - 1; i >= 0; i-- {
		tmpCount += len(say[i].Content)
		if tmpCount > maxlen {
			if i == len(say)-1 {
				// 第一个就超出
				return ChatCompletionRequest{}, errors.Errorf("输入内容过长; 最长%v, 当前%v", maxlen, tmpCount)
			}
			say = say[i+1:] // 丢弃此段内容
			break CLIP
		}
	}
	return ChatCompletionRequest{
		Model:     c.defaultEngine,
		Messages:  say,
		Stop:      c.stop,
		MaxTokens: &c.maxtokens,
	}, nil
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
	if maxlen, l := 3896-syslen, len(tstr); l > maxlen {
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
