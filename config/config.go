package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime configuration values (loaded from env / defaults).
type Config struct {
	ServiceName string
	Env         string

	HTTPPort string
	GRPCPort string

	// Postgres URL handling:
	// PostgresURL        - raw POSTGRES_URL value (as provided; may be empty)
	// PostgresURLMode    - mode: raw | build | merge | auto
	// PostgresURLFinal   - final DSN used for connecting
	// PostgresURLWarning - conflict warning when raw URL and parts differ (raw wins)
	PostgresURL        string
	PostgresURLMode    string
	PostgresURLFinal   string
	PostgresURLWarning string

	// Individual Postgres parts (always populated: either parsed from URL or from env/defaults)
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string

	PostgresMaxOpenConns int
	PostgresMaxIdleConns int
	PostgresConnLifetime time.Duration

	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
}

// Load builds a Config from environment variables (with sensible defaults).
func Load() (*Config, error) {
	cfg := &Config{
		ServiceName: getString("SERVICE_NAME", "pizza_service"),
		Env:         getString("ENV", "dev"),

		HTTPPort: getString("HTTP_PORT", "8080"),
		GRPCPort: getString("GRPC_PORT", "9090"),

		PostgresURL:          strings.TrimSpace(getString("POSTGRES_URL", "")),
		PostgresURLMode:      strings.ToLower(strings.TrimSpace(getString("POSTGRES_URL_MODE", ""))),
		PostgresHost:         getString("POSTGRES_HOST", "localhost"),
		PostgresPort:         getString("POSTGRES_PORT", "5434"),
		PostgresUser:         getString("POSTGRES_USER", "postgres"),
		PostgresPassword:     getString("POSTGRES_PASSWORD", "postgres"),
		PostgresDatabase:     getString("POSTGRES_DATABASE", getString("POSTGRES_DB", "postgres")),
		PostgresSSLMode:      getString("POSTGRES_SSLMODE", "disable"),
		PostgresMaxOpenConns: getInt("POSTGRES_MAX_OPEN_CONNS", 20),
		PostgresMaxIdleConns: getInt("POSTGRES_MAX_IDLE_CONNS", 10),

		JWTSecret: getString("JWT_SECRET", "dev_secret"),
	}

	// Normalize ports (ensure leading colon for net/http / gRPC listeners).
	cfg.HTTPPort = normalizePort(cfg.HTTPPort)
	cfg.GRPCPort = normalizePort(cfg.GRPCPort)

	// Parse Postgres connection max lifetime.
	rawLifetime := getString("POSTGRES_CONN_MAX_LIFETIME", "1h")
	dur, err := parseDurationFlexible(rawLifetime)
	if err != nil {
		return nil, fmt.Errorf("parse POSTGRES_CONN_MAX_LIFETIME: %w", err)
	}
	cfg.PostgresConnLifetime = dur

	// Parse JWT access TTL.
	rawAccess := getString("JWT_ACCESS_TTL", "1h")
	accessTTL, err := parseDurationFlexible(rawAccess)
	if err != nil {
		return nil, fmt.Errorf("parse JWT_ACCESS_TTL: %w", err)
	}
	cfg.JWTAccessTTL = accessTTL

	// Parse JWT refresh TTL (optional).
	if rawRefresh, ok := os.LookupEnv("JWT_REFRESH_TTL"); ok && rawRefresh != "" {
		refreshTTL, err := parseDurationFlexible(rawRefresh)
		if err != nil {
			return nil, fmt.Errorf("parse JWT_REFRESH_TTL: %w", err)
		}
		cfg.JWTRefreshTTL = refreshTTL
	}

	// Determine Postgres URL resolution mode.
	cfg.decideURLMode()

	// Build / parse final Postgres URL.
	if err := cfg.resolvePostgresURL(); err != nil {
		return nil, err
	}

	// Validate final configuration.
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// decideURLMode finalizes the PostgresURLMode (defaulting / auto logic).
func (c *Config) decideURLMode() {
	mode := c.PostgresURLMode
	switch mode {
	case "", "raw":
		mode = "raw"
	case "build", "merge", "auto":
		// allowed
	default:
		mode = "raw"
	}
	if mode == "auto" {
		if c.PostgresHost != "" && c.PostgresPort != "" && c.PostgresUser != "" &&
			c.PostgresPassword != "" && c.PostgresDatabase != "" {
			mode = "build"
		} else {
			mode = "raw"
		}
	}
	c.PostgresURLMode = mode
}

// resolvePostgresURL computes PostgresURLFinal based on the selected mode.
func (c *Config) resolvePostgresURL() error {
	switch c.PostgresURLMode {
	case "build":
		// Ignore raw POSTGRES_URL; build from parts.
		c.PostgresURLFinal = c.buildURLFromParts()

	case "merge":
		// Merge: parse raw, fill missing fields from parts, then rebuild.
		if c.PostgresURL == "" {
			c.PostgresURLFinal = c.buildURLFromParts()
			return nil
		}
		parsed, err := parsePostgresURL(c.PostgresURL)
		if err != nil {
			return fmt.Errorf("parse POSTGRES_URL (merge): %w", err)
		}
		if parsed.Host != "" {
			c.PostgresHost = parsed.Host
		}
		if parsed.Port != "" {
			c.PostgresPort = parsed.Port
		}
		if parsed.User != "" {
			c.PostgresUser = parsed.User
		}
		if parsed.Password != "" {
			c.PostgresPassword = parsed.Password
		}
		if parsed.Database != "" {
			c.PostgresDatabase = parsed.Database
		}
		if parsed.SSLMode != "" {
			c.PostgresSSLMode = parsed.SSLMode
		}
		c.PostgresURLFinal = c.buildURLFromParts()

	case "raw":
		// Use the raw URL if present; parse to populate parts.
		if c.PostgresURL == "" {
			c.PostgresURLFinal = c.buildURLFromParts()
			return nil
		}
		parsed, err := parsePostgresURL(c.PostgresURL)
		if err != nil {
			return fmt.Errorf("parse POSTGRES_URL: %w", err)
		}
		conflicts := c.detectConflicts(parsed)
		if len(conflicts) > 0 {
			c.PostgresURLWarning = "conflicts between POSTGRES_URL and individual vars: " +
				strings.Join(conflicts, "; ")
		}
		// Overwrite parts with parsed values (raw wins).
		if parsed.Host != "" {
			c.PostgresHost = parsed.Host
		}
		if parsed.Port != "" {
			c.PostgresPort = parsed.Port
		}
		if parsed.User != "" {
			c.PostgresUser = parsed.User
		}
		if parsed.Password != "" {
			c.PostgresPassword = parsed.Password
		}
		if parsed.Database != "" {
			c.PostgresDatabase = parsed.Database
		}
		if parsed.SSLMode != "" {
			c.PostgresSSLMode = parsed.SSLMode
		}
		c.PostgresURLFinal = c.PostgresURL

	default:
		// Fallback: build from parts.
		c.PostgresURLFinal = c.buildURLFromParts()
	}
	return nil
}

// parsedURL is an intermediate representation of a Postgres URL.
type parsedURL struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
}

