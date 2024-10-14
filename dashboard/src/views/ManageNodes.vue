<template>
    <div>
      <h3 class="text-xl font-bold mb-4">Nodes Overview</h3>
      <table class="w-full table-auto">
        <thead>
          <tr class="bg-stone-700 text-left">
            <th class="p-2">Node ID</th>
            <th class="p-2">Status</th>
            <th class="p-2">Load</th>
            <th class="p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="node in nodes" :key="node.id" class="hover:bg-stone-600">
            <td class="p-2">{{ node.id }}</td>
            <td class="p-2">{{ node.status }}</td>
            <td class="p-2">{{ node.load }}%</td>
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
  
  interface Node {
    id: string;
    status: string;
    load: number;
  }
  
  export default defineComponent({
    name: 'ManageNodes',
    data() {
      return {
        nodes: [
          { id: 'Node-1', status: 'Active', load: 70 },
          { id: 'Node-2', status: 'Active', load: 50 },
          { id: 'Node-3', status: 'Inactive', load: 0 },
        ] as Node[],
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
  