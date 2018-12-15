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

	// Config
	ConfigFile       string
	EncryptPasswords bool
	EncryptKey       string

	// Collector Settings
	Concurrency       int
	Interval          time.Duration
	Timeout           time.Duration
	AllowInsecureSSH  bool
	DisableCollection bool
	Daemonize         bool

	// Git
	GitPush       bool
	GitAuth       bool
	GitUser       string
	GitPrivateKey string

	// User Settings
	Workspace   string
	RunResult   string
	DefaultUser string
	DefaultPass string

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

	// Webserver
	WebserverEnabled bool
	HTTPListen       string

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
	UnlockPass     string
	Port           int
	Ciphers        string
	Disabled       bool
	CustomTimeout  time.Duration
	CommandTimeout time.Duration
	// Channel to track device collection failures
	FailChan chan bool
	// Number of failures to trigger an alert. A value of 0 disables alerts
	FailureWarning int
	// SSH Auth method
	SSHAuthMethod string
	SSHPrivateKey string
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
		SMTPHost:         "smtphost",
		SMTPPort:         25,
		HTTPListen:       "localhost:5002",
		GitAuth:          true,
		GitUser:          "git",
		GitPrivateKey:    ".ssh/id_rsa",
	}
}
