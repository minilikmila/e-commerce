package cloudinary

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type Client struct {
	CloudName    string
	APIKey       string
	APISecret    string
	UploadPreset string
	Folder       string
	HTTPClient   *http.Client
}

func NewClient(cloudName, apiKey, apiSecret, uploadPreset, folder string) *Client {
	// Use a custom transport with longer timeouts and better DNS handling
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &Client{
		CloudName:    cloudName,
		APIKey:       apiKey,
		APISecret:    apiSecret,
		UploadPreset: uploadPreset,
		Folder:       folder,
		HTTPClient: &http.Client{
			Timeout:   60 * time.Second,
			Transport: transport,
		},
	}
}

// UploadUnsigned uploads a file using an unsigned upload preset. Returns the secure_url.
func (c *Client) UploadUnsigned(ctx context.Context, file io.Reader, filename string) (string, error) {
	if c.UploadPreset == "" {
		return "", fmt.Errorf("upload preset required for unsigned upload")
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	_ = writer.WriteField("upload_preset", c.UploadPreset)
	if c.Folder != "" {
		_ = writer.WriteField("folder", c.Folder)
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", url.PathEscape(c.CloudName))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// Provide more context for DNS/network errors
		if netErr, ok := err.(net.Error); ok {
			if netErr.Timeout() {
				return "", fmt.Errorf("cloudinary upload timeout: %w", err)
			}
			if dnsErr, ok := netErr.(*net.DNSError); ok {
				return "", fmt.Errorf("cloudinary DNS resolution failed (check network/Docker DNS): %w", dnsErr)
			}
		}
		return "", fmt.Errorf("cloudinary upload network error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("cloudinary upload failed (status %d): %s", resp.StatusCode, string(b))
	}

	type uploadResp struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	}
	var ur uploadResp
	b, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &ur); err != nil {
		return "", fmt.Errorf("decode cloudinary response: %w", err)
	}
	if ur.SecureURL != "" {
		return ur.SecureURL, nil
	}
	if ur.URL != "" {
		return ur.URL, nil
	}
	return "", fmt.Errorf("cloudinary response missing url")
}

// UploadSigned uploads a file using signed parameters (api_key + signature + timestamp).
// Signature is computed as sha1 of the concatenated, sorted params and api secret, per Cloudinary spec.
func (c *Client) UploadSigned(ctx context.Context, file io.Reader, filename string, opts map[string]string) (string, error) {
	if c.APIKey == "" || c.APISecret == "" {
		return "", fmt.Errorf("api key/secret required for signed upload")
	}
	// base params
	params := map[string]string{}
	// optional folder
	if c.Folder != "" {
		params["folder"] = c.Folder
	}
	// merge opts
	for k, v := range opts {
		if v != "" {
			params[k] = v
		}
	}
	// mandatory timestamp
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	params["timestamp"] = ts

	// compute signature
	signature := c.sign(params)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	// file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	// params
	for k, v := range params {
		_ = writer.WriteField(k, v)
	}
	_ = writer.WriteField("api_key", c.APIKey)
	_ = writer.WriteField("signature", signature)

	if err := writer.Close(); err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", url.PathEscape(c.CloudName))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// Provide more context for DNS/network errors
		if netErr, ok := err.(net.Error); ok {
			if netErr.Timeout() {
				return "", fmt.Errorf("cloudinary upload timeout: %w", err)
			}
			if dnsErr, ok := netErr.(*net.DNSError); ok {
				return "", fmt.Errorf("cloudinary DNS resolution failed (check network/Docker DNS): %w", dnsErr)
			}
		}
		return "", fmt.Errorf("cloudinary upload network error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("cloudinary upload failed (status %d): %s", resp.StatusCode, string(b))
	}

	type uploadResp struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	}
	var ur uploadResp
	b, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &ur); err != nil {
		return "", fmt.Errorf("decode cloudinary response: %w", err)
	}
	if ur.SecureURL != "" {
		return ur.SecureURL, nil
	}
	if ur.URL != "" {
		return ur.URL, nil
	}
	return "", fmt.Errorf("cloudinary response missing url")
}

// sign computes SHA1 hex signature for provided params using API secret.
// Build string "key=value&..." sorted by key, then append api_secret, sha1 hex of the result.
func (c *Client) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b bytes.Buffer
	for i, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(params[k])
		if i < len(keys)-1 {
			b.WriteByte('&')
		}
	}
	b.WriteString(c.APISecret)
	sum := sha1.Sum(b.Bytes())
	return hex.EncodeToString(sum[:])
}
