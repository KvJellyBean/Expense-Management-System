export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()
  
  if (typeof window !== 'undefined') {
    authStore.loadFromStorage()
  }

  if (!authStore.isManager) {
    return navigateTo('/dashboard')
  }
})
