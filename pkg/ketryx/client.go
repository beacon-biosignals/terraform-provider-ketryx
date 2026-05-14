package ketryx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	KETRYX_API_URL string = "https://app.ketryx.com"
)

type KetryxAPIProject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type KetryxAPIRepository struct {
	AuthToken  string `json:"auth-token"`
	AuthUser   string `json:"auth-user"`
	MainRef    string `json:"main-ref"`
	ReleaseRef string `json:"release-ref"`
	Url        string `json:"url"`
}

type KetryxAPIProjectsGetResponse struct {
	Projects []KetryxAPIProject `json:"projects"`
}

type KetryxAPIProjectsPostResponse struct {
	Id string `json:"id"`
}

type KetryxAPIProjectsPostRequest struct {
	Name         string                `json:"name"`
	Repositories []KetryxAPIRepository `json:"repository"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Client
//
////////////////////////////////////////////////////////////////////////////////////////////////////

type Client struct {
	bearer string
	url    string
}

func NewClient(api_key string) (*Client, error) {
	return &Client{bearer: api_key, url: KETRYX_API_URL}, nil
}

func (c *Client) do(ctx context.Context, path, method string, response, request any) error {
	var (
		api_url  string
		err      error
		req      *http.Request
		req_body []byte
		resp     *http.Response

		client *http.Client = &http.Client{}
	)

	api_url, err = url.JoinPath(c.url, path)
	ctx = tflog.SetField(tflog.SetField(ctx, "url", api_url), "method", method)

	tflog.Trace(ctx, "Executing Ketryx API request")

	// Construct the request
	if request != nil {
		// This request has a request body
		req_body, err = json.Marshal(request)
		if err == nil {
			req, err = http.NewRequest(method, api_url, bytes.NewBuffer(req_body))
		}
	} else {
		// This request do not have a request body
		req, err = http.NewRequest(method, api_url, nil)
	}
	if err != nil {
		// There was in error in constructing the request
		tflog.Error(ctx, "Error executing Ketryx API request", map[string]any{
			"request":  fmt.Sprintf("%v", req),
			"response": nil,
		})
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.bearer)
	req.Header.Set("Content-Type", "application/json")

	// Query Ketryx
	resp, err = client.Do(req)
	if err != nil {
		// There was an error in making the request
		tflog.Error(ctx, "Error in Ketryx API request", map[string]any{
			"request":  fmt.Sprintf("%v", req),
			"response": fmt.Sprintf("%v", resp),
		})
		return err
	}
	if resp.StatusCode != http.StatusOK {
		tflog.Error(ctx, "Error in Ketryx API request", map[string]any{
			"request":  fmt.Sprintf("%v", req),
			"response": fmt.Sprintf("%v", resp),
		})
		return fmt.Errorf("Error in Ketryx API response: %s", resp.Status)
	}
	defer resp.Body.Close()

	tflog.Trace(ctx, "Recieved Ketryx API response", map[string]any{
		"request":  fmt.Sprintf("%v", req),
		"response": fmt.Sprintf("%v", resp),
	})

	// Decode into the response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tflog.Debug(ctx, "Could not decode Ketryx API response", map[string]any{
			"request":  fmt.Sprintf("%v", req),
			"response": fmt.Sprintf("%v", resp),
		})
		return err
	}

	return nil
}

func (c *Client) ProjectsGet(ctx context.Context) (KetryxAPIProjectsGetResponse, error) {
	var (
		projects KetryxAPIProjectsGetResponse
		err      error
	)

	err = c.do(ctx, "api/v1/projects", http.MethodGet, &projects, nil)
	return projects, err
}

func (c *Client) ProjectsPost(ctx context.Context, req KetryxAPIProjectsPostRequest) (KetryxAPIProjectsPostResponse, error) {
	var (
		project_id KetryxAPIProjectsPostResponse
		err        error
	)

	err = c.do(ctx, "api/v1/projects", http.MethodPost, &project_id, req)
	return project_id, err
}
