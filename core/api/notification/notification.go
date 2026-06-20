// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package notification

import (
	"context"

	"billionmail-core/api/notification/v1"
)

type INotificationV1 interface {
	SendMobile(ctx context.Context, req *v1.SendMobileReq) (res *v1.SendMobileRes, err error)
	CreateBatch(ctx context.Context, req *v1.CreateBatchReq) (res *v1.CreateBatchRes, err error)
	BatchTaskInfo(ctx context.Context, req *v1.BatchTaskInfoReq) (res *v1.BatchTaskInfoRes, err error)
	BatchTaskList(ctx context.Context, req *v1.BatchTaskListReq) (res *v1.BatchTaskListRes, err error)
}
