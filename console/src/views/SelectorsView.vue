<script setup lang="ts">
import Layout from '@/components/Layout.vue';
import { SelectorService } from '@/services'; // Assume SelectorService handles API requests
import { ref, computed } from 'vue';
import { useNotificationStore } from '@/stores/notification';
import type { Selector, ListResponse, CreateRequest, CreateResponse } from '@/types/SelectorTypes';

const notificationStore = useNotificationStore();
const selectors = ref<Array<Selector>>([]);
const name = ref('');
const value = ref('');
const createError = ref('');
const creating = ref(false);

const fetchSelectors = async () => {
    const response: ListResponse = await SelectorService.list();
    selectors.value = response.selectors;
};

const createSelector = async () => {
    creating.value = true;
    const requestData: CreateRequest = { name: name.value, value: value.value, visibility: 'private' };
    const res: CreateResponse = await SelectorService.create(requestData);
    creating.value = false;

    if (res.status == 'error') {
        createError.value = res.message;
        return;
    }

    notificationStore.createNotification('Selector created successfully', 'success');
    name.value = '';
    value.value = '';
    fetchSelectors();
};

fetchSelectors();

</script>

<template>
    <Layout>
        <div class="max-w-7xl">
            <div class="grid grid-cols-4 p-4 gap-4">
                <div class="col-span-2 border dark:border-gray-700 rounded-lg p-4">
                    <h2 class="text-lg font-semibold dark:text-gray-300">Create Selector</h2>

                    <div v-if="createError" class="text-red-500 text-sm mt-2">{{ createError }}</div>

                    <input v-model="name" type="text" placeholder="Selector Name"
                        class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2" />

                    <input v-model="value" type="text" placeholder="Selector Value"
                        class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2" />

                    <button @click="createSelector" :disabled="creating"
                        class="w-full bg-emerald-800 text-white rounded-lg p-2 mt-2">Create Selector</button>
                </div>
            </div>
            <div class="w-full p-4">
                <div class="w-full border rounded-lg text-gray-800 dark:text-gray-300 dark:border-gray-700 p-4">
                    Selectors List

                    <table class="min-w-full divide-y divide-gray-700">
                        <thead>
                            <tr>
                                <th scope="col"
                                    class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">ID
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold dark:text-white">Name
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold dark:text-white">Value
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-800">
                            <tr v-for="selector in selectors" :key="selector.id">
                                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium dark:text-white sm:pl-0">
                                    {{ selector.id.substring(0, 8) }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm dark:text-gray-300">
                                    {{ selector.name }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm dark:text-gray-300">
                                    {{ selector.value }}
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </Layout>
</template>
