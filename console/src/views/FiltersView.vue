<script setup lang="ts">
import Layout from '@/components/Layout.vue';
import { FieldService, FilterService } from '@/services';
import { ref } from 'vue';
import { useNotificationStore } from '@/stores/notification';
import type { Field, ListResponse as FieldListResponse } from '@/types/FieldTypes';
import type { Filter, ListResponse as FilterListResponse, CreateRequest, CreateResponse } from '@/types/FilterTypes';

const notificationStore = useNotificationStore();
const fields = ref<Array<Field>>([]);
const filters = ref<Array<Filter>>([]);
const selectedFieldId = ref('');
const filterName = ref('');
const filterType = ref('string_equals'); // default type
const filterValue = ref('');
const createError = ref('');
const creating = ref(false);

const fetchFields = async () => {
    const response: FieldListResponse = await FieldService.list();
    fields.value = response.fields;
};

const fetchFilters = async () => {
    const response: FilterListResponse = await FilterService.list();
    filters.value = response.filters;
};

const createFilter = async () => {
    if (!selectedFieldId.value) {
        createError.value = 'Please select a field';
        return;
    }
    creating.value = true;
    const requestData: CreateRequest = { 
        name: filterName.value, 
        value: filterValue.value, 
        type: filterType.value, 
        field_id: selectedFieldId.value 
    };
    const res: CreateResponse = await FilterService.create(requestData);
    creating.value = false;

    if (res.status === 'error') {
        createError.value = res.message;
        return;
    }

    notificationStore.createNotification('Filter created successfully', 'success');
    filterName.value = '';
    filterValue.value = '';
    selectedFieldId.value = '';
    filterType.value = 'string_equals'; // reset to default
    fetchFilters();
};

fetchFields();
fetchFilters();

</script>

<template>
    <Layout>
        <div class="max-w-7xl">
            <div class="grid grid-cols-4 p-4 gap-4">
                <div class="col-span-2 border dark:border-gray-700 rounded-lg p-4">
                    <h2 class="text-lg font-semibold dark:text-gray-300">Create Filter</h2>

                    <div v-if="createError" class="text-red-500 text-sm mt-2">{{ createError }}</div>

                    <select v-model="selectedFieldId" class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2">
                        <option value="" disabled>Select a Field</option>
                        <option v-for="field in fields" :key="field.id" :value="field.id">
                            {{ field.name }}
                        </option>
                    </select>

                    <input v-model="filterName" type="text" placeholder="Filter Name"
                        class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2" />

                    <select v-model="filterType" class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2">
                        <option value="string_equals">String Equals</option>
                        <option value="string_contains">String Contains</option>
                    </select>

                    <input v-model="filterValue" type="text" placeholder="Filter Value"
                        class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2" />

                    <button @click="createFilter" :disabled="creating"
                        class="w-full bg-emerald-800 text-white rounded-lg p-2 mt-2">Create Filter</button>
                </div>
            </div>
            <div class="w-full p-4">
                <div class="w-full border rounded-lg text-gray-800 dark:text-gray-300 dark:border-gray-700 p-4">
                    Filters List

                    <table class="min-w-full divide-y divide-gray-700">
                        <thead>
                            <tr>
                                <th scope="col"
                                    class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">ID
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold dark:text-white">Name
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold dark:text-white">Field ID
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold dark:text-white">Type
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold dark:text-white">Value
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-800">
                            <tr v-for="filter in filters" :key="filter.id">
                                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium dark:text-white sm:pl-0">
                                    {{ filter.id.substring(0, 8) }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm dark:text-gray-300">
                                    {{ filter.name }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm dark:text-gray-300">
                                    {{ filter.field_id }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm dark:text-gray-300">
                                    {{ filter.type }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm dark:text-gray-300">
                                    {{ filter.value }}
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </Layout>
</template>
