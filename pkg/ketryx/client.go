package ketryx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

type KetryxAPIProjectDeleteResponse struct {
	Id string `json:"id"`
}

type KetryxAPIProjectGetResponse struct {
	Id           string                `json:"id"`
	Name         string                `json:"name"`
	Repositories []KetryxAPIRepository `json:"repository"`
}

type KetryxAPIProjectPostResponse struct {
	Id string `json:"id"`
}

type KetryxAPIProjectsGetResponse struct {
	Projects []KetryxAPIProject `json:"projects"`
}

type KetryxAPIProjectsPostResponse struct {
	Id string `json:"id"`
}

type KetryxAPIProjectPostRequest struct {
	Repositories []KetryxAPIRepository `json:"repository"`
	Settings     string                `json:"settings"`
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

// Ketryx client
type Client struct {
	bearer string
	url    string
}

// Creates a new Ketryx client
func NewClient(api_key string) (*Client, error) {
	return &Client{bearer: api_key, url: KETRYX_API_URL}, nil
}

// Workhorse function for requests against Ketryx
func (c *Client) do(
	ctx context.Context,
	path, method string,
	expectedStatusCode int,
	response, request any,
) error {
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

	tflog.Debug(ctx, "Executing Ketryx API request")

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
	if resp.StatusCode != expectedStatusCode {
		tflog.Error(ctx, "Error in Ketryx API request", map[string]any{
			"request":  fmt.Sprintf("%v", req),
			"response": fmt.Sprintf("%v", resp),
		})
		return fmt.Errorf("Error in Ketryx API response: %s", resp.Status)
	}
	defer resp.Body.Close()

	tflog.Debug(ctx, "Recieved Ketryx API response", map[string]any{
		"request":  fmt.Sprintf("%v", req),
		"response": fmt.Sprintf("%v", resp),
	})

	// Decode the response, if one exists
	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			tflog.Debug(ctx, "Could not decode Ketryx API response", map[string]any{
				"request":  fmt.Sprintf("%v", req),
				"response": fmt.Sprintf("%v", resp),
			})
			return err
		}
	}

	return nil
}

// Deletes a Ketryx project
func (c *Client) ProjectDelete(ctx context.Context, projectId string) (KetryxAPIProjectDeleteResponse, error) {
	var (
		resp KetryxAPIProjectDeleteResponse
		err  error
	)

	tflog.Debug(ctx, "Deleting Ketryx Project", map[string]any{
		"ketryx_project_id": projectId,
	})

	err = c.do(ctx, fmt.Sprintf("api/v1/projects/%s", projectId), http.MethodDelete, http.StatusNoContent, &resp, nil)
	return resp, err
}

// Retrieves a Ketryx project
func (c *Client) ProjectGet(ctx context.Context, projectId string) (KetryxAPIProjectGetResponse, error) {
	var (
		resp KetryxAPIProjectGetResponse
		err  error
	)

	tflog.Debug(ctx, "Retrieving Ketryx Project", map[string]any{
		"ketryx_project_id": projectId,
	})

	err = c.do(ctx, fmt.Sprintf("api/v1/projects/%s", projectId), http.MethodGet, http.StatusOK, &resp, nil)
	return resp, err
}

// Updates a Ketryx project
func (c *Client) ProjectPost(
	ctx context.Context,
	req KetryxAPIProjectPostRequest,
	projectId string,
) (KetryxAPIProjectPostResponse, error) {
	var (
		resp KetryxAPIProjectPostResponse
		err  error
	)

	tflog.Debug(ctx, "Updating Ketryx Project", map[string]any{
		"ketryx_project_id": projectId,
	})

	err = c.do(ctx, fmt.Sprintf("api/v1/project/%s", projectId), http.MethodPost, http.StatusOK, &resp, req)
	return resp, err
}

// Lists Ketryx projects
func (c *Client) ProjectsGet(ctx context.Context) (KetryxAPIProjectsGetResponse, error) {
	var (
		resp KetryxAPIProjectsGetResponse
		err  error
	)

	tflog.Debug(ctx, "Retrieving Ketryx Projects")

	err = c.do(ctx, "api/v1/projects", http.MethodGet, http.StatusOK, &resp, nil)
	return resp, err
}

// Creates a Ketryx project
func (c *Client) ProjectsPost(ctx context.Context, req KetryxAPIProjectsPostRequest) (KetryxAPIProjectsPostResponse, error) {
	var (
		resp KetryxAPIProjectsPostResponse
		err  error
	)

	tflog.Debug(ctx, "Creating Ketryx Project")

	// XXX The Ketryx API returns 200 OK for project creation, not 201 CREATED
	err = c.do(ctx, "api/v1/projects", http.MethodPost, http.StatusOK, &resp, req)
	return resp, err
}
