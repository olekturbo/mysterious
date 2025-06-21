package predictor

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type HTTP struct {
	client *http.Client
	url    string
}

type Payload struct {
	Message string `json:"message"`
}

type Response struct {
	Result string `json:"result"`
}

func NewHTTP(client *http.Client, url string) *HTTP {
	return &HTTP{
		client: client,
		url:    url,
	}
}

func (h *HTTP) Predict(ctx context.Context, text string) (string, error) {
	payload := Payload{
		Message: text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res Response
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return "", err
	}

	return res.Result, nil
}
