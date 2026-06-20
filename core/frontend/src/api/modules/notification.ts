import { i18n } from '@/i18n'
import { instance } from '@/api'

const { t } = i18n.global

export interface CreateBatchParams {
	title: string
	body: number[]
	targets: string[]
	threads: number
}

/**
 * 创建批量通知任务
 */
export const createBatchTask = (data: CreateBatchParams) => {
	return instance.post('/notification/batch/create', data, {
		fetchOptions: {
			loading: t('notification.loading.creating'),
			successMessage: true,
		},
	})
}

/**
 * 批量通知任务列表
 */
export const getBatchTaskList = (params: { page: number; page_size: number; keyword: string }) => {
	return instance.get('/notification/batch/list', { params })
}

/**
 * 批量通知任务详情
 */
export const getBatchTaskInfo = (params: { id: number }) => {
	return instance.get('/notification/batch/info', { params })
}
