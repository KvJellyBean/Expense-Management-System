// Nuxt 3 auto-imports
declare module '#app' {
  export const useRuntimeConfig: () => any
  export const definePageMeta: (meta: any) => void
  export const navigateTo: (to: any) => void
}

declare module 'vue' {
  export const ref: <T>(value: T) => { value: T }
  export const computed: <T>(getter: () => T) => { value: T }
  export const watch: (...args: any[]) => void
  export const onMounted: (fn: () => void) => void
  export const onUnmounted: (fn: () => void) => void
  export const reactive: <T extends object>(target: T) => T
  export const nextTick: () => Promise<void>
}

declare module 'vue-router' {
  export const useRouter: () => any
  export const useRoute: () => any
}

declare module '~/stores/auth' {
  export const useAuthStore: () => any
}

declare module '~/composables/useApi' {
  export const useApi: () => any
}

declare module '~/composables/useFormat' {
  export const useFormat: () => any
}

export {}
