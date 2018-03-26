package configuration

import (
	"time"
)

// Config holds the general application settings.
type Config struct {
	// Debug
	Debug      bool
	Copyrights bool

	// Config
	ConfigFile string

	// Collector Settings
	Concurrency       int
	Interval          time.Duration
	Timeout           time.Duration
	AllowInsecureSSH  bool
	DisableCollection bool

	// Git
	GitPush bool

	// User Settings
	Workspace   string
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
	Name           string
	Host           string
	Type           string
	User           string
	Pass           string
	Port           int
	Disabled       bool
	CustomTimeout  time.Duration
	CommandTimeout time.Duration
}

// GetName provides a simple implementation for the Collector interface.
func (d DeviceConfig) GetName() string {
	return d.Name
}

// getConfigDefaults provides reasonable defaults for Contra!
func getConfigDefaults() *Config {
	return &Config{
		false,
		false,
		"contra.conf",
		30,
		300 * time.Second,
		120 * time.Second,
		false,
		false,
		false,
		"/workspace",
		"",
		"",
		false,
		"",
		"",
		"Changes from Contra!",
		"smtphost",
		25,
		"",
		"",
		false,
		"localhost:5002",
		nil,
	}
}
