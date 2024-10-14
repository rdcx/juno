<template>
  <div class="p-6">
    <h2 class="text-4xl font-bold mb-6">Welcome to the Juno Dashboard</h2>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- Nodes Overview -->
      <div class="bg-stone-800 p-6 rounded-lg shadow-lg">
        <h3 class="text-xl font-semibold mb-4">Nodes Overview</h3>
        <p class="text-sm">Total Active Nodes: <span class="font-bold">{{ activeNodes }}</span></p>
        <p class="text-sm">Total Inactive Nodes: <span class="font-bold">{{ inactiveNodes }}</span></p>
        <button @click="goToNodes" class="mt-4 px-4 py-2 bg-stone-600 rounded hover:bg-stone-500">
          Manage Nodes
        </button>
      </div>

      <!-- Balancers Overview -->
      <div class="bg-stone-800 p-6 rounded-lg shadow-lg">
        <h3 class="text-xl font-semibold mb-4">Balancers Overview</h3>
        <p class="text-sm">Total Active Balancers: <span class="font-bold">{{ activeBalancers }}</span></p>
        <p class="text-sm">Total Inactive Balancers: <span class="font-bold">{{ inactiveBalancers }}</span></p>
        <button @click="goToBalancers" class="mt-4 px-4 py-2 bg-stone-600 rounded hover:bg-stone-500">
          Manage Balancers
        </button>
      </div>

      <!-- Logs Overview -->
      <div class="bg-stone-800 p-6 rounded-lg shadow-lg">
        <h3 class="text-xl font-semibold mb-4">System Logs</h3>
        <p class="text-sm">Total Logs: <span class="font-bold">{{ totalLogs }}</span></p>
        <button @click="goToLogs" class="mt-4 px-4 py-2 bg-stone-600 rounded hover:bg-stone-500">
          View Logs
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { useRouter } from 'vue-router';

export default defineComponent({
  name: 'HomeView',
  setup() {
    const router = useRouter();

    // Dummy data for nodes and balancers
    const activeNodes = 12;
    const inactiveNodes = 3;
    const activeBalancers = 5;
    const inactiveBalancers = 1;
    const totalLogs = 125;

    const goToNodes = () => {
      router.push('/nodes');
    };

    const goToBalancers = () => {
      router.push('/balancers');
    };

    const goToLogs = () => {
      router.push('/logs');
    };

    return {
      activeNodes,
      inactiveNodes,
      activeBalancers,
      inactiveBalancers,
      totalLogs,
      goToNodes,
      goToBalancers,
      goToLogs,
    };
  },
});
</script>

<style scoped>
h2 {
  color: #b39ddb;
}
</style>