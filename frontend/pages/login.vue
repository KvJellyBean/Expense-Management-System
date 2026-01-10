<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-500 to-blue-700 px-4">
    <div class="max-w-md w-full">
      <div class="card">
        <div class="text-center mb-6 sm:mb-8">
          <h1 class="text-2xl sm:text-3xl font-bold text-gray-800 mb-2">Expense Management</h1>
          <p class="text-sm sm:text-base text-gray-600">Sign in to your account</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4 sm:space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
            <input
              v-model="form.email"
              type="email"
              required
              class="input"
              placeholder="name@company.com"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Password</label>
            <input
              v-model="form.password"
              type="password"
              required
              class="input"
              placeholder="••••••••"
            />
          </div>

          <div v-if="error" class="p-3 bg-red-100 border border-red-400 text-red-700 rounded text-sm">
            {{ error }}
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full btn btn-primary"
          >
            {{ loading ? 'Processing...' : 'Sign In' }}
          </button>
        </form>

        <div class="mt-4 sm:mt-6 p-3 sm:p-4 bg-gray-50 rounded-lg">
          <p class="text-xs sm:text-sm text-gray-600 mb-2">Demo Accounts:</p>
          <div class="text-xs space-y-1">
            <p><strong>Employee:</strong> employee1@example.com | employee2@example.com </p>
            <p><strong>Manager:</strong> manager@example.com</p>
            <p class="text-gray-500">Password: password123</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const router = useRouter()

const form = ref({
  email: '',
  password: ''
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  try {
    loading.value = true
    error.value = ''
    
    await authStore.login(form.value.email, form.value.password)
    
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err.message || 'Login failed'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (authStore.isAuthenticated) {
    router.push('/dashboard')
  }
})
</script>
