<template>
    <div class="fixed top-0 right-0 m-4">
        <div @click="notificationsStore.removeNotification(notification.id)" v-for="notification in notifications"
            :key="notification.id"
            class="mb-4 w-72 flex items-center border-2 justify-center gap-2 text-sm border rounded p-3 text-white cursor-pointer"
            :class="notification.type === 'error' ? 'bg-red-700 border-red-900' : 'bg-emerald-700 border-emerald-900'">
            <div v-if="notification.type === 'error'">
                <ExclamationTriangleIcon class="w-6 h-6" />
            </div>
            <div v-else>
                <CheckIcon class="w-6 h-6" />
            </div>
            {{ notification.message }}
            <div class="ml-auto">
                <XMarkIcon class="w-4 h-4" />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { ExclamationTriangleIcon } from '@heroicons/vue/24/solid';
import { CheckIcon } from '@heroicons/vue/16/solid';
import { XMarkIcon } from '@heroicons/vue/16/solid';
import { useNotificationStore } from '@/stores/notification';

const notificationsStore = useNotificationStore();

const notifications = ref(notificationsStore.notifications);
</script>