<template>
	<div class="p-24px">
		<bt-table-layout>
			<template #toolsLeft>
				<n-button type="primary" @click="handleAdd">
					{{ t('notification.actions.add') }}
				</n-button>
			</template>
			<template #toolsRight>
				<bt-search
					v-model:value="tableParams.keyword"
					:width="320"
					:placeholder="t('notification.search.titlePlaceholder')"
					@search="() => fetchTable(true)">
				</bt-search>
			</template>
			<template #table>
				<n-data-table v-bind="tableProps" :columns="columns">
					<template #empty>
						<n-empty :description="t('notification.empty')" />
					</template>
				</n-data-table>
			</template>
			<template #pageRight>
				<bt-table-page v-bind="pageProps" @refresh="() => fetchTable()"> </bt-table-page>
			</template>
			<template #modal>
				<form-modal></form-modal>
			</template>
		</bt-table-layout>
	</div>
</template>

<script lang="tsx" setup>
import { DataTableColumns, NTag } from 'naive-ui'
import { formatTime } from '@/utils'
import { useDataTable } from '@/hooks/useDataTable'
import { useModal } from '@/hooks/modal/useModal'
import { getBatchTaskList } from '@/api/modules/notification'
import type { NotificationTask, NotificationTaskParams } from './interface'

import NotificationBatchForm from './components/NotificationBatchForm.vue'

const { t } = useI18n()

const { tableProps, pageProps, tableParams, fetchTable } = useDataTable<
	NotificationTask,
	NotificationTaskParams
>({
	params: {
		page: 1,
		page_size: 10,
		keyword: '',
	},
	fetchFn: getBatchTaskList,
})

const columns: DataTableColumns<NotificationTask> = [
	{
		key: 'title',
		title: t('notification.columns.title'),
		ellipsis: { tooltip: true },
	},
	{
		key: 'target_count',
		title: t('notification.columns.targetCount'),
		width: 100,
	},
	{
		key: 'sent_count',
		title: t('notification.columns.sentCount'),
		width: 100,
	},
	{
		key: 'failed_count',
		title: t('notification.columns.failedCount'),
		width: 100,
	},
	{
		key: 'task_process',
		title: t('notification.columns.status'),
		width: 110,
		render: row => {
			if (row.task_process === 0)
				return (
					<NTag size="small" bordered={false} type="warning">
						{t('notification.status.pending')}
					</NTag>
				)
			if (row.task_process === 1)
				return (
					<NTag size="small" bordered={false} type="info">
						{t('notification.status.running')}
					</NTag>
				)
			return (
				<NTag size="small" bordered={false} type="success">
					{t('notification.status.completed')}
				</NTag>
			)
		},
	},
	{
		key: 'create_time',
		title: t('notification.columns.createTime'),
		width: 170,
		render: row => formatTime(row.create_time),
	},
]

const [FormModal, formModalApi] = useModal({
	component: NotificationBatchForm,
	state: {
		refresh: () => fetchTable(true),
	},
})

const handleAdd = () => {
	formModalApi.open()
}

let pollTimer: ReturnType<typeof setInterval> | undefined

onMounted(() => {
	pollTimer = setInterval(() => {
		if (tableProps.value.data.some((row: NotificationTask) => row.task_process !== 2)) {
			fetchTable()
		}
	}, 5000)
})

onUnmounted(() => {
	if (pollTimer) clearInterval(pollTimer)
})
</script>
