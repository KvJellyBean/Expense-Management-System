import { defineNuxtPlugin } from "#app"

export default defineNuxtPlugin(() => {
  const authStore = useAuthStore()
  
  // Load auth state from localStorage on app start
  authStore.loadFromStorage()
})
