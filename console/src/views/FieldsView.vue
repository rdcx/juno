<script setup lang="ts">
import Layout from '@/components/Layout.vue';
import { SelectorService, FieldService } from '@/services';
import { ref } from 'vue';
import { useNotificationStore } from '@/stores/notification';
import type { Selector, ListResponse as SelectorListResponse } from '@/types/SelectorTypes';
import type { Field, ListResponse as FieldListResponse, CreateRequest, CreateResponse } from '@/types/FieldTypes';

const notificationStore = useNotificationStore();
const selectors = ref<Array<Selector>>([]);
const fields = ref<Array<Field>>([]);
const selectedSelectorId = ref('');
const fieldName = ref('');
const fieldType = ref('string'); // default type
const createError = ref('');
const creating = ref(false);

const fetchSelectors = async () => {
    const response: SelectorListResponse = await SelectorService.list();
    selectors.value = response.selectors;
};

const fetchFields = async () => {
    const response: FieldListResponse = await FieldService.list();
    fields.value = response.fields;
};

const createField = async () => {
    if (!selectedSelectorId.value) {
        createError.value = 'Please select a selector';
        return;
    }
    creating.value = true;
    const requestData: CreateRequest = { name: fieldName.value, selector_id: selectedSelectorId.value, visibility: 'private', type: fieldType.value };
    const res: CreateResponse = await FieldService.create(requestData);
    creating.value = false;

    if (res.status === 'error') {
        createError.value = res.message;
        return;
    }

    notificationStore.createNotification('Field created successfully', 'success');
    fieldName.value = '';
    selectedSelectorId.value = '';
    fieldType.value = 'string'; // reset to default
    fetchFields();
};

fetchSelectors();
fetchFields();

</script>

<template>
    <Layout>
        <div class="max-w-7xl">
            <div class="grid grid-cols-4 p-4 gap-4">
                <div class="col-span-2 border dark:border-gray-700 rounded-lg p-4">
                    <h2 class="text-lg font-semibold dark:text-gray-300">Create Field</h2>

                    <div v-if="createError" class="text-red-500 text-sm mt-2">{{ createError }}</div>

                    <select v-model="selectedSelectorId" class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2">
                        <option value="" disabled>Select a Selector</option>
                        <option v-for="selector in selectors" :key="selector.id" :value="selector.id">
                            {{ selector.name }}
                        </option>
                    </select>

                    <input v-model="fieldName" type="text" placeholder="Field Name"
                        class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2" />

                    <select v-model="fieldType" class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2">
                        <option value="string">String</option>
                        <option value="number">Number</option>
                    </select>

                    <button @click="createField" :disabled="creating"
                        class="w-full bg-emerald-800 text-white rounded-lg p-2 mt-2">Create Field</button>
                </div>
            </div>
            <div class="w-full p-4">
                <div class="w-full border rounded-lg text-gray-800 dark:text-gray-300 dark:border-gray-700 p-4">
                    Fields List

                    <table class="min-w-full divide-y divide-gray-700">
                        <thead>
                            <tr>
                                <th scope="col"
                                    class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">ID
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold dark:text-white">Name
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold dark:text-white">Selector ID
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold dark:text-white">Type
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-800">
                            <tr v-for="field in fields" :key="field.id">
                                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium dark:text-white sm:pl-0">
                                    {{ field.id.substring(0, 8) }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm dark:text-gray-300">
                                    {{ field.name }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm dark:text-gray-300">
                                    {{ field.selector_id }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm dark:text-gray-300">
                                    {{ field.type }}
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </Layout>
</template>
