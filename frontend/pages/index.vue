<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="text-center">
      <h1 class="text-4xl font-bold text-gray-800 mb-4">Expense Management System</h1>
      <p class="text-gray-600 mb-8">Selamat datang! Silakan login untuk melanjutkan.</p>
      <NuxtLink to="/login" class="btn btn-primary">
        Login
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()

onMounted(() => {
  if (process.client) {
    authStore.loadFromStorage()
    if (authStore.isAuthenticated) {
      navigateTo('/dashboard')
    }
  }
})
</script>
