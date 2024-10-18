<script setup lang="ts">
import Layout from '@/components/Layout.vue';
import { TokenService } from '@/services';
import { TransactionService } from '@/services';
import { computed, ref } from 'vue';
import { useNotificationStore } from '@/stores/notification';

const notificationStore = useNotificationStore();
const balance = ref(0);
const amount = ref(0);
const depositError = ref('');
const depositing = ref(false);

const fetchBalance = async () => {
    balance.value = (await TokenService.balance()).balance;
};

const deposit = async () => {
    depositing.value = true;
    const res = await TokenService.deposit(amount.value);
    depositing.value = false;

    if (res.status == 'error') {
        depositError.value = res.message;
        return;
    }

    amount.value = 0;

    notificationStore.createNotification('Deposit successful', 'success');

    fetchBalance();
    fetchTransactions();
};

const transactions = ref([]);

const fetchTransactions = async () => {
    transactions.value = (await TransactionService.list()).transactions;
};

const formattedBalance = computed(() => {
    return balance.value.toLocaleString();
});

fetchBalance();

fetchTransactions();

</script>

<template>
    <Layout>
        <div class="max-w-7xl">
            <div class="grid grid-cols-4 p-4 gap-4">
                <div class="col-span-2 border dark:border-gray-700 rounded-lg p-4">
                    <h2 class="text-lg font-semibold dark:text-gray-300">Tokens</h2>
                    <div class="text-xl dark:text-gray-300">{{ formattedBalance }}</div>
                </div>

                <div class="col-span-2 border dark:border-gray-700 rounded-lg p-4 ml-4">
                    <h2 class="text-lg font-semibold dark:text-gray-300">Deposit</h2>

                    <div v-if="depositError" class="text-red-500 text-sm mt-2">{{ depositError }}</div>

                    <input v-model="amount" type="number"
                        class="w-full border dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 rounded-lg p-2 mt-2" />

                    <button @click="deposit" :disabled="depositing"
                        class="w-full bg-emerald-800 text-white rounded-lg p-2 mt-2">Deposit</button>
                </div>
            </div>
            <div class="w-full p-4">
                <div class="w-full border rounded-lg dark:border-gray-700 p-4">
                    Transactions

                    <table class="min-w-full divide-y divide-gray-700">
                        <thead>
                            <tr>
                                <th scope="col"
                                    class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">ID
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">Amount
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">Type
                                </th>
                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">Meta</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-800">
                            <tr v-for="tran in transactions">
                                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
                                    {{ tran.id }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-300">
                                    {{ tran.amount.toLocaleString() }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-300">
                                    {{ tran.key }}
                                </td>
                                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-300">
                                    {{ tran.meta }}
                                </td>
                            </tr>

                            <!-- More people... -->
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </Layout>
</template>