// @ts-ignore - Nuxt auto-imports
export const useApi = () => {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const apiFetch = async (endpoint: string, options: any = {}) => {
    const headers: any = {
      'Content-Type': 'application/json',
      ...options.headers
    }

    if (authStore.token) {
      headers.Authorization = `Bearer ${authStore.token}`
    }

    const response = await fetch(`${config.public.apiBase}${endpoint}`, {
      ...options,
      headers
    })

    if (response.status === 401) {
      authStore.logout()
      navigateTo('/login')
      throw new Error('Unauthorized')
    }

    if (!response.ok) {
      const error = await response.text()
      throw new Error(error || 'Request failed')
    }

    return response.json()
  }

  return { apiFetch }
}
