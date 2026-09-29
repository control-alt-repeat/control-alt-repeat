package xero

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	authorizeURL   = "https://login.xero.com/identity/connect/authorize"
	tokenURL       = "https://identity.xero.com/connect/token"
	connectionsURL = "https://api.xero.com/connections"

	defaultRedirectURI = "http://localhost:8765/callback"
	defaultScopes      = "offline_access accounting.transactions accounting.contacts accounting.settings.read"
)

// Config holds the Xero app credentials. Values come from environment variables so
// no secrets are ever written into the repository.
//
//	XERO_CLIENT_ID      (required)
//	XERO_CLIENT_SECRET  (required)
//	XERO_TENANT_ID      (optional - auto-detected when the app has one connection)
//	XERO_TOKEN_FILE     (optional - default ~/.config/car/xero-token.json)
//	XERO_SCOPES         (optional - space separated)
//	XERO_REDIRECT_URI   (optional - only used by `car xero login`)
type Config struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	TokenFile    string
	Scopes       string
	RedirectURI  string
}

func ConfigFromEnv() (Config, error) {
	cfg := Config{
		ClientID:     os.Getenv("XERO_CLIENT_ID"),
		ClientSecret: os.Getenv("XERO_CLIENT_SECRET"),
		TenantID:     os.Getenv("XERO_TENANT_ID"),
		TokenFile:    os.Getenv("XERO_TOKEN_FILE"),
		Scopes:       os.Getenv("XERO_SCOPES"),
		RedirectURI:  os.Getenv("XERO_REDIRECT_URI"),
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return cfg, errors.New("XERO_CLIENT_ID and XERO_CLIENT_SECRET must be set")
	}
	if cfg.TokenFile == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return cfg, err
		}
		cfg.TokenFile = filepath.Join(dir, "car", "xero-token.json")
	}
	if cfg.Scopes == "" {
		cfg.Scopes = defaultScopes
	}
	if cfg.RedirectURI == "" {
		cfg.RedirectURI = defaultRedirectURI
	}

	return cfg, nil
}

type storedToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	TenantID     string    `json:"tenant_id,omitempty"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// TokenSource hands out valid access tokens. It supports both Xero app types:
//   - Custom connection: client_credentials grant, no refresh token needed.
//   - Web app: refresh_token grant, seeded once by `car xero login`. Xero rotates
//     refresh tokens on every use, so the new one is persisted immediately.
type TokenSource struct {
	cfg   Config
	mu    sync.Mutex
	token storedToken
}

func NewTokenSource(cfg Config) (*TokenSource, error) {
	ts := &TokenSource{cfg: cfg}
	data, err := os.ReadFile(cfg.TokenFile)
	switch {
	case err == nil:
		if err := json.Unmarshal(data, &ts.token); err != nil {
			return nil, fmt.Errorf("reading token file %s: %w", cfg.TokenFile, err)
		}
	case !errors.Is(err, os.ErrNotExist):
		return nil, err
	}
	return ts, nil
}

func (ts *TokenSource) TenantID() string {
	if ts.cfg.TenantID != "" {
		return ts.cfg.TenantID
	}
	return ts.token.TenantID
}

func (ts *TokenSource) AccessToken(ctx context.Context) (string, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.token.AccessToken != "" && time.Now().Add(time.Minute).Before(ts.token.ExpiresAt) {
		return ts.token.AccessToken, nil
	}

	form := url.Values{}
	if ts.token.RefreshToken != "" {
		form.Set("grant_type", "refresh_token")
		form.Set("refresh_token", ts.token.RefreshToken)
	} else {
		form.Set("grant_type", "client_credentials")
	}

	resp, err := ts.requestToken(ctx, form)
	if err != nil {
		return "", err
	}
	if err := ts.save(resp); err != nil {
		return "", err
	}

	return ts.token.AccessToken, nil
}

func (ts *TokenSource) requestToken(ctx context.Context, form url.Values) (tokenResponse, error) {
	var out tokenResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return out, err
	}
	req.SetBasicAuth(ts.cfg.ClientID, ts.cfg.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}
	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("xero token request (%s) failed: %d %s", form.Get("grant_type"), resp.StatusCode, body)
	}

	return out, json.Unmarshal(body, &out)
}

func (ts *TokenSource) save(resp tokenResponse) error {
	ts.token.AccessToken = resp.AccessToken
	ts.token.ExpiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	if resp.RefreshToken != "" {
		ts.token.RefreshToken = resp.RefreshToken
	}

	if err := os.MkdirAll(filepath.Dir(ts.cfg.TokenFile), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ts.token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ts.cfg.TokenFile, data, 0o600)
}

// Login runs the OAuth2 authorization code flow for a Xero "Web app". The user
// opens the printed URL, approves access, then pastes the URL they were redirected
// to (it does not matter if the browser shows an error for localhost).
func Login(ctx context.Context, cfg Config, in io.Reader, out io.Writer) error {
	ts, err := NewTokenSource(cfg)
	if err != nil {
		return err
	}

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return err
	}
	state := hex.EncodeToString(stateBytes)

	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", cfg.ClientID)
	q.Set("redirect_uri", cfg.RedirectURI)
	q.Set("scope", cfg.Scopes)
	q.Set("state", state)

	fmt.Fprintf(out, "Open this URL, approve access, then paste the full URL you are redirected to:\n\n%s?%s\n\n> ", authorizeURL, q.Encode())

	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	redirected, err := url.Parse(strings.TrimSpace(line))
	if err != nil {
		return err
	}
	if redirected.Query().Get("state") != state {
		return errors.New("state mismatch - please retry the login")
	}
	code := redirected.Query().Get("code")
	if code == "" {
		return fmt.Errorf("no code in redirect URL (error: %s)", redirected.Query().Get("error"))
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", cfg.RedirectURI)

	resp, err := ts.requestToken(ctx, form)
	if err != nil {
		return err
	}
	if err := ts.save(resp); err != nil {
		return err
	}

	conns, err := ListConnections(ctx, ts)
	if err != nil {
		return err
	}
	for _, c := range conns {
		fmt.Fprintf(out, "Connected organisation: %s (%s)\n", c.TenantName, c.TenantID)
	}
	if len(conns) == 1 {
		ts.token.TenantID = conns[0].TenantID
		return ts.save(tokenResponse{AccessToken: ts.token.AccessToken, ExpiresIn: int(time.Until(ts.token.ExpiresAt).Seconds())})
	}
	fmt.Fprintln(out, "Multiple organisations connected - set XERO_TENANT_ID to choose one.")
	return nil
}

type Connection struct {
	TenantID   string `json:"tenantId"`
	TenantName string `json:"tenantName"`
	TenantType string `json:"tenantType"`
}

func ListConnections(ctx context.Context, ts *TokenSource) ([]Connection, error) {
	token, err := ts.AccessToken(ctx)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, connectionsURL, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("listing xero connections failed: %d %s", resp.StatusCode, body)
	}

	var conns []Connection
	return conns, json.Unmarshal(body, &conns)
}