// parsePostgresURL parses a Postgres DSN into components.
func parsePostgresURL(raw string) (*parsedURL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return nil, fmt.Errorf("unsupported scheme %s", u.Scheme)
	}
	res := &parsedURL{}
	if u.User != nil {
		res.User = u.User.Username()
		if pw, ok := u.User.Password(); ok {
			res.Password = pw
		}
	}
	res.Host = u.Hostname()
	res.Port = u.Port()
	res.Database = strings.TrimPrefix(u.Path, "/")
	q := u.Query()
	res.SSLMode = q.Get("sslmode")
	return res, nil
}

// detectConflicts finds mismatches between raw URL parts and env-provided parts.
func (c *Config) detectConflicts(p *parsedURL) []string {
	var conflicts []string
	if p.Host != "" && c.PostgresHost != "" && !strings.EqualFold(p.Host, c.PostgresHost) {
		conflicts = append(conflicts, fmt.Sprintf("host(%s!=%s)", p.Host, c.PostgresHost))
	}
	if p.Port != "" && c.PostgresPort != "" && p.Port != c.PostgresPort {
		conflicts = append(conflicts, fmt.Sprintf("port(%s!=%s)", p.Port, c.PostgresPort))
	}
	if p.User != "" && c.PostgresUser != "" && p.User != c.PostgresUser {
		conflicts = append(conflicts, fmt.Sprintf("user(%s!=%s)", p.User, c.PostgresUser))
	}
	if p.Database != "" && c.PostgresDatabase != "" && p.Database != c.PostgresDatabase {
		conflicts = append(conflicts, fmt.Sprintf("db(%s!=%s)", p.Database, c.PostgresDatabase))
	}
	if p.SSLMode != "" && c.PostgresSSLMode != "" && p.SSLMode != c.PostgresSSLMode {
		conflicts = append(conflicts, fmt.Sprintf("sslmode(%s!=%s)", p.SSLMode, c.PostgresSSLMode))
	}
	return conflicts
}

