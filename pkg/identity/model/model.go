// Package model defines the canonical identity payload for me JSON/YAML output.
package model

import "time"

// SourceStatus is the per-provider run status.
type SourceStatus string

const (
	StatusOK          SourceStatus = "ok"
	StatusPartial     SourceStatus = "partial"
	StatusError       SourceStatus = "error"
	StatusUnavailable SourceStatus = "unavailable"
)

// Subject is normalized primary identity (often from osaccount).
type Subject struct {
	Username    string `json:"username" yaml:"username"`
	DisplayName string `json:"display_name" yaml:"display_name"`
	UID         string `json:"uid" yaml:"uid"`
	GID         string `json:"gid" yaml:"gid"`
	HomeDir     string `json:"home_dir" yaml:"home_dir"`
	Shell       string `json:"shell" yaml:"shell"`
}

// SourceEnvelope is one provider's contribution.
type SourceEnvelope struct {
	Name       string       `json:"name" yaml:"name"`
	Status     SourceStatus `json:"status" yaml:"status"`
	DurationMs int64        `json:"duration_ms" yaml:"duration_ms"`
	Data       any          `json:"data,omitempty" yaml:"data,omitempty"`
	Warnings   []string     `json:"warnings,omitempty" yaml:"warnings,omitempty"`
}

// Meta is run-level metadata.
type Meta struct {
	Hostname   string `json:"hostname" yaml:"hostname"`
	Timestamp  string `json:"timestamp" yaml:"timestamp"`
	DurationMs int64  `json:"duration_ms" yaml:"duration_ms"`
	BestEffort bool   `json:"best_effort" yaml:"best_effort"`
}

// ErrorEntry is a tracked issue (unknown source, provider note, etc.).
type ErrorEntry struct {
	Source  string `json:"source" yaml:"source"`
	Code    string `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
}

// Payload is the top-level canonical document for me --json/--yaml.
type Payload struct {
	Subject Subject          `json:"subject" yaml:"subject"`
	Sources []SourceEnvelope `json:"sources" yaml:"sources"`
	Meta    Meta             `json:"meta" yaml:"meta"`
	Errors  []ErrorEntry     `json:"errors" yaml:"errors"`
}

// OsAccountData is the osaccount provider data contract.
type OsAccountData struct {
	Username string `json:"username" yaml:"username"`
	UID      string `json:"uid" yaml:"uid"`
	GID      string `json:"gid" yaml:"gid"`
	HomeDir  string `json:"home_dir" yaml:"home_dir"`
	Shell    string `json:"shell" yaml:"shell"`
	// GroupIDs are supplementary group ids (including primary on Unix) when available.
	GroupIDs []string `json:"group_ids,omitempty" yaml:"group_ids,omitempty"`
	Groups   []string `json:"groups,omitempty" yaml:"groups,omitempty"`
}

// EnvContextData is the envcontext provider data contract.
type EnvContextData struct {
	SudoUser string     `json:"sudo_user,omitempty" yaml:"sudo_user,omitempty"`
	SudoUID  string     `json:"sudo_uid,omitempty" yaml:"sudo_uid,omitempty"`
	SSHUser  string     `json:"ssh_user,omitempty" yaml:"ssh_user,omitempty"`
	CI       *EnvCIData `json:"ci,omitempty" yaml:"ci,omitempty"`
}

// EnvCIData describes CI context when present.
type EnvCIData struct {
	IsCI     bool   `json:"is_ci" yaml:"is_ci"`
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"`
	Actor    string `json:"actor,omitempty" yaml:"actor,omitempty"`
}

// NetworkData is the network provider data contract.
type NetworkData struct {
	Hostname       string   `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	FQDN           string   `json:"fqdn,omitempty" yaml:"fqdn,omitempty"`
	Domain         string   `json:"domain,omitempty" yaml:"domain,omitempty"`
	Workgroup      string   `json:"workgroup,omitempty" yaml:"workgroup,omitempty"`
	LocalAddresses []string `json:"local_addresses,omitempty" yaml:"local_addresses,omitempty"`
}

// SysInfoData is the sysinfo provider contract (GOOS, GOARCH, friendly OS name/version).
type SysInfoData struct {
	Platform  string `json:"platform" yaml:"platform"`
	Arch      string `json:"arch" yaml:"arch"`
	OSName    string `json:"os_name,omitempty" yaml:"os_name,omitempty"`
	OSVersion string `json:"os_version,omitempty" yaml:"os_version,omitempty"`
}

// AuthProvidersData is the authproviders provider data contract.
type AuthProvidersData struct {
	Git   *GitAuthData   `json:"git,omitempty" yaml:"git,omitempty"`
	Cloud *CloudAuthData `json:"cloud,omitempty" yaml:"cloud,omitempty"`
}

// GitAuthData is local git identity hints.
type GitAuthData struct {
	UserName  string `json:"user_name,omitempty" yaml:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty" yaml:"user_email,omitempty"`
}

// CloudAuthData groups cloud hints (best-effort; often empty).
type CloudAuthData struct {
	AWS   *AWSAuth   `json:"aws,omitempty" yaml:"aws,omitempty"`
	GCP   *GCPAuth   `json:"gcp,omitempty" yaml:"gcp,omitempty"`
	Azure *AzureAuth `json:"azure,omitempty" yaml:"azure,omitempty"`
}

// AWSAuth is AWS-related hints when available.
type AWSAuth struct {
	Configured bool   `json:"configured" yaml:"configured"`
	AccountID  string `json:"account_id,omitempty" yaml:"account_id,omitempty"`
	ARN        string `json:"arn,omitempty" yaml:"arn,omitempty"`
}

// GCPAuth is GCP-related hints when available.
type GCPAuth struct {
	Configured bool   `json:"configured" yaml:"configured"`
	Account    string `json:"account,omitempty" yaml:"account,omitempty"`
	Project    string `json:"project,omitempty" yaml:"project,omitempty"`
}

// AzureAuth is Azure-related hints when available.
type AzureAuth struct {
	Configured     bool   `json:"configured" yaml:"configured"`
	TenantID       string `json:"tenant_id,omitempty" yaml:"tenant_id,omitempty"`
	SubscriptionID string `json:"subscription_id,omitempty" yaml:"subscription_id,omitempty"`
	User           string `json:"user,omitempty" yaml:"user,omitempty"`
}

// NowRFC3339 returns the current UTC time in RFC3339 format.
func NowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
