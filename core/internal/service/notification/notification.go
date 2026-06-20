package notification

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

const mobilePushTimeout = 10 * time.Second

type notificationService struct{}

func newNotificationService() INotification {
	return &notificationService{}
}

// mobilePushPayload is the JSON body forwarded to the external mobile backend.
// Body is a []int rather than []byte so it serializes as a JSON number array.
type mobilePushPayload struct {
	Target string `json:"target"`
	Title  string `json:"title"`
	Body   []int  `json:"body"`
}

func (s *notificationService) SendMobile(ctx context.Context, target, title string, body []byte) error {
	url := os.Getenv("MOBILE_NOTIFY_URL")
	if url == "" {
		return fmt.Errorf("MOBILE_NOTIFY_URL is not configured")
	}

	bodyInts := make([]int, len(body))
	for i, v := range body {
		bodyInts[i] = int(v)
	}

	client := g.Client().Timeout(mobilePushTimeout).ContentJson()
	if token := os.Getenv("MOBILE_NOTIFY_TOKEN"); token != "" {
		client = client.SetHeader("Authorization", "Bearer "+token)
	}

	resp, err := client.Post(ctx, url, mobilePushPayload{
		Target: target,
		Title:  title,
		Body:   bodyInts,
	})
	if err != nil {
		return err
	}
	defer resp.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mobile backend responded with status %d: %s", resp.StatusCode, resp.ReadAllString())
	}

	return nil
}
