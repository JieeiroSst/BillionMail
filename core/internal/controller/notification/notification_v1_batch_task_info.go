package notification

import (
	"context"

	v1 "billionmail-core/api/notification/v1"
	"billionmail-core/internal/service/notification"
	"billionmail-core/internal/service/public"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) BatchTaskInfo(ctx context.Context, req *v1.BatchTaskInfoReq) (res *v1.BatchTaskInfoRes, err error) {
	res = &v1.BatchTaskInfoRes{}

	task, err := notification.Notification().GetBatchTaskInfo(ctx, req.Id)
	if err != nil {
		res.SetError(gerror.Newf(public.LangCtx(ctx, "Failed to get batch notification task: {}", err.Error())))
		return res, nil
	}

	res.Data = v1.BatchTaskInfoData{
		Id:          task.Id,
		Title:       task.Title,
		TargetCount: task.TargetCount,
		SentCount:   task.SentCount,
		FailedCount: task.FailedCount,
		TaskProcess: task.TaskProcess,
		CreateTime:  task.CreateTime,
		UpdateTime:  task.UpdateTime,
	}
	res.SetSuccess(public.LangCtx(ctx, "OK"))
	return res, nil
}
