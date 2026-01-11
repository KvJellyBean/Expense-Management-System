<template>
  <div class="flex flex-wrap gap-2 sm:gap-4 mb-4 sm:mb-6">
    <NuxtLink 
      to="/dashboard" 
      :class="isActive('/dashboard') ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100'"
      class="px-3 py-2 sm:px-4 rounded-lg text-sm sm:text-base"
    >
      {{ authStore.user?.role === 'manager' ? 'All Expenses' : 'My Expenses' }}
    </NuxtLink>
    <NuxtLink 
      to="/expenses/new" 
      :class="isActive('/expenses/new') ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100'"
      class="px-3 py-2 sm:px-4 rounded-lg text-sm sm:text-base"
    >
      Submit New
    </NuxtLink>
    <NuxtLink 
      v-if="authStore.user?.role === 'manager'" 
      to="/approvals" 
      :class="isActive('/approvals') ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100'"
      class="px-3 py-2 sm:px-4 rounded-lg text-sm sm:text-base"
    >
      Pending Approvals
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const route = useRoute()

const isActive = (path: string) => route.path === path
</script>
