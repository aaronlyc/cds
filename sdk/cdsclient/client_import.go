package cdsclient

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ovh/cds/sdk"

	"github.com/ovh/cds/sdk/exportentities"
)

func (c *client) PipelineImport(projectKey string, content io.Reader, format string, force bool) ([]string, error) {
	var url string
	url = fmt.Sprintf("/project/%s/import/pipeline?format=%s", projectKey, format)

	if force {
		url += "&forceUpdate=true"
	}

	btes, _, _, err := c.Request(context.Background(), "POST", url, content)
	messages := []string{}
	_ = json.Unmarshal(btes, &messages)
	return messages, err
}

func (c *client) ApplicationImport(projectKey string, content io.Reader, format string, force bool) ([]string, error) {
	var url string
	url = fmt.Sprintf("/project/%s/import/application", projectKey)
	if force {
		url += "?force=true"
	}

	mods := []RequestModifier{}
	switch format {
	case "json":
		mods = []RequestModifier{
			func(r *http.Request) {
				r.Header.Set("Content-Type", "application/json")
			},
		}
	case "yaml", "yml":
		mods = []RequestModifier{
			func(r *http.Request) {
				r.Header.Set("Content-Type", "application/x-yaml")
			},
		}
	default:
		return nil, exportentities.ErrUnsupportedFormat
	}

	btes, _, _, err := c.Request(context.Background(), "POST", url, content, mods...)
	messages := []string{}
	_ = json.Unmarshal(btes, &messages)

	return messages, err
}

func (c *client) EnvironmentImport(projectKey string, content io.Reader, format string, force bool) ([]string, error) {
	var url string
	url = fmt.Sprintf("/project/%s/import/environment", projectKey)
	if force {
		url += "?force=true"
	}

	mods := []RequestModifier{}
	switch format {
	case "json":
		mods = []RequestModifier{
			func(r *http.Request) {
				r.Header.Set("Content-Type", "application/json")
			},
		}
	case "yaml", "yml":
		mods = []RequestModifier{
			func(r *http.Request) {
				r.Header.Set("Content-Type", "application/x-yaml")
			},
		}
	default:
		return nil, exportentities.ErrUnsupportedFormat
	}

	btes, _, _, err := c.Request(context.Background(), "POST", url, content, mods...)
	messages := []string{}
	_ = json.Unmarshal(btes, &messages)

	return messages, err
}

// WorkerModelImport import a worker model via as code
func (c *client) WorkerModelImport(content io.Reader, format string, force bool) (*sdk.Model, error) {
	url := "/worker/model/import"
	if force {
		url += "?force=true"
	}

	var mods []RequestModifier
	switch format {
	case "json":
		mods = []RequestModifier{
			func(r *http.Request) {
				r.Header.Set("Content-Type", "application/json")
			},
		}
	case "yaml", "yml":
		mods = []RequestModifier{
			func(r *http.Request) {
				r.Header.Set("Content-Type", "application/x-yaml")
			},
		}
	default:
		return nil, exportentities.ErrUnsupportedFormat
	}

	btes, _, code, err := c.Request(context.Background(), "POST", url, content, mods...)
	if err != nil {
		return nil, err
	}

	if code >= 400 {
		return nil, fmt.Errorf("HTTP Status code %d", code)
	}

	var wm sdk.Model
	if err := json.Unmarshal(btes, &wm); err != nil {
		return nil, err
	}

	return &wm, nil
}

func (c *client) WorkflowImport(projectKey string, content io.Reader, format string, force bool) ([]string, error) {
	var url string
	url = fmt.Sprintf("/project/%s/import/workflows", projectKey)
	if force {
		url += "?force=true"
	}

	mods := []RequestModifier{}
	switch format {
	case "json":
		mods = []RequestModifier{
			func(r *http.Request) {
				r.Header.Set("Content-Type", "application/json")
			},
		}
	case "yaml", "yml":
		mods = []RequestModifier{
			func(r *http.Request) {
				r.Header.Set("Content-Type", "application/x-yaml")
			},
		}
	default:
		return nil, exportentities.ErrUnsupportedFormat
	}

	btes, _, _, err := c.Request(context.Background(), "POST", url, content, mods...)
	messages := []string{}
	_ = json.Unmarshal(btes, &messages)

	return messages, err
}

func (c *client) WorkflowPush(projectKey string, tarContent io.Reader, mods ...RequestModifier) ([]string, *tar.Reader, error) {
	url := fmt.Sprintf("/project/%s/push/workflows", projectKey)

	mods = append(mods,
		func(r *http.Request) {
			r.Header.Set("Content-Type", "application/tar")
		})

	btes, headers, code, err := c.Request(context.Background(), "POST", url, tarContent, mods...)
	if err != nil {
		return nil, nil, err
	}

	if code >= 400 {
		return nil, nil, fmt.Errorf("HTTP Status code %d", code)
	}

	messages := []string{}
	if err := json.Unmarshal(btes, &messages); err != nil {
		return nil, nil, err
	}

	wName := headers.Get(sdk.ResponseWorkflowNameHeader)
	if wName == "" {
		return messages, nil, nil
	}
	tarReader, err := c.WorkflowPull(projectKey, wName, mods...)
	if err != nil {
		return nil, nil, err
	}

	return messages, tarReader, nil
}
