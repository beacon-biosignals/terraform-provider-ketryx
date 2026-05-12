package ketryx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type KetryxAPIProjectsGetResponse struct {
	Projects []KetryxAPIProject `json:"projects"`
}

type KetryxAPIProject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	bearer string
	url    string
}

func NewClient(api_key string) (*Client, error) {
	return &Client{bearer: api_key, url: "https://app.ketryx.com"}, nil
}

func (c *Client) ProjectsGet(ctx context.Context) (KetryxAPIProjectsGetResponse, error) {
	var (
		api_url  string
		err      error
		projects KetryxAPIProjectsGetResponse
		req      *http.Request
		resp     *http.Response

		client *http.Client = &http.Client{}
		method string       = http.MethodGet
	)

	api_url, err = url.JoinPath(c.url, "api/v1/projects")
	ctx = tflog.SetField(tflog.SetField(ctx, "url", api_url), "method", method)

	tflog.Trace(ctx, "Executing Ketryx API request")

	// Query Ketryx for projects
	req, err = http.NewRequest(http.MethodGet, api_url, nil)
	req.Header.Set("Authorization", "Bearer "+c.bearer)
	if err != nil {
		// There was in error in constructing the request
		tflog.Error(ctx, "Error executing Ketryx API request", map[string]any{
			"request":  fmt.Sprintf("%v", req),
			"response": nil,
		})
		return KetryxAPIProjectsGetResponse{}, err
	}
	resp, err = client.Do(req)
	if err != nil {
		// There was an error in making the request
		tflog.Error(ctx, "Error in Ketryx API request", map[string]any{
			"request":  fmt.Sprintf("%v", req),
			"response": fmt.Sprintf("%v", resp),
		})
		return KetryxAPIProjectsGetResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		tflog.Error(ctx, "Error in Ketryx API request", map[string]any{
			"request":  fmt.Sprintf("%v", req),
			"response": fmt.Sprintf("%v", resp),
		})
		return KetryxAPIProjectsGetResponse{},
			fmt.Errorf("Error in Ketryx API response: %s", resp.Status)
	}
	defer resp.Body.Close()

	tflog.Trace(ctx, "Recieved Ketryx API response", map[string]any{
		"request":  fmt.Sprintf("%v", req),
		"response": fmt.Sprintf("%v", resp),
	})

	// Decode into the response
	if err = json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		tflog.Debug(ctx, "Could not decode Ketryx API response", map[string]any{
			"request":  fmt.Sprintf("%v", req),
			"response": fmt.Sprintf("%v", resp),
		})
		return KetryxAPIProjectsGetResponse{}, err
	}

	return projects, nil
}
