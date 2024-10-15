<template>
    <div>
      <h3 class="text-xl font-bold mb-4">System Logs</h3>
      
      <!-- Log Filters -->
      <div class="mb-4 flex justify-between items-center">
        <div>
          <button @click="filterLogs('All')" class="px-4 py-2 mr-2 bg-gray-700 rounded hover:bg-gray-600">All</button>
          <button @click="filterLogs('Info')" class="px-4 py-2 mr-2 bg-gray-700 rounded hover:bg-gray-600">Info</button>
          <button @click="filterLogs('Warning')" class="px-4 py-2 mr-2 bg-yellow-600 rounded hover:bg-yellow-500">Warning</button>
          <button @click="filterLogs('Error')" class="px-4 py-2 bg-red-600 rounded hover:bg-red-500">Error</button>
        </div>
        <button @click="clearLogs" class="px-4 py-2 bg-red-600 rounded hover:bg-red-500">Clear Logs</button>
      </div>
  
      <!-- Logs Table -->
      <table class="w-full table-auto">
        <thead>
          <tr class="bg-gray-700 text-left">
            <th class="p-2">Node</th>
            <th class="p-2">Timestamp</th>
            <th class="p-2">Level</th>
            <th class="p-2">Message</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in filteredLogs" :key="log.id" class="hover:bg-gray-600">
            <td class="p-2">{{ log.node }}</td>
            <td class="p-2">{{ log.timestamp }}</td>
            <td class="p-2">
              <span :class="getLogLevelClass(log.level)">
                {{ log.level }}
              </span>
            </td>
            <td class="p-2">{{ log.message }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent } from 'vue';
  
  interface Log {
    id: string;
    node: string; 
    timestamp: string;
    level: string;
    message: string;
  }
  
  export default defineComponent({
    name: 'ManageLogs',
    data() {
      return {
        logs: [
          { id: '1', node: 'node-1', timestamp: '2024-10-14 12:00:00', level: 'Info', message: 'Node started successfully' },
          { id: '2', node: 'node-2', timestamp: '2024-10-14 12:05:00', level: 'Warning', message: 'High CPU usage detected' },
          { id: '3', node: 'node-3', timestamp: '2024-10-14 12:10:00', level: 'Error', message: 'Balancer failed to respond' },
        ] as Log[],
        filter: 'All',
      };
    },
    computed: {
      filteredLogs(): Log[] {
        if (this.filter === 'All') {
          return this.logs;
        }
        return this.logs.filter(log => log.level === this.filter);
      },
    },
    methods: {
      filterLogs(level: string) {
        this.filter = level;
      },
      clearLogs() {
        this.logs = [];
      },
      getLogLevelClass(level: string) {
        switch (level) {
          case 'Info':
            return 'text-blue-400';
          case 'Warning':
            return 'text-yellow-400';
          case 'Error':
            return 'text-red-400';
          default:
            return '';
        }
      },
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
  