import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

import type { Notification } from '@/types/SystemTypes'
import { rand, timestamp } from '@vueuse/core'

export const useNotificationStore = defineStore('notification', () => {
    const notifications = ref<Array<Notification>>([])
    
    const addNotification = (notification: Notification) => {
        notifications.value.push(notification)
    }
    
    const removeNotification = (id: number) => {
        const index = notifications.value.findIndex((n) => n.id === id)
        if (index !== -1) {
            notifications.value.splice(index, 1)
        }
    }

    const createNotification = (message: string, type: string) => {
        const id = timestamp() + rand(0, 20000)  
        addNotification({ id, message, type })
    }
    
    return {
        notifications: computed(() => notifications.value),
        addNotification,
        removeNotification,
        createNotification,
    }
})
