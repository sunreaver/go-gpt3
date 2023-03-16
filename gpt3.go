package gpt3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type EngineType string

// Engine Types
const (
	TextAda001Engine     EngineType = "text-ada-001"
	TextBabbage001Engine EngineType = "text-babbage-001"
	TextCurie001Engine   EngineType = "text-curie-001"
	TextDavinci001Engine EngineType = "text-davinci-001"
	TextDavinci002Engine EngineType = "text-davinci-002"
	TextDavinci003Engine EngineType = "text-davinci-003"
	Gpt35TurboEngine     EngineType = "gpt-3.5-turbo"
	Gpt35Turbo0301Engine EngineType = "gpt-3.5-turbo-0301"
	ImageEngineCreate    EngineType = "images/generations"
	AdaEngine            EngineType = "ada"
	BabbageEngine        EngineType = "babbage"
	CurieEngine          EngineType = "curie"
	DavinciEngine        EngineType = "davinci"
	DefaultEngine        EngineType = Gpt35TurboEngine
)

type EmbeddingEngine string

const (
	TextSimilarityAda001      = "text-similarity-ada-001"
	TextSimilarityBabbage001  = "text-similarity-babbage-001"
	TextSimilarityCurie001    = "text-similarity-curie-001"
	TextSimilarityDavinci001  = "text-similarity-davinci-001"
	TextSearchAdaDoc001       = "text-search-ada-doc-001"
	TextSearchAdaQuery001     = "text-search-ada-query-001"
	TextSearchBabbageDoc001   = "text-search-babbage-doc-001"
	TextSearchBabbageQuery001 = "text-search-babbage-query-001"
	TextSearchCurieDoc001     = "text-search-curie-doc-001"
	TextSearchCurieQuery001   = "text-search-curie-query-001"
	TextSearchDavinciDoc001   = "text-search-davinci-doc-001"
	TextSearchDavinciQuery001 = "text-search-davinci-query-001"
	CodeSearchAdaCode001      = "code-search-ada-code-001"
	CodeSearchAdaText001      = "code-search-ada-text-001"
	CodeSearchBabbageCode001  = "code-search-babbage-code-001"
	CodeSearchBabbageText001  = "code-search-babbage-text-001"
	TextEmbeddingAda002       = "text-embedding-ada-002"
)

const (
	defaultBaseURL        = "https://api.openai.com/v1"
	defaultUserAgent      = "go-gpt3"
	defaultTimeoutSeconds = 30
)

func getEngineURL(engine string) string {
	return fmt.Sprintf("%s/engines/%s/completions", defaultBaseURL, engine)
}

type CompletionResponseInterface interface {
	CanContinue() bool
	Text() string
	Role() string
	Reset()
	TotalTokens() int
}

// A Client is an API client to communicate with the OpenAI gpt-3 APIs
type Client interface {
	// Engines lists the currently available engines, and provides basic information about each
	// option such as the owner and availability.
	// Engines(ctx context.Context) (*EnginesResponse, error)

	// Engine retrieves an engine instance, providing basic information about the engine such
	// as the owner and availability.
	// Engine(ctx context.Context, engine EngineType) (*EngineObject, error)

	// Completion creates a completion with the default engine. This is the main endpoint of the API
	// which auto-completes based on the given prompt.
	// Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)

	// CompletionStream creates a completion with the default engine and streams the results through
	// multiple calls to onData.
	// CompletionStream(ctx context.Context, request CompletionRequest, onData func(CompletionResponseInterface)) error

	// CompletionWithEngine is the same as Completion except allows overriding the default engine on the client
	CompletionWithEngine(ctx context.Context, engine EngineType, request CompletionRequest) (*CompletionResponse, error)

	// CompletionStreamWithEngine is the same as CompletionStream except allows overriding the default engine on the client
	CompletionStreamWithEngine(ctx context.Context, engine EngineType, request CompletionRequest, onData func(CompletionResponseInterface)) error

	ChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error)

	// CompletionStream creates a completion with the default engine and streams the results through
	// multiple calls to onData.
	ChatCompletionStream(ctx context.Context, request ChatCompletionRequest, onData func(CompletionResponseInterface)) error

	// Given a prompt and an instruction, the model will return an edited version of the prompt.
	// Edits(ctx context.Context, request EditsRequest) (*EditsResponse, error)

	// Search performs a semantic search over a list of documents with the default engine.
	// Search(ctx context.Context, request SearchRequest) (*SearchResponse, error)

	// SearchWithEngine performs a semantic search over a list of documents with the specified engine.
	// SearchWithEngine(ctx context.Context, engine EngineType, request SearchRequest) (*SearchResponse, error)

	// Returns an embedding using the provided request.
	// Embeddings(ctx context.Context, request EmbeddingsRequest) (*EmbeddingsResponse, error)

	CreateImage(ctx context.Context, request CreateImageReq) (*CreateImageResp, error)
}

