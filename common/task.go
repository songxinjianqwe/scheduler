package common

import "time"

type Task struct {
	Id string `json:"id"`
	TaskType string `json:"taskType"`
	time time.Duration `json:"time"`
	script string `json:"script"`
	results []string `json:"results"`
 }