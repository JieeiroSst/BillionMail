export interface NotificationTask {
	id: number
	title: string
	target_count: number
	sent_count: number
	failed_count: number
	task_process: number
	create_time: number
	update_time: number
}

export interface NotificationTaskParams {
	page: number
	page_size: number
	keyword: string
}
