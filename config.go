package main

type CloudWatchLogGroup struct {
	Name         string   `json:"name"`
	FlowLog      bool     `json:"flowlog"`
	GroupColumns []string `json:"group_columns,omitempty"`
}

type Config struct {
	CloudWatchLogs []CloudWatchLogGroup `json:"cloudwatch_logs"`
}