// buildURLFromParts constructs a Postgres DSN from individual fields.
func (c *Config) buildURLFromParts() string {
	// Format: postgres://user:pass@host:port/db?sslmode=...
	userPart := url.UserPassword(c.PostgresUser, c.PostgresPassword).String()
	var sb strings.Builder
	sb.WriteString("postgres://")
	sb.WriteString(userPart)
	sb.WriteString("@")
	sb.WriteString(c.PostgresHost)
	if c.PostgresPort != "" {
		sb.WriteString(":" + c.PostgresPort)
	}
	sb.WriteString("/")
	sb.WriteString(c.PostgresDatabase)

	params := url.Values{}
	if c.PostgresSSLMode != "" {
		params.Set("sslmode", c.PostgresSSLMode)
	}
	if q := params.Encode(); q != "" {
		sb.WriteString("?")
		sb.WriteString(q)
	}
	return sb.String()
}

// Validate performs basic semantic checks on the final configuration.
func (c *Config) Validate() error {
	if c.PostgresURLFinal == "" {
		return errors.New("final postgres url is empty")
	}
	if !strings.HasPrefix(c.PostgresURLFinal, "postgres://") &&
		!strings.HasPrefix(c.PostgresURLFinal, "postgresql://") {
		return fmt.Errorf("final postgres url must start with postgres:// got %s", c.PostgresURLFinal)
	}
	if c.JWTSecret == "" {
		return errors.New("jwt secret is empty")
	}
	if c.HTTPPort == "" {
		return errors.New("http port is empty")
	}
	if c.GRPCPort == "" {
		return errors.New("grpc port is empty")
	}
	if c.JWTAccessTTL <= 0 {
		return errors.New("jwt access ttl must be > 0")
	}
	envAllowed := map[string]struct{}{"dev": {}, "stage": {}, "prod": {}}
	if _, ok := envAllowed[strings.ToLower(c.Env)]; !ok {
		return fmt.Errorf("unsupported ENV=%s (allowed: dev|stage|prod)", c.Env)
	}
	return nil
}

// IsProd returns true if environment is production.
func (c *Config) IsProd() bool { return strings.EqualFold(c.Env, "prod") }

// IsDev returns true if environment is development.
func (c *Config) IsDev() bool { return strings.EqualFold(c.Env, "dev") }

// IsStage returns true if environment is staging.
func (c *Config) IsStage() bool { return strings.EqualFold(c.Env, "stage") }

// MaskedPostgresURL returns final DSN with password masked (for logging).
func (c *Config) MaskedPostgresURL() string {
	raw := c.PostgresURLFinal
	if raw == "" {
		return raw
	}
	if !(strings.HasPrefix(raw, "postgres://") || strings.HasPrefix(raw, "postgresql://")) {
		return raw
	}
	i := strings.Index(raw, "://")
	if i == -1 {
		return raw
	}
	rest := raw[i+3:]
	at := strings.Index(rest, "@")
	if at == -1 {
		return raw
	}
	cred := rest[:at]
	colon := strings.Index(cred, ":")
	if colon == -1 {
		return raw
	}
	masked := cred[:colon] + ":***"
	return raw[:i+3] + masked + "@" + rest[at+1:]
}

// parseDurationFlexible accepts either Go duration syntax (e.g. "1h30m")
// or a plain integer treated as seconds.
func parseDurationFlexible(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("empty duration string")
	}
	if isAllDigits(s) {
		sec, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(sec) * time.Second, nil
	}
	return time.ParseDuration(s)
}

// isAllDigits returns true if s contains only decimal digits.
func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// getString returns an env var or a default.
func getString(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

// getInt returns an env var parsed as int or a default.
func getInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return def
}

// normalizePort ensures an address begins with ":" if non-empty.
func normalizePort(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return p
	}
	if !strings.HasPrefix(p, ":") {
		return ":" + p
	}
	return p
}