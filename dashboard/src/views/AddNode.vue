<template>
    <div class="max-w-lg mx-auto p-6 bg-gray-900 text-white rounded-lg shadow-lg">
      <h2 class="text-3xl font-bold mb-4">Add New Node</h2>
  
      <form @submit.prevent="handleSubmit">
        <!-- Address Input -->
        <div class="mb-4">
          <label class="block mb-2">Address</label>
          <input v-model="address" type="text" class="w-full p-2 rounded bg-gray-800 text-white" required />
        </div>
  
        <!-- Status Select -->
        <div class="mb-4">
          <label class="block mb-2">Status</label>
          <select v-model="status" class="w-full p-2 rounded bg-gray-800 text-white" required>
            <option value="Active">Active</option>
            <option value="Inactive">Inactive</option>
          </select>
        </div>
  
        <!-- Load Input -->
        <div class="mb-4">
          <label class="block mb-2">Shard Offset</label>
          <input v-model="shardOffset" type="number" min="0" max="100000" class="w-full p-2 rounded bg-gray-800 text-white" required />
        </div>
  
        <!-- Shards Slider -->
        <div class="mb-4">
          <label class="block mb-2">Total Shards: {{ shards }}</label>
          <input v-model="shards" type="range" min="1" max="100000" class="w-full" />
        </div>
  
        <!-- Submit Button -->
        <button type="submit" class="w-full bg-green-600 p-2 rounded hover:bg-green-500">
          Add Node
        </button>
      </form>
  
      <!-- Success/Error Messages -->
      <p v-if="successMessage" class="mt-4 text-green-400">{{ successMessage }}</p>
      <p v-if="errorMessage" class="mt-4 text-red-400">{{ errorMessage }}</p>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref } from 'vue';
  import axios from 'axios';
  import { useUserStore } from '@/stores/user';
  
  export default defineComponent({
    name: 'AddNode',
    setup() {
      const address = ref('');
      const status = ref('Active');
      const load = ref(0);
      const shardOffset = ref(0);
      const shards = ref(10); // Default value for shards
      const successMessage = ref('');
      const errorMessage = ref('');
  
      const handleSubmit = async () => {

        let shardList = [];
        for (let i = 0; i < shards.value; i++) {
          shardList.push(i + shardOffset.value);
        } 
        try {
          await axios.post('http://localhost:8080/nodes', {
            address: address.value,
            shards: shardList, // Send the selected number of shards
            status: status.value,
            load: load.value,
          }, {
            headers: {
              Authorization: useUserStore().token,
            },
          });
  
          // Clear form fields
          address.value = '';
          status.value = 'Active';
          load.value = 0;
          shards.value = 10;
  
          // Set success message
          successMessage.value = 'Node added successfully!';
          errorMessage.value = ''; // Clear error message if previously set
        } catch (error) {
          console.error('Error adding node:', error);
          errorMessage.value = 'Failed to add node.';
          successMessage.value = ''; // Clear success message if previously set
        }
      };
  
      return {
        address,
        status,
        load,
        shards,
        successMessage,
        errorMessage,
        handleSubmit,
        shardOffset,
      };
    },
  });
  </script>
  
  <style scoped>
  form {
    display: flex;
    flex-direction: column;
  }
  
  input,
  select {
    margin-bottom: 1rem;
  }
  
  button {
    margin-top: 1rem;
  }
  </style>
  