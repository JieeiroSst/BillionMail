package notification

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"billionmail-core/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	defaultBatchThreads = 5
	maxBatchThreads     = 50
	batchTargetPageSize = 100
)

// activeBatchTasks tracks task IDs currently being processed, so ProcessBatchTasks
// (polled every few seconds) doesn't start a second goroutine for the same task.
var activeBatchTasks sync.Map

// normalizeTargets trims, drops empty entries, and de-duplicates a target list.
func normalizeTargets(targets []string) []string {
	result := make([]string, 0, len(targets))
	seen := make(map[string]struct{}, len(targets))

	for _, t := range targets {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		result = append(result, t)
	}

	return result
}

func (s *notificationService) CreateBatch(ctx context.Context, title string, body []byte, targets []string, threads int) (taskId int64, err error) {
	cleanTargets := normalizeTargets(targets)
	if len(cleanTargets) == 0 {
		return 0, fmt.Errorf("targets must not be empty")
	}

	if threads <= 0 {
		threads = defaultBatchThreads
	} else if threads > maxBatchThreads {
		threads = maxBatchThreads
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model("bm_notification_tasks").Data(g.Map{
			"title":        title,
			"body":         body,
			"target_count": len(cleanTargets),
			"threads":      threads,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		taskId = id

		rows := make([]g.Map, 0, len(cleanTargets))
		for _, t := range cleanTargets {
			rows = append(rows, g.Map{"task_id": taskId, "target": t})
		}

		_, err = tx.Model("bm_notification_targets").Data(rows).Insert()
		return err
	})

	return taskId, err
}

func (s *notificationService) ListBatchTasks(ctx context.Context, page, pageSize int, keyword string) ([]*entity.NotificationTask, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	model := g.DB().Model("bm_notification_tasks").Ctx(ctx)
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		model = model.WhereLike("title", "%"+keyword+"%")
	}

	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	var tasks []*entity.NotificationTask
	if total > 0 {
		err = model.OrderDesc("id").Page(page, pageSize).Scan(&tasks)
		if err != nil {
			return nil, 0, err
		}
	}

	return tasks, total, nil
}

func (s *notificationService) GetBatchTaskInfo(ctx context.Context, taskId int64) (*entity.NotificationTask, error) {
	var task entity.NotificationTask

	if err := g.DB().Model("bm_notification_tasks").Ctx(ctx).Where("id", taskId).Scan(&task); err != nil {
		return nil, err
	}

	if task.Id == 0 {
		return nil, fmt.Errorf("notification task %d not found", taskId)
	}

	return &task, nil
}

// ProcessBatchTasks picks up pending batch tasks and processes each in its own
// goroutine. Intended to be polled periodically (see internal/service/timers).
func ProcessBatchTasks(ctx context.Context) {
	var tasks []*entity.NotificationTask

	err := g.DB().Model("bm_notification_tasks").
		Where("task_process IN (0,1)").
		OrderAsc("id").
		Limit(5).
		Scan(&tasks)

	if err != nil {
		g.Log().Warning(ctx, "query pending notification tasks failed:", err)
		return
	}

	for _, task := range tasks {
		if _, alreadyRunning := activeBatchTasks.LoadOrStore(task.Id, struct{}{}); alreadyRunning {
			continue
		}

		task := task
		taskCtx := gctx.New()

		go func() {
			defer activeBatchTasks.Delete(task.Id)
			processBatchTask(taskCtx, task)
		}()
	}
}

func processBatchTask(ctx context.Context, task *entity.NotificationTask) {
	if task.TaskProcess == 0 {
		if _, err := g.DB().Model("bm_notification_tasks").Where("id", task.Id).Data(g.Map{"task_process": 1}).Update(); err != nil {
			g.Log().Error(ctx, "failed to mark notification task running:", err)
		}
	}

	threads := task.Threads
	if threads <= 0 {
		threads = defaultBatchThreads
	}

	sem := make(chan struct{}, threads)
	var wg sync.WaitGroup
	lastId := 0

	for {
		var targets []*entity.NotificationTarget

		err := g.DB().Model("bm_notification_targets").
			Where("task_id", task.Id).
			Where("is_sent", 0).
			Where("id > ?", lastId).
			OrderAsc("id").
			Limit(batchTargetPageSize).
			Scan(&targets)

		if err != nil {
			g.Log().Error(ctx, "failed to fetch notification targets:", err)
			break
		}

		if len(targets) == 0 {
			break
		}

		lastId = targets[len(targets)-1].Id

		ids := make([]int, len(targets))
		for i, t := range targets {
			ids[i] = t.Id
		}

		if _, err := g.DB().Model("bm_notification_targets").WhereIn("id", ids).Data(g.Map{"is_sent": 2}).Update(); err != nil {
			g.Log().Error(ctx, "failed to mark notification targets as fetched:", err)
		}

		for _, target := range targets {
			target := target

			sem <- struct{}{}
			wg.Add(1)

			go func() {
				defer wg.Done()
				defer func() { <-sem }()
				sendBatchTarget(ctx, task, target)
			}()
		}
	}

	wg.Wait()

	if _, err := g.DB().Model("bm_notification_tasks").Where("id", task.Id).Data(g.Map{"task_process": 2}).Update(); err != nil {
		g.Log().Error(ctx, "failed to mark notification task completed:", err)
	}
}

func sendBatchTarget(ctx context.Context, task *entity.NotificationTask, target *entity.NotificationTarget) {
	err := Notification().SendMobile(ctx, target.Target, task.Title, task.Body)
	now := time.Now().Unix()

	if err != nil {
		_, updateErr := g.DB().Model("bm_notification_targets").Where("id", target.Id).Data(g.Map{
			"is_sent":       3,
			"sent_time":     now,
			"error_message": err.Error(),
		}).Update()
		if updateErr != nil {
			g.Log().Error(ctx, "failed to record failed notification target:", updateErr)
		}

		if _, err := g.DB().Model("bm_notification_tasks").Where("id", task.Id).Increment("failed_count", 1); err != nil {
			g.Log().Error(ctx, "failed to increment notification task failed_count:", err)
		}

		return
	}

	if _, updateErr := g.DB().Model("bm_notification_targets").Where("id", target.Id).Data(g.Map{
		"is_sent":   1,
		"sent_time": now,
	}).Update(); updateErr != nil {
		g.Log().Error(ctx, "failed to record sent notification target:", updateErr)
	}

	if _, err := g.DB().Model("bm_notification_tasks").Where("id", task.Id).Increment("sent_count", 1); err != nil {
		g.Log().Error(ctx, "failed to increment notification task sent_count:", err)
	}
}
