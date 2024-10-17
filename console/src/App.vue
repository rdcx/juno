<script setup lang="ts">
import Notifications from './components/Notifications.vue';
import ToggleMode from './components/ToggleMode.vue';
</script>

<template>
  <router-view />

  <Notifications />
  <ToggleMode /> 
</template>
