package entity

// NotificationTask Entity
type NotificationTask struct {
	Id          int    `json:"id"           dc:"Task ID"`
	Title       string `json:"title"        dc:"Notification Title"`
	Body        []byte `json:"-"            dc:"Notification Body (raw bytes)"`
	TargetCount int    `json:"target_count" dc:"Total Target Count"`
	SentCount   int    `json:"sent_count"   dc:"Sent Count"`
	FailedCount int    `json:"failed_count" dc:"Failed Count"`
	Threads     int    `json:"threads"      dc:"Worker Threads"`
	TaskProcess int    `json:"task_process" dc:"Task Status(0:Pending 1:Running 2:Completed)"`
	CreateTime  int    `json:"create_time"  dc:"Create Time"`
	UpdateTime  int    `json:"update_time"  dc:"Update Time"`
}

// NotificationTarget Entity
type NotificationTarget struct {
	Id           int    `json:"id"            dc:"Target Row ID"`
	TaskId       int    `json:"task_id"       dc:"Task ID"`
	Target       string `json:"target"        dc:"Target device token / recipient identifier"`
	IsSent       int    `json:"is_sent"       dc:"Send Status(0:Pending 2:Fetched 1:Sent 3:Failed)"`
	ErrorMessage string `json:"error_message" dc:"Error Message"`
	SentTime     int    `json:"sent_time"     dc:"Send Time"`
	CreateTime   int    `json:"create_time"   dc:"Create Time"`
}
