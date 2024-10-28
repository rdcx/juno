<script setup lang="ts">
import { ref, onMounted } from 'vue';
import JobService from '@/services/job/JobService';
import StrategyService from '@/services/strategy/StrategyService';
import type { Job, ListResponse as JobListResponse, CreateRequest, CreateResponse } from '@/types/JobTypes';
import type { Strategy, ListResponse as StrategyListResponse } from '@/types/StrategyTypes';
import { useNotificationStore } from '@/stores/notification';
import Layout from '@/components/Layout.vue';

const notificationStore = useNotificationStore();
const jobs = ref<Array<Job>>([]);
const strategies = ref<Array<Strategy>>([]);
const selectedStrategyId = ref('');

const fetchJobs = async () => {
    const response: JobListResponse = await JobService.list();
    jobs.value = response.jobs;
};

const fetchStrategies = async () => {
    const response: StrategyListResponse = await StrategyService.list();
    strategies.value = response.strategies;
};

const createJob = async () => {
    if (!selectedStrategyId.value) {
        notificationStore.createNotification('Please select a strategy', 'error');
        return;
    }

    const payload: CreateRequest = { strategy_id: selectedStrategyId.value };
    const response: CreateResponse = await JobService.create(payload);

    if (response.status === 'error') {
        notificationStore.createNotification(response.message, 'error');
        return;
    }

    notificationStore.createNotification('Job created successfully', 'success');
    selectedStrategyId.value = '';
    fetchJobs(); // Refresh the job list after creation
};

onMounted(() => {
    fetchJobs();
    fetchStrategies();
});

</script>

<template>
    <Layout>
        <div class="max-w-7xl p-4">
            <!-- Create Job Form -->
            <div class="mb-6 p-4 border dark:border-gray-700 rounded-lg">
                <h3 class="text-lg font-semibold dark:text-gray-300">Create New Job</h3>
                <select v-model="selectedStrategyId"
                    class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2">
                    <option value="" disabled>Select a Strategy</option>
                    <option v-for="strategy in strategies" :key="strategy.id" :value="strategy.id">
                        {{ strategy.name }}
                    </option>
                </select>
                <button @click="createJob" class="bg-emerald-800 text-white rounded-lg p-2 mt-2 w-full">Create
                    Job</button>
            </div>

            <!-- Jobs List -->
            <div>
                <table class="min-w-full divide-y divide-gray-700 border dark:border-gray-700">
                    <thead>
                        <tr>
                            <th class="px-6 py-3 text-left text-sm font-semibold dark:text-gray-300">ID</th>
                            <th class="px-6 py-3 text-left text-sm font-semibold dark:text-gray-300">Strategy</th>
                            <th class="px-6 py-3 text-left text-sm font-semibold dark:text-gray-300">Status</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-800">
                        <tr v-for="job in jobs" :key="job.id">
                            <td class="px-6 py-4 text-sm dark:text-gray-300">{{ job.id }}</td>
                            <td class="px-6 py-4 text-sm dark:text-gray-300">
                                {{ strategies.find(strategy => strategy.id === job.strategy_id)?.name || 'Unknown' }}
                            </td>
                            <td class="px-6 py-4 text-sm flex items-center">
                                <span class="dark:text-gray-300">{{ job.status }}</span>
                                <span v-if="job.status === 'running'"
                                    class="ml-2 animate-spin rounded-full h-4 w-4 border-t-2 border-blue-500"></span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </Layout>
</template>

<style>
/* Spinner animation */
.animate-spin {
    animation: spin 1s linear infinite;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}
</style>
