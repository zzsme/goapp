// Package utils provides HTTP utility functions for making common HTTP requests
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HttpClient is the default HTTP client used for requests
var HttpClient = &http.Client{}

// HttpGet performs an HTTP GET request to the specified URL
func HttpGet(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// HttpPostJson performs an HTTP POST request with a JSON body
func HttpPostJson(url string, body interface{}, headers map[string]string) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// HttpPutJson performs an HTTP PUT request with a JSON body
func HttpPutJson(url string, body interface{}, headers map[string]string) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// HttpDelete performs an HTTP DELETE request
func HttpDelete(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// HttpPatchJson performs an HTTP PATCH request with a JSON body
func HttpPatchJson(url string, body interface{}, headers map[string]string) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