type client struct {
	baseURL    string
	apiKey     string
	userAgent  string
	httpClient *http.Client
	idOrg      string

	gpt3 *GPT3client
}

// NewClient returns a new OpenAI GPT-3 API client. An apiKey is required to use the client
func NewClient(apiKey string, gpt3 *GPT3client, options ...ClientOption) Client {
	httpClient := &http.Client{
		Timeout: time.Duration(defaultTimeoutSeconds * time.Second),
	}

	c := &client{
		userAgent:  defaultUserAgent,
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: httpClient,
		idOrg:      "",

		gpt3: gpt3,
	}
	for _, o := range options {
		o(c)
	}
	return c
}

func (c *client) Engines(ctx context.Context) (*EnginesResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/engines", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(EnginesResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *client) Engine(ctx context.Context, engine EngineType) (*EngineObject, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/engines/%s", engine), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(EngineObject)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *client) CompletionWithEngine(ctx context.Context, engine EngineType, request CompletionRequest) (*CompletionResponse, error) {
	request.Stream = false
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/engines/%s/completions", engine), request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(CompletionResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *client) CompletionStreamWithEngine(
	ctx context.Context,
	engine EngineType,
	request CompletionRequest,
	onData func(CompletionResponseInterface),
) error {
	request.Stream = true
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/engines/%s/completions", engine), request)
	if err != nil {
		return err
	}
	return c.sendAndOnData(req, new(CompletionResponse), onData)
}

func (c *client) sendAndOnData(req *http.Request, output CompletionResponseInterface, onData func(CompletionResponseInterface)) error {
	resp, err := c.performRequest(req)
	if err != nil {
		return err
	}

	reader := newEventStreamReader(resp.Body, 1<<16)
	defer resp.Body.Close()

LOOP:
	for {
		event, err := reader.ReadEvent()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			return errors.Wrap(err, "ReadEvent")
		}

		// If we get an error, ignore it.
		var msg *Event
		if msg, err = processEvent(event); err != nil {
			return errors.Wrap(err, "ProcessEvent")
		}
		output.Reset()
		if bytes.Equal(msg.Data, doneSequence) {
			break LOOP
		}
		if err := json.Unmarshal(msg.Data, output); err != nil {
			return errors.Errorf("invalid json stream data: %v", err)
		}

		onData(output)
	}
	return nil
}

func (c *client) Edits(ctx context.Context, request EditsRequest) (*EditsResponse, error) {
	req, err := c.newRequest(ctx, "POST", "/edits", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(EditsResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

// func (c *client) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
// 	return c.SearchWithEngine(ctx, c.defaultEngine, request)
// }

func (c *client) SearchWithEngine(ctx context.Context, engine EngineType, request SearchRequest) (*SearchResponse, error) {
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/engines/%s/search", engine), request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}
	output := new(SearchResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

// Embeddings creates text embeddings for a supplied slice of inputs with a provided model.
//
// See: https://beta.openai.com/docs/api-reference/embeddings
func (c *client) Embeddings(ctx context.Context, request EmbeddingsRequest) (*EmbeddingsResponse, error) {
	req, err := c.newRequest(ctx, "POST", "/embeddings", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := EmbeddingsResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

func (c *client) performRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkForSuccess(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// returns an error if this response includes an error.
func checkForSuccess(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read from body: %w", err)
	}
	var result APIErrorResponse
	if err := json.Unmarshal(data, &result); err != nil {
		// if we can't decode the json error then create an unexpected error
		apiError := APIError{
			StatusCode: resp.StatusCode,
			Type:       "Unexpected",
			Message:    string(data),
		}
		return apiError
	}
	result.Error.StatusCode = resp.StatusCode
	return result.Error
}

func getResponseObject(rsp *http.Response, v interface{}) error {
	defer rsp.Body.Close()
	if err := json.NewDecoder(rsp.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}
	return nil
}

func jsonBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer(nil), nil
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed encoding json: %w", err)
	}
	return bytes.NewBuffer(raw), nil
}

func (c *client) newRequest(ctx context.Context, method, path string, payload interface{}) (*http.Request, error) {
	bodyReader, err := jsonBodyReader(payload)
	if err != nil {
		return nil, err
	}
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	if len(c.idOrg) > 0 {
		req.Header.Set("OpenAI-Organization", c.idOrg)
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	return req, nil
}
