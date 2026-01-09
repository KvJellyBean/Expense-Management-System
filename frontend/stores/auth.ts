// @ts-ignore - Nuxt auto-imports defineStore
interface User {
  id: number
  email: string
  name: string
  role: string
}

interface AuthState {
  user: User | null
  token: string | null
}

// @ts-ignore
export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null
  }),

  getters: {
    isAuthenticated: (state: AuthState) => !!state.token,
    isManager: (state: AuthState) => state.user?.role === 'manager'
  },

  actions: {
    async login(email: string, password: string) {
      const config = useRuntimeConfig()
      
      const response = await fetch(`${config.public.apiBase}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, password })
      })

      if (!response.ok) {
        const error = await response.text()
        throw new Error(error || 'Login failed')
      }

      const data = await response.json()
      
      this.token = data.token
      this.user = data.user

      if (typeof window !== 'undefined') {
        localStorage.setItem('token', data.token)
        localStorage.setItem('user', JSON.stringify(data.user))
      }
    },

    logout() {
      this.token = null
      this.user = null
      
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
      }
    },

    loadFromStorage() {
      if (typeof window !== 'undefined') {
        const token = localStorage.getItem('token')
        const user = localStorage.getItem('user')
        
        if (token && user) {
          this.token = token
          this.user = JSON.parse(user)
        }
      }
    }
  }
})
