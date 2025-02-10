package payloadcontent

type ScanMetadata struct {
        ScannerVersion  string   `json:"scanner_version"`
        PoliciesVersion string  `json:"policies_version"`
        ScanningRules   []string `json:"scanning_rules"`
        ExcludedPaths   []string `json:"excluded_paths"`
}

type Summary struct {
        TotalVulnerabilities int            `json:"total_vulnerabilities"`
        SeverityCounts       map[string]int `json:"severity_counts"`
        FixableCount         int            `json:"fixable_count"`
        Compliant            bool           `json:"compliant"`
}

type Vulnerability struct {
        ID            string    `json:"id"`
        Severity      string    `json:"severity"`
        CVSS          float64   `json:"cvss"`
        Status        string    `json:"status"`
        PackageName   string    `json:"package_name"`
        CurrentVersion string   `json:"current_version"`
        FixedVersion  string    `json:"fixed_version"`
        Description   string    `json:"description"`
        PublishedDate string    `json:"published_date"`
        Link          string    `json:"link"`
        RiskFactors   []string  `json:"risk_factors"`
}


type ScanResults struct {
        ScanID           string          `json:"scan_id"`
        Timestamp        string          `json:"timestamp"`
        ScanStatus       string          `json:"scan_status"`
        ResourceType     string          `json:"resource_type"`
        ResourceName     string          `json:"resource_name"`
        Vulnerabilities  []Vulnerability `json:"vulnerabilities"`
        Summary          Summary         `json:"summary"`
        ScanMetadata     ScanMetadata    `json:"scan_metadata"`
}

type ScanArray struct {
	KeyScanResult ScanResults `json:"scanResults"`
}

//content for the scan API
type PostScan struct {
        Repo      string   `json:"repo"` 
        Files    []string `json:"files"`
} 

//content for the query API
type PostQuery struct {
        Filters struct {
                Severity string `json:"severity"`
        } `json:"filters"`
}

type Config struct {
	ServerAddr string `json:"server_addr" yaml:"server_addr"`
	ServerPort string `json:"server_port" yaml:"server_port"`
	GitMaxRetries int `json:"git_maxretries" yaml:"git_maxretries"`
}

