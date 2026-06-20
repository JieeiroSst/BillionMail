<template>
	<modal :title="t('notification.form.title')" width="560">
		<div class="max-h-640px pt-8px overflow-auto">
			<bt-form ref="formRef" :model="form" :rules="rules" class="pr-12px">
				<n-form-item :label="t('notification.form.notifyTitle')" path="title">
					<n-input v-model:value="form.title" :placeholder="t('notification.form.notifyTitlePlaceholder')">
					</n-input>
				</n-form-item>
				<n-form-item :label="t('notification.form.message')" path="message">
					<n-input
						v-model:value="form.message"
						type="textarea"
						:rows="4"
						:placeholder="t('notification.form.messagePlaceholder')">
					</n-input>
				</n-form-item>
				<n-form-item :label="t('notification.form.targets')" path="targets">
					<n-input
						v-model:value="form.targets"
						type="textarea"
						:rows="6"
						:placeholder="t('notification.form.targetsPlaceholder')">
					</n-input>
				</n-form-item>
				<n-form-item :label="t('notification.form.threads')" path="threads">
					<n-input-number v-model:value="form.threads" :min="1" :max="50" class="w-full">
					</n-input-number>
				</n-form-item>
			</bt-form>
		</div>
	</modal>
</template>

<script lang="ts" setup>
import { FormRules } from 'naive-ui'
import { createBatchTask } from '@/api/modules/notification'
import { useModal } from '@/hooks/modal/useModal'

const { t } = useI18n()

const formRef = useTemplateRef('formRef')

const form = reactive({
	title: '',
	message: '',
	targets: '',
	threads: 5,
})

const rules: FormRules = {
	title: {
		required: true,
		message: t('notification.form.validation.titleRequired'),
		trigger: ['input', 'blur'],
	},
	message: {
		required: true,
		message: t('notification.form.validation.messageRequired'),
		trigger: ['input', 'blur'],
	},
	targets: {
		required: true,
		trigger: ['input', 'blur'],
		validator: (_rule, value: string) => {
			const lines = value
				.split('\n')
				.map(item => item.trim())
				.filter(item => item.length > 0)
			if (lines.length === 0) {
				return new Error(t('notification.form.validation.targetsRequired'))
			}
			return true
		},
	},
}

const resetForm = () => {
	form.title = ''
	form.message = ''
	form.targets = ''
	form.threads = 5
}

const getTargets = () => {
	return form.targets
		.split('\n')
		.map(item => item.trim())
		.filter(item => item.length > 0)
}

const [Modal, modalApi] = useModal({
	onChangeState: isOpen => {
		if (!isOpen) resetForm()
	},
	onConfirm: async () => {
		await formRef.value?.validate()

		await createBatchTask({
			title: form.title,
			body: Array.from(new TextEncoder().encode(form.message)),
			targets: getTargets(),
			threads: form.threads,
		})

		const state = modalApi.getState<{ refresh: () => void }>()
		state.refresh?.()
	},
})
</script>
