package models

import "time"

// Parsed version with org/app/pipeline split
type PipelineRecord struct {
	Org        string        `json:"org"`
	App        string        `json:"app"`
	Pipeline   string        `json:"pipeline"`
	JobName    string        `json:"jobName"`
	Status     string        `json:"status"`
	Duration   time.Duration `json:"duration"`
	InitiateAt time.Time     `json:"initiateAt"`
	FinishAt   time.Time     `json:"finishAt"`
}

type PipelineRequest struct {
	RunSpecID  string `json:"runSpecId"`
	FinishAt   string `json:"finishAt"`
	JobName    string `json:"jobName"`
	Status     string `json:"status"`
	InitiateAt string `json:"initiateAt"`
}
