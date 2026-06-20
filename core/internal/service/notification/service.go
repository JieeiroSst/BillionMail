package notification

import (
	"context"

	"billionmail-core/internal/model/entity"
)

type INotification interface {
	// SendMobile forwards a notification to the external mobile backend.
	// body is sent as a JSON array of byte values, not a base64 string.
	SendMobile(ctx context.Context, target, title string, body []byte) error

	// CreateBatch creates an async batch task that broadcasts the same title/body
	// to many targets. Sending happens in the background, see ProcessBatchTasks.
	CreateBatch(ctx context.Context, title string, body []byte, targets []string, threads int) (taskId int64, err error)

	// GetBatchTaskInfo returns the current progress of a batch task.
	GetBatchTaskInfo(ctx context.Context, taskId int64) (*entity.NotificationTask, error)

	// ListBatchTasks returns a paginated list of batch tasks, newest first.
	ListBatchTasks(ctx context.Context, page, pageSize int, keyword string) (tasks []*entity.NotificationTask, total int, err error)
}

var iNotificationService INotification

func Notification() INotification {
	if iNotificationService == nil {
		iNotificationService = newNotificationService()
	}
	return iNotificationService
}
