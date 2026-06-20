import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
	path: '/notification',
	name: 'NotificationLayout',
	component: Layout,
	meta: {
		sort: 8,
		key: 'notification',
		title: 'Notification',
		titleKey: 'layout.menu.notification',
	},
	children: [
		{
			path: '/notification',
			name: 'Notification',
			meta: { title: 'Notification', titleKey: 'layout.menu.notification' },
			component: () => import('@/views/notification/index.vue'),
		},
	],
}

export default route
