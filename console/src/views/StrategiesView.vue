<script setup lang="ts">
import { ref } from 'vue';
import Layout from '@/components/Layout.vue';
import StrategyService from '@/services/strategy/StrategyService';
import type { Strategy, ListResponse, CreateRequest, CreateResponse } from '@/types/StrategyTypes';
import StrategyEditor from '@/components/StrategyEditor.vue';
import { useNotificationStore } from '@/stores/notification';

const notificationStore = useNotificationStore();
const strategies = ref<Array<Strategy>>([]);
const strategyName = ref('');
const selectedStrategy = ref<Strategy | null>(null);
const showEditor = ref(false);

const fetchStrategies = async () => {
    const response: ListResponse = await StrategyService.list();
    strategies.value = response.strategies;
};

const createStrategy = async () => {
    if (!strategyName.value) {
        notificationStore.createNotification('Strategy name is required', 'error');
        return;
    }

    const payload: CreateRequest = { name: strategyName.value };
    const response: CreateResponse = await StrategyService.create(payload);

    if (response.status === 'error') {
        notificationStore.createNotification(response.message, 'error');
        return;
    }

    notificationStore.createNotification('Strategy created successfully', 'success');
    strategyName.value = '';
    fetchStrategies(); // Refresh the list to include the new strategy
};

const openEditor = (strategy: Strategy) => {
    selectedStrategy.value = strategy;
    showEditor.value = true;
};

const closeEditor = () => {
    showEditor.value = false;
    selectedStrategy.value = null;
};

fetchStrategies();

</script>

<template>
    <Layout>
        <div class="max-w-7xl p-4">
            <!-- Create Strategy Form -->
            <div class="mb-6 p-4 border dark:border-gray-700 rounded-lg">
                <h3 class="text-lg font-semibold dark:text-gray-300">Create New Strategy</h3>
                <input v-model="strategyName" type="text" placeholder="Strategy Name"
                       class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2" />
                <button @click="createStrategy" class="bg-emerald-800 text-white rounded-lg p-2 mt-2 w-full">
                    Create Strategy
                </button>
            </div>

            <!-- Strategies List -->
            <div>
                <table class="min-w-full divide-y divide-gray-700 border dark:border-gray-700">
                    <thead>
                        <tr>
                            <th class="px-6 py-3 text-left text-sm font-semibold dark:text-gray-300">ID</th>
                            <th class="px-6 py-3 text-left text-sm font-semibold dark:text-gray-300">Name</th>
                            <th class="px-6 py-3 text-left text-sm font-semibold dark:text-gray-300">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-800">
                        <tr v-for="strategy in strategies" :key="strategy.id">
                            <td class="px-6 py-4 text-sm dark:text-gray-300">{{ strategy.id }}</td>
                            <td class="px-6 py-4 text-sm dark:text-gray-300">{{ strategy.name }}</td>
                            <td class="px-6 py-4 text-sm">
                                <button @click="openEditor(strategy)" class="text-blue-500 hover:underline">
                                    Edit
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Popup editor for strategy editing -->
        <StrategyEditor v-if="showEditor" :strategy="selectedStrategy" @close="closeEditor" />
    </Layout>
</template>
