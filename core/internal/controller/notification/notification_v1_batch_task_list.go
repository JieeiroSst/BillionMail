package notification

import (
	"context"

	v1 "billionmail-core/api/notification/v1"
	"billionmail-core/internal/service/notification"
	"billionmail-core/internal/service/public"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) BatchTaskList(ctx context.Context, req *v1.BatchTaskListReq) (res *v1.BatchTaskListRes, err error) {
	res = &v1.BatchTaskListRes{}

	tasks, total, err := notification.Notification().ListBatchTasks(ctx, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		res.SetError(gerror.Newf(public.LangCtx(ctx, "Failed to list batch notification tasks: {}", err.Error())))
		return res, nil
	}

	list := make([]*v1.BatchTaskInfoData, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, &v1.BatchTaskInfoData{
			Id:          task.Id,
			Title:       task.Title,
			TargetCount: task.TargetCount,
			SentCount:   task.SentCount,
			FailedCount: task.FailedCount,
			TaskProcess: task.TaskProcess,
			CreateTime:  task.CreateTime,
			UpdateTime:  task.UpdateTime,
		})
	}

	res.Data.Total = total
	res.Data.List = list
	res.SetSuccess(public.LangCtx(ctx, "OK"))
	return res, nil
}
