package main

type CloudWatchLogGroup struct {
	Name    string `json:"name"`
	FlowLog bool   `json:"flowlog"`
}

type Config struct {
	CloudWatchLogs []CloudWatchLogGroup `json:"cloudwatch_logs"`
}
