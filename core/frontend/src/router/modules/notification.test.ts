import { describe, it, expect, vi } from 'vitest'

vi.mock('@/router/constant', () => ({
  Layout: () => Promise.resolve({ template: '<router-view />' }),
}))

import route from './notification'

describe('notification route', () => {
  it('has correct path', () => {
    expect(route.path).toBe('/notification')
  })

  it('has correct meta', () => {
    expect(route.meta).toMatchObject({
      sort: 8,
      key: 'notification',
      title: 'Notification',
      titleKey: 'layout.menu.notification',
    })
  })

  it('has one child named Notification', () => {
    expect(route.children).toHaveLength(1)
    expect(route.children![0].name).toBe('Notification')
    expect(route.children![0].path).toBe('/notification')
  })
})
