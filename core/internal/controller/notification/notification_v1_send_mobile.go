package notification

import (
	"context"

	v1 "billionmail-core/api/notification/v1"
	"billionmail-core/internal/service/notification"
	"billionmail-core/internal/service/public"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) SendMobile(ctx context.Context, req *v1.SendMobileReq) (res *v1.SendMobileRes, err error) {
	res = &v1.SendMobileRes{}

	if err := notification.Notification().SendMobile(ctx, req.Target, req.Title, []byte(req.Body)); err != nil {
		res.SetError(gerror.Newf(public.LangCtx(ctx, "Failed to send mobile notification: {}", err.Error())))
		return res, nil
	}

	res.SetSuccess(public.LangCtx(ctx, "Notification sent"))
	return res, nil
}
