import type { Ref, ComputedRef } from 'vue'

// Global type declarations for Nuxt auto-imports
declare global {
  // Vue composables
  const ref: <T = any>(value?: T) => Ref<T>
  const computed: <T = any>(getter: () => T) => ComputedRef<T>
  const watch: typeof import('vue')['watch']
  const onMounted: (fn: () => void) => void
  const onUnmounted: (fn: () => void) => void
  const reactive: <T extends object>(target: T) => T
  const nextTick: () => Promise<void>
  
  // Vue Router
  const useRouter: () => import('vue-router').Router
  const useRoute: () => import('vue-router').RouteLocationNormalizedLoaded
  
  // Nuxt
  const useRuntimeConfig: () => {
    public: {
      apiBase: string
    }
  }
  const navigateTo: (to: string | object) => Promise<void>
  const definePageMeta: (meta: {
    middleware?: string | string[]
    layout?: string
  }) => void
  
  // Custom composables - declare as variable instead of const to avoid redeclaration
  var useAuthStore: () => any
  var useApi: () => {
    apiFetch: <T = any>(url: string, options?: any) => Promise<T>
  }
  var useFormat: () => {
    formatIDR: (amount: number) => string
    parseIDR: (formatted: string) => number
    formatDate: (date: string) => string
  }
}

export {}
