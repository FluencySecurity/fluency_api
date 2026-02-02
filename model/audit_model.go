package model

type AuditEvent struct {
	Account  string `json:"account"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	ClientIP string `json:"clientIP"`
	Country  string `json:"country"`
	City     string `json:"city"`

	CreatedOn  int64  `json:"createdOn"`
	Dur        int64  `json:"dur"`
	Category   string `json:"category"`
	Error      string `json:"error"`
	DayIndex   string `json:"dayIndex"`
	Partition  string `json:"partition"`
	Action     string `json:"action"`
	ActionType string `json:"action_type"`
	ActionInfo string `json:"action_info"`
}