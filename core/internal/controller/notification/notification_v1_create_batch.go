package notification

import (
	"context"

	v1 "billionmail-core/api/notification/v1"
	"billionmail-core/internal/service/notification"
	"billionmail-core/internal/service/public"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) CreateBatch(ctx context.Context, req *v1.CreateBatchReq) (res *v1.CreateBatchRes, err error) {
	res = &v1.CreateBatchRes{}

	taskId, err := notification.Notification().CreateBatch(ctx, req.Title, []byte(req.Body), req.Targets, req.Threads)
	if err != nil {
		res.SetError(gerror.Newf(public.LangCtx(ctx, "Failed to create batch notification task: {}", err.Error())))
		return res, nil
	}

	res.Data.Id = taskId
	res.SetSuccess(public.LangCtx(ctx, "Batch notification task created"))
	return res, nil
}
