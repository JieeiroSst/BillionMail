package v1

import (
	"billionmail-core/utility/types/api_v1"
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

// ByteArray marshals as a JSON array of byte values (e.g. [72,105]) instead of
// the base64 string Go's encoding/json produces by default for []byte.
type ByteArray []byte

func (b ByteArray) MarshalJSON() ([]byte, error) {
	ints := make([]int, len(b))

	for i, v := range b {
		ints[i] = int(v)
	}

	return json.Marshal(ints)
}

func (b *ByteArray) UnmarshalJSON(data []byte) error {
	var ints []int

	if err := json.Unmarshal(data, &ints); err != nil {
		return err
	}

	out := make(ByteArray, len(ints))

	for i, v := range ints {
		if v < 0 || v > 255 {
			return fmt.Errorf("byte array value out of range at index %d: %d", i, v)
		}
		out[i] = byte(v)
	}

	*b = out
	return nil
}

type SendMobileReq struct {
	g.Meta        `path:"/notification/send_mobile" method:"post" tags:"Notification" summary:"Send a notification to the mobile app"`
	Authorization string    `json:"authorization" dc:"Authorization" in:"header"`
	Target        string    `json:"target" v:"required" dc:"Target device token / recipient identifier on the mobile backend"`
	Title         string    `json:"title" v:"required" dc:"Notification title"`
	Body          ByteArray `json:"body" v:"required" dc:"Notification message body, as an array of byte values"`
}

type SendMobileRes struct {
	api_v1.StandardRes
}

type CreateBatchReq struct {
	g.Meta        `path:"/notification/batch/create" method:"post" tags:"Notification" summary:"Create a batch notification task for many mobile targets"`
	Authorization string    `json:"authorization" dc:"Authorization" in:"header"`
	Title         string    `json:"title" v:"required" dc:"Notification title"`
	Body          ByteArray `json:"body" v:"required" dc:"Notification message body, as an array of byte values"`
	Targets       []string  `json:"targets" v:"required" dc:"Target device tokens / recipient identifiers"`
	Threads       int       `json:"threads" v:"min:0" dc:"Concurrent send workers" default:"5"`
}

type CreateBatchRes struct {
	api_v1.StandardRes
	Data struct {
		Id int64 `json:"id" dc:"Batch task ID"`
	} `json:"data"`
}

type BatchTaskInfoReq struct {
	g.Meta        `path:"/notification/batch/info" method:"get" tags:"Notification" summary:"Get batch notification task progress"`
	Authorization string `json:"authorization" dc:"Authorization" in:"header"`
	Id            int64  `json:"id" v:"required" dc:"Batch task ID"`
}

type BatchTaskInfoData struct {
	Id          int    `json:"id" dc:"Batch task ID"`
	Title       string `json:"title" dc:"Notification title"`
	TargetCount int    `json:"target_count" dc:"Total target count"`
	SentCount   int    `json:"sent_count" dc:"Sent count"`
	FailedCount int    `json:"failed_count" dc:"Failed count"`
	TaskProcess int    `json:"task_process" dc:"Task status (0:Pending 1:Running 2:Completed)"`
	CreateTime  int    `json:"create_time" dc:"Create time"`
	UpdateTime  int    `json:"update_time" dc:"Update time"`
}

type BatchTaskInfoRes struct {
	api_v1.StandardRes
	Data BatchTaskInfoData `json:"data"`
}

type BatchTaskListReq struct {
	g.Meta        `path:"/notification/batch/list" method:"get" tags:"Notification" summary:"List batch notification tasks"`
	Authorization string `json:"authorization" dc:"Authorization" in:"header"`
	Page          int    `json:"page" v:"required|min:1" dc:"page"`
	PageSize      int    `json:"page_size" v:"required|min:1" dc:"page size"`
	Keyword       string `json:"keyword" dc:"search keyword (matches title)"`
}

type BatchTaskListRes struct {
	api_v1.StandardRes
	Data struct {
		Total int                  `json:"total" dc:"total"`
		List  []*BatchTaskInfoData `json:"list" dc:"task list"`
	} `json:"data"`
}
