<template>
    <div>
      <h3 class="text-xl font-bold mb-4">Balancers Overview</h3>
      <table class="w-full table-auto">
        <thead>
          <tr class="bg-stone-700 text-left">
            <th class="p-2">Balancer ID</th>
            <th class="p-2">Status</th>
            <th class="p-2">Load</th>
            <th class="p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="balancer in balancers" :key="balancer.id" class="hover:bg-stone-600">
            <td class="p-2">{{ balancer.id }}</td>
            <td class="p-2">{{ balancer.status }}</td>
            <td class="p-2">{{ balancer.load }}%</td>
            <td class="p-2">
              <button class="px-2 py-1 bg-stone-700 rounded hover:bg-stone-600">
                View
              </button>
              <button class="ml-2 px-2 py-1 bg-red-600 rounded hover:bg-red-500">
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent } from 'vue';
  
  interface Balancer {
    id: string;
    status: string;
    load: number;
  }
  
  export default defineComponent({
    name: 'ManageBalancers',
    data() {
      return {
        balancers: [
          { id: 'Balancer-1', status: 'Active', load: 60 },
          { id: 'Balancer-2', status: 'Active', load: 45 },
          { id: 'Balancer-3', status: 'Inactive', load: 0 },
        ] as Balancer[],
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
  th, td {
    padding: 10px;
  }
  tbody tr {
    transition: background-color 0.2s ease;
  }
  </style>
  