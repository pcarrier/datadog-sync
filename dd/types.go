package dd

type id int64

type thresholds struct {
	OK       *float64 `json:"ok,omitempty" yaml:"ok,omitempty"`
	Critical *float64 `json:"critical,omitempty" yaml:"critical,omitempty"`
	Warning  *float64 `json:"warning,omitempty" yaml:"warning,omitempty"`
}

type options struct {
	NotifyNoData            bool                `json:"notify_no_data,omitempty" yaml:"notify_no_data,omitempty"`
	NotifyAudit             bool                `json:"notify_audit,omitempty" yaml:"notify_audit,omitempty"`
	IncludeTags             bool                `json:"include_tags,omitempty" yaml:"include_tags,omitempty"`
	NoDataTimeFrameMinutes  *float64            `json:"no_data_timeframe,omitempty" yaml:"no_data_timeframe,omitempty"`
	TimeoutHours            *float64            `json:"timeout_h,omitempty" yaml:"timeout_h,omitempty"`
	RenotifyIntervalMinutes *float64            `json:"renotify_interval,omitempty" yaml:"renotify_interval,omitempty"`
	EscalationMessage       string              `json:"escalation_message,omitempty" yaml:"escalation_message,omitempty"`
	Thresholds              thresholds          `json:"thresholds,omitempty" yaml:"thresholds,omitempty"`
	Silenced                map[string]*float64 `json:"silenced,omitempty" yaml:"silenced,omitempty"`
}

// Monitor represents a Datadog monitor definition as exposed through their API.
// ID is optional but always provided on API reads.
type Monitor struct {
	Type    string   `json:"type" yaml:"type"`
	Query   string   `json:"query" yaml:"query"`
	Name    string   `json:"name,omitempty" yaml:"name,omitempty"`
	Message string   `json:"message,omitempty" yaml:"message,omitempty"`
	Tags    []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Options options  `json:"options,omitempty" yaml:"options,omitempty"`
	ID      *id      `json:"id,omitempty" yaml:"id,omitempty"`
}
