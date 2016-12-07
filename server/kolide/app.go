package kolide

import "golang.org/x/net/context"

// AppConfigStore contains method for saving and retrieving
// application configuration
type AppConfigStore interface {
	NewAppConfig(info *AppConfig) (*AppConfig, error)
	AppConfig() (*AppConfig, error)
	SaveAppConfig(info *AppConfig) error
}

// AppConfigService provides methods for configuring
// the Kolide application
type AppConfigService interface {
	NewAppConfig(ctx context.Context, p AppConfigPayload) (info *AppConfig, err error)
	AppConfig(ctx context.Context) (info *AppConfig, err error)
	ModifyAppConfig(ctx context.Context, p AppConfigPayload) (info *AppConfig, err error)
}

// SMTPAuthenticationType defines whether or not a password will be required
// for STMP
type SMTPAuthenticationType uint

const (
	UserNamePassword SMTPAuthenticationType = iota
	AuthNone
)

// STMPAuthenticationMethod values used to define method to authenticate
// user
type SMTPAuthenticationMethod uint

const (
	Plain SMTPAuthenticationMethod = iota
	Login
	GSSAPI
	DigestMD5
	MD5
	CramMD5
)

type SMTPConfig struct {
	// Configured is a flag that indicates if smtp has been successfully
	// tested with the settings provided by an admin user.
	Configured bool `json:"configured" db:"smtp_configured"`
	// SenderAddress is the email address that will appear in emails sent
	// from Kolide
	SenderAddress string `json:"sender_address" db:"smtp_sender_address"`
	// Server is the host name of the SMTP server Kolide will use to send mail
	Server string `json:"server" db:"smtp_server"`
	// Port port SMTP server will use
	Port uint `json:"port" db:"smtp_port"`
	// AuthenticationType type of authentication for SMTP
	AuthenticationType SMTPAuthenticationType `json:"authentication_type" db:"smtp_authentication_type"`
	// UserName must be provided if SMTPAuthenticationType is UserNamePassword
	UserName string `json:"user_name" db:"smtp_user_name"`
	// Password must be provided if SMTPAuthenticationType is UserNamePassword
	Password string `json:"password" db:"smtp_password"`
	// EnableSSLTLS whether to use SSL/TLS for SMTP
	EnableSSLTLS bool `json:"enable_ssl_tls" db:"smtp_enable_ssl_tls"`
	// SMTPAuthenticationMethod authentication method smtp server will use
	AuthenticationMethod SMTPAuthenticationMethod `json:"authentication_method" db:"smtp_authentication_method"`
	// Advanced SMTP Options
	// SMTPDomain optional domain for SMTP
	Domain string `json:"domain,omitempty" db:"smtp_domain"`
	// VerifySSLCerts defaults to true but can be turned off of self signed
	// SSL certs are used by the SMTP server
	VerifySSLCerts bool `json:"verify_ssl_certs" db:"smtp_verify_ssl_certs"`
	// EnableStartTLS detects of TLS is enabled on mail server and starts to use it (default true)
	EnableStartTLS bool `json:"enable_start_tls" db:"smtp_enable_start_tls"`
}

// AppConfig holds configuration about the Kolide application.
// AppConfig data can be managed by a Kolide API user.
type AppConfig struct {
	ID              int64  `json:"-"`
	OrgName         string `json:"org_name,omitempty" db:"org_name"`
	OrgLogoURL      string `json:"org_logo_url,omitempty" db:"org_logo_url"`
	KolideServerURL string `json:"kolide_server_url" db:"kolide_server_url"`
	SMTPConfig
}

// AppConfigPayload contains request and response format of
// the AppConfig struct.
type AppConfigPayload struct {
	OrgInfo        *OrgInfo        `json:"org_info,omitempty"`
	ServerSettings *ServerSettings `json:"server_settings,omitempty"`
}

// OrgInfo contains general info about the organization using Kolide.
type OrgInfo struct {
	OrgName    *string `json:"org_name,omitempty" db:"org_name"`
	OrgLogoURL *string `json:"org_logo_url,omitempty" db:"org_logo_url"`
}

// ServerSettings contains general settings about the kolide App.
type ServerSettings struct {
	KolideServerURL *string `json:"kolide_server_url,omitempty"`
}

type OrderDirection int

const (
	OrderAscending OrderDirection = iota
	OrderDescending
)

// ListOptions defines options related to paging and ordering to be used when
// listing objects
type ListOptions struct {
	// Which page to return (must be positive integer)
	Page uint
	// How many results per page (must be positive integer, 0 indicates
	// unlimited)
	PerPage uint
	// Key to use for ordering
	OrderKey string
	// Direction of ordering
	OrderDirection OrderDirection
}
