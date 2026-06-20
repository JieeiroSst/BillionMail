package database_initialization

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {

	registerHandler(func() {
		notificationSQLList := []string{

			`CREATE TABLE IF NOT EXISTS bm_notification_tasks (
                id SERIAL PRIMARY KEY,
                title VARCHAR(255) NOT NULL,
                body BYTEA NOT NULL,
                target_count INTEGER NOT NULL DEFAULT 0,
                sent_count INTEGER NOT NULL DEFAULT 0,
                failed_count INTEGER NOT NULL DEFAULT 0,
                threads INTEGER NOT NULL DEFAULT 5,
                task_process SMALLINT NOT NULL DEFAULT 0, -- 0: Pending, 1: Running, 2: Completed
                create_time INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
                update_time INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())
            )`,

			`CREATE TABLE IF NOT EXISTS bm_notification_targets (
                id SERIAL PRIMARY KEY,
                task_id INTEGER NOT NULL,
                target VARCHAR(255) NOT NULL,
                is_sent SMALLINT NOT NULL DEFAULT 0, -- 0: Pending, 2: Fetched, 1: Sent, 3: Failed
                error_message TEXT NOT NULL DEFAULT '',
                sent_time INTEGER NOT NULL DEFAULT 0,
                create_time INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
                FOREIGN KEY (task_id) REFERENCES bm_notification_tasks(id) ON DELETE CASCADE
            )`,

			`CREATE INDEX IF NOT EXISTS idx_bm_notification_tasks_task_process ON bm_notification_tasks(task_process)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_notification_targets_task_id ON bm_notification_targets(task_id)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_notification_targets_is_sent ON bm_notification_targets(is_sent)`,
			`CREATE INDEX IF NOT EXISTS idx_bm_notification_targets_task_sent ON bm_notification_targets(task_id, is_sent)`,
		}

		for _, sql := range notificationSQLList {
			_, err := g.DB().Exec(context.Background(), sql)
			if err != nil {
				g.Log().Error(context.Background(), "Failed to execute notification SQL:", err, sql)
				return
			}
		}

		g.Log().Info(context.Background(), "Notification tables initialized successfully")
	})
}
