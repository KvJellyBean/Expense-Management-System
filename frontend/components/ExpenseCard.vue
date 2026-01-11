<template>
  <div
    class="border rounded-lg p-3 sm:p-4 hover:shadow-md transition"
    :class="clickable ? 'cursor-pointer' : 'cursor-default'"
    role="button"
    tabindex="0"
    @click="onSelect"
    @keydown.enter="onSelect"
  >
    <div class="flex justify-between items-start gap-2">
      <div class="flex-1 min-w-0">
        <div class="flex flex-wrap items-center gap-1 sm:gap-2 mb-2">
          <span class="font-semibold text-base sm:text-lg break-all">{{ formatIDR(expense.amount_idr) }}</span>
          <span :class="getStatusClass(expense.status)" class="text-xs px-2 py-1 rounded whitespace-nowrap">
            {{ getStatusLabel(expense.status) }}
          </span>
          <span v-if="expense.auto_approved" class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded whitespace-nowrap">Auto</span>
        </div>
        <p class="text-gray-700 mb-1 text-sm sm:text-base line-clamp-2">{{ expense.description }}</p>
        <p class="text-xs text-gray-500">{{ formatDate(expense.submitted_at) }}</p>
      </div>

      <div class="flex-shrink-0">
        <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </div>
    </div>

    <div class="mt-3" v-if="$slots.default">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{ expense: any; clickable?: boolean }>(), { clickable: true })
const emit = defineEmits(['select'])

const { formatIDR, formatDate } = useFormat()

const getStatusClass = (status: string) => {
  const classes: any = {
    'awaiting_approval': 'bg-yellow-100 text-yellow-800',
    'completed': 'bg-green-100 text-green-800',
    'rejected': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getStatusLabel = (status: string) => {
  const labels: any = {
    'awaiting_approval': 'Pending',
    'completed': 'Approved',
    'rejected': 'Rejected'
  }
  return labels[status] || status
}

const onSelect = () => {
  if (props.clickable === false) return
  emit('select', props.expense.id)
}
</script>
