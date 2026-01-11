<template>
  <header class="bg-white shadow-sm">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <h1 class="text-xl font-bold text-gray-800">Expense Management</h1>
        </div>
        <nav class="flex items-center space-x-4" aria-label="User menu">
          <span class="text-sm text-gray-600">{{ authStore.user?.name }}</span>
          <span class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">{{ authStore.user?.role }}</span>
          <button @click="handleLogout" class="text-sm text-red-600 hover:text-red-800">
            Logout
          </button>
        </nav>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
