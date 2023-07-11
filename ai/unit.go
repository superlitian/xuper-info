package ai

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"xuper-info/common"
)

type GPTRequest struct {
	// A list of string prompts to use.
	Model string `json:"model"`
	// TODO there are other prompt types here for using token integers that we could add support for.
	Prompt string `json:"prompt"`
	// How many tokens to complete up to. Max of 512
	MaxTokens *int `json:"max_tokens,omitempty"`
	// Sampling temperature to use
	Temperature *float32 `json:"temperature,omitempty"`
	// Alternative to temperature for nucleus sampling
	TopP *float32 `json:"top_p,omitempty"`
	// How many choice to create for each prompt
	N *int `json:"n"`
	// Include the probabilities of most likely tokens
	LogProbs *int `json:"logprobs"`
	// Echo back the prompt in addition to the completion
	Echo bool `json:"echo"`
	// Up to 4 sequences where the API will stop generating tokens. Response will not contain the stop sequence.
	Stop []string `json:"stop,omitempty"`
	// PresencePenalty number between 0 and 1 that penalizes tokens that have already appeared in the text so far.
	PresencePenalty float32 `json:"presence_penalty"`
	// FrequencyPenalty number between 0 and 1 that penalizes tokens on existing frequency in the text so far.
	FrequencyPenalty float32 `json:"frequency_penalty"`

	// Whether to stream back results or not. Don't set this value in the request yourself
	// as it will be overriden depending on if you use CompletionStream or Completion methods.
	Stream bool `json:"stream,omitempty"`
}

type GPTResponse struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	Created int                 `json:"created"`
	Model   string              `json:"model"`
	Choices []GPTResponseChoice `json:"choices"`
	Usage   GPTResponseUsage    `json:"usage"`
}

type GPTResponseChoice struct {
	Text         string           `json:"text"`
	Index        int              `json:"index"`
	LogProbs     GPTLogprobResult `json:"logprobs"`
	FinishReason string           `json:"finish_reason"`
}

type GPTResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GPTLogprobResult struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float32            `json:"token_logprobs"`
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

type GPTAPIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
}

type GPTAPIErrorResponse struct {
	Error GPTAPIError `json:"error"`
}

func Unit(text string) {
	baseURL := "https://api.openai.com/v1/completions"
	gptRequest := GPTRequest{
		Model:       "text-davinci-001",
		Prompt:      text,
		MaxTokens:   intPtr(2000),
		Temperature: float32Ptr(0),
		Stream:      false,
	}
	bodyReader, err := jsonBodyReader(gptRequest)
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequestWithContext(context.Background(), "POST", baseURL, bodyReader)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hclient := &http.Client{Transport: tr, Timeout: time.Duration(60 * time.Second)}
	resp, err := hclient.Do(req)
	if err != nil {
		common.Res.Msg = "服务错误，请稍后重试"
		common.SendMsg()
		return
	}
	if err := checkForSuccess(resp); err != nil {
		common.Res.Msg = "服务错误，请稍后重试"
		common.SendMsg()
		return
	}
	output := new(GPTResponse)
	if err := getResponseObject(resp, output); err != nil {
		common.Res.Msg = "服务错误，请稍后重试"
		common.SendMsg()
		return
	}
	common.Res.Msg = strings.Replace(output.Choices[0].Text, "\n", "", 2)
	common.SendMsg()
}

func (e GPTAPIError) Error() string {
	return fmt.Sprintf("[%d:%s] %s", e.StatusCode, e.Type, e.Message)
}

func intPtr(i int) *int {
	return &i
}

// Float32Ptr converts a float32 to a *float32 as a convenience
func float32Ptr(f float32) *float32 {
	return &f
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
	var result GPTAPIErrorResponse
	if err := json.Unmarshal(data, &result); err != nil {
		// if we can't decode the json error then create an unexpected error
		apiError := GPTAPIError{
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
