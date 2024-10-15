<template>
    <div>
      <h3 class="text-xl font-bold mb-4">Balancers Overview</h3>
      <table class="w-full table-auto">
        <thead>
          <tr class="bg-gray-700 text-left">
            <th class="p-2">Balancer ID</th>
            <th class="p-2">Address</th>
            <th class="p-2">Status</th>
            <th class="p-2">Load</th>
            <th class="p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading" class="text-center">
            <td colspan="4" class="p-2">Loading...</td>
          </tr>
          <tr v-for="balancer in balancers" :key="balancer.id" class="hover:bg-gray-600">
            <td class="p-2">{{ balancer.id }}</td>
            <td class="p-2">{{ balancer.address }}</td>
            <td class="p-2">{{ balancer.status }}</td>
            <td class="p-2">{{ balancer.load }}%</td>
            <td class="p-2">
              <button class="px-2 py-1 bg-gray-700 rounded hover:bg-gray-600">View</button>
              <button class="ml-2 px-2 py-1 bg-red-600 rounded hover:bg-red-500">Delete</button>
            </td>
          </tr>
          <tr v-if="!loading && balancers.length === 0" class="text-center">
            <td colspan="4" class="p-2">No balancers found</td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref, onMounted } from 'vue';
  import axios from 'axios';
  import { useUserStore } from '@/stores/user';
  
  interface Balancer {
    id: string;
    address: string;
    status: string;
    load: number;
  }
  
  export default defineComponent({
    name: 'ManageBalancers',
    setup() {
      const balancers = ref<Balancer[]>([]);
      const loading = ref(true);
  
      // Fetch balancers data when the component is mounted
      const fetchBalancers = async () => {
        try {
          const response = await axios.get('http://localhost:8080/balancers', {
            headers: {
              Authorization: useUserStore().token,
            },
          });
          balancers.value = response.data.balancers; // Assuming the data is in the expected format
        } catch (error) {
          console.error('Error fetching balancers:', error);
        } finally {
          loading.value = false;
        }
      };
  
      onMounted(() => {
        fetchBalancers();
      });
  
      return {
        balancers,
        loading,
      };
    },
  });
  </script>
  
  <style scoped>
  table {
    width: 100%;
    border-spacing: 0;
  }
  thead {
    background-color: #5e35b1;
  }
  th,
  td {
    padding: 10px;
  }
  tbody tr {
    transition: background-color 0.2s ease;
  }
  </style>
  