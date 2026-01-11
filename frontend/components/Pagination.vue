<template>
  <div class="mt-4 flex flex-col sm:flex-row justify-center items-center gap-2 sm:gap-0 sm:space-x-2">
    <button @click="$emit('prev')" :disabled="page === 1" class="btn btn-secondary w-full sm:w-auto">Previous</button>
    <span class="px-4 py-2 text-sm sm:text-base">Page {{ page }} of {{ totalPages }}</span>
    <button @click="$emit('next')" :disabled="page >= totalPages" class="btn btn-secondary w-full sm:w-auto">Next</button>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ page: number; total: number; limit: number }>()
const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.limit)))
</script>
