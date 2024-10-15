<template>
  <div>
    <h3 class="text-xl font-bold mb-4">Nodes Overview</h3>
    <table class="w-full table-auto">
      <thead>
        <tr class="bg-gray-700 text-left">
          <th class="p-2">Node ID</th>
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
        <tr v-for="node in nodes" :key="node.id" class="hover:bg-gray-600">
          <td class="p-2">{{ node.id }}</td>
          <td class="p-2">{{ node.address }}</td>
          <td class="p-2">{{ node.status }}</td>
          <td class="p-2">{{ node.load }}%</td>
          <td class="p-2">
            <button class="px-2 py-1 bg-gray-700 rounded hover:bg-gray-600">
              View
            </button>
            <button class="ml-2 px-2 py-1 bg-red-600 rounded hover:bg-red-500">
              Delete
            </button>
          </td>
        </tr>
        <tr v-if="!loading && nodes.length === 0" class="text-center">
          <td colspan="4" class="p-2">No nodes available</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import axios from 'axios';
import { useUserStore } from '@/stores/user';

interface Node {
  id: string;
  address: string;
  status: string;
  load: number;
}

export default defineComponent({
  name: 'ManageNodes',
  setup() {
    const nodes = ref<Node[]>([]);
    const loading = ref(true);

    // Function to fetch nodes data from the API
    const fetchNodes = async () => {
      try {
        const response = await axios.get('http://localhost:8080/nodes', {
          headers: {
            Authorization: useUserStore().token,
          },
        });
        nodes.value = response.data.nodes; // Assuming the API returns an array of nodes
      } catch (error) {
        console.error('Error fetching nodes:', error);
      } finally {
        loading.value = false;
      }
    };

    // Fetch nodes when the component is mounted
    onMounted(fetchNodes);

    return {
      nodes,
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
th, td {
  padding: 10px;
}
tbody tr {
  transition: background-color 0.2s ease;
}
</style>
