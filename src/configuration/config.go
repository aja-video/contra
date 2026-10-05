package configuration

import (
	"time"
)

// Config holds the general application settings.
type Config struct {
	// Debug
	Debug      bool
	Copyrights bool
	Quiet      bool
	Version    bool
	ConfigTest bool

	// Config
	ConfigFile       string
	EncryptPasswords bool
	EncryptKey       string

	// Collector Settings
	Concurrency int
	Interval    time.Duration
	Timeout     time.Duration
	Daemonize   bool

	// Git
	GitPush       bool
	GitAuth       bool
	GitUser       string
	GitPrivateKey string
	// GitVerifyHostKey verifies the remote SSH host key against known_hosts
	// during git push. Set to false to skip verification.
	GitVerifyHostKey bool

	// User Settings
	Workspace string
	RunResult string

	// Mail
	EmailEnabled bool
	EmailTo      string
	EmailFrom    string
	EmailSubject string

	// SMTP Plain Auth
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string

	// Devices
	Devices []DeviceConfig
}

// DeviceConfig holds the device specific settings.
type DeviceConfig struct {
	Name string
	Host string
	Type string
	User string
	Pass string
	// UnlockPass is used to access xtd-cli-mode on certain hp/comware devices, or enable for Cisco/etc..
	UnlockPass string
	Port       int
	Ciphers    string
	Disabled   bool
	// Number of consecutive failed runs before the first alert is sent. A value of 0 disables alerts.
	FailureWarning int
	// FailureBackoffCount is the number of alerts sent during a single failure episode before
	// backoff engages. Once reached, further alerts are rate limited to FailureBackoffInterval.
	FailureBackoffCount int
	// FailureBackoffInterval is the minimum spacing between alerts once backoff has engaged.
	FailureBackoffInterval time.Duration
	// SSH settings
	SSHTimeout       time.Duration
	SSHAuthMethod    string
	SSHPrivateKey    string
	AllowInsecureSSH bool
	// Route53 settings
	ZoneID string
	// http-json settings
	Endpoint string
}

// GetName provides a simple implementation for the Collector interface.
func (d DeviceConfig) GetName() string {
	return d.Name
}

// getConfigDefaults provides reasonable defaults for Contra!
func getConfigDefaults() *Config {
	return &Config{
		ConfigFile:       "contra.conf",
		EncryptPasswords: true,
		Concurrency:      30,
		Interval:         300 * time.Second,
		Timeout:          120 * time.Second,
		Workspace:        "/workspace",
		EmailSubject:     "Changes from Contra!",
		EmailEnabled:     false,
		SMTPHost:         "smtphost",
		SMTPPort:         25,
		GitAuth:          true,
		GitUser:          "git",
		GitPrivateKey:    ".ssh/id_rsa",
		GitVerifyHostKey: true,
	}
}
