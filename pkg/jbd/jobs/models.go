package jobs

import "time"

type JobsResponse struct {
	Jobs   []Jobs `json:"jobs"`
	Count  int    `json:"count"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type TaskDetail struct {
	Git                           string `json:"git"`
	BuildEnv                      string `json:"buildEnv"`
	Pool                          string `json:"pool"`
	LocalConfig                   string `json:"localConfig"`
	CacheNpm                      bool   `json:"cacheNpm"`
	UseJukwaapackCloudForCodePath bool   `json:"useJukwaapackCloudForCodePath"`
	SourceMap                     bool   `json:"sourceMap"`
	Codepath                      string `json:"codepath"`
	Branch                        string `json:"branch"`
	App                           string `json:"app"`
	Category                      string `json:"category"`
}

type Tasks struct {
	TaskType   string     `json:"taskType"`
	TaskDetail TaskDetail `json:"taskDetail"`
}

type Jobs struct {
	Tasks     []Tasks   `json:"tasks"`
	Status    string    `json:"status"`
	ID        string    `json:"_id"`
	Issuer    string    `json:"issuer"`
	JobType   string    `json:"jobType"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	V         int       `json:"__v"`
	SlackTs   string    `json:"slackTs"`
}
