package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	DefaultVersion = "v25.0"
	DefaultBaseUrl = "https://graph.facebook.com"
)

type Client struct {
	dio     *http.Client
	token   string
	phone   string
	version string
	baseUrl string
}

func NewClient(config Config) *Client {
	cfg := config.WithDefaults()
	dio := &http.Client{
		Timeout:   cfg.DioTimeout,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}

	return &Client{
		dio:     dio,
		token:   cfg.ApiToken,
		phone:   cfg.ApiPhone,
		version: cfg.ApiVersion,
		baseUrl: cfg.ApiBaseUrl,
	}
}

func (c *Client) Get(ctx context.Context, path string) error {
	method := http.MethodGet
	endpoint := c.endpoint(path)
	request, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return err
	}

	// t.decorate(request)
	return c.execute(request)
}

func (c *Client) Post(ctx context.Context, path string, body any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	log.Printf("b: %v", string(b))
	log.Printf("uri: %s", c.endpoint(path))

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(path), bytes.NewReader(b))
	if err != nil {
		return err
	}

	return c.execute(request)
}

func (c *Client) PhoneNumberId() string { return c.phone }

// HELPERS

func (c *Client) execute(req *http.Request) error {
	response, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= http.StatusBadRequest {

		// IMPRIMIR ERROR
		return fmt.Errorf("error: %s", string(body))
		//var envelope APIErrorEnvelope
		//if decodeErr := decodeStrict(body, &envelope); decodeErr == nil && envelope.Error.Message != "" {
		//	return &RequestError{StatusCode: response.StatusCode, Body: string(body), Err: envelope.Error}
		//}
		//return &RequestError{StatusCode: response.StatusCode, Body: string(body), Err: fmt.Errorf("unexpected error response")}
		//

	}

	fmt.Printf("error: %s", string(body))

	return nil

	//if out == nil || len(body) == 0 {
	//	return nil
	//}
	//
	//return decodeStrict(body, out)
}

func (c *Client) endpoint(path string) string {
	return fmt.Sprintf("%s/%s/%s", c.baseUrl, c.version, strings.TrimLeft(path, "/"))
}

// doRequest maneja la inyección del token Bearer y la ejecución HTTP
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	return c.dio.Do(req.WithContext(ctx))
}

//
//func (c *Client) messages() string {
//	return fmt.Sprintf("%s/%s/%s/messages", c.ApiBaseUrl, c.ApiVersion, c.ApiPhone)
//}
