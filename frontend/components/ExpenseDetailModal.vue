<template>
  <div v-if="visible && expense" @click.self="emitClose" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-lg max-w-2xl w-full p-4 sm:p-6 max-h-[80vh] overflow-y-auto">
      <div class="flex justify-between items-start mb-4">
        <h3 class="text-xl sm:text-2xl font-bold">Expense Detail</h3>
        <button @click="emitClose" class="text-gray-400 hover:text-gray-600 p-1">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-500">Amount</label>
          <p class="text-xl sm:text-2xl font-bold text-blue-600 break-all">{{ formatIDR(expense.amount_idr) }}</p>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-500">Status</label>
          <div class="mt-1">
            <span :class="getStatusClass(expense.status)" class="text-sm px-3 py-1 rounded">{{ getStatusLabel(expense.status) }}</span>
            <span v-if="expense.auto_approved" class="ml-2 text-sm bg-blue-100 text-blue-800 px-3 py-1 rounded">Auto</span>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-500">Description</label>
          <p class="mt-1 text-gray-700 break-words whitespace-pre-wrap">{{ expense.description }}</p>
        </div>

        <div v-if="expense.receipt_url">
          <label class="block text-sm font-medium text-gray-500 mb-2">Receipt</label>
          <div class="space-y-2">
            <div class="border rounded-lg overflow-hidden bg-gray-50">
              <img
                :src="expense.receipt_url"
                :alt="'Receipt for ' + expense.description"
                class="max-w-full sm:max-w-md max-h-96 mx-auto"
                @error="onImageError"
              />
            </div>
            <a :href="expense.receipt_url" target="_blank" class="inline-block text-sm text-blue-600 hover:underline">🔗 Open Receipt in New Tab</a>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-500">Submitted At</label>
          <p class="mt-1 text-gray-700">{{ formatDate(expense.submitted_at) }}</p>
        </div>

        <div v-if="expense.approval && expense.approval.notes" class="border-t pt-4">
          <label class="block text-sm font-medium text-gray-500 mb-2">{{ expense.status === 'rejected' ? 'Rejection Notes' : 'Approval Notes' }}</label>
          <div :class="expense.status === 'rejected' ? 'bg-red-50 border border-red-200' : 'bg-gray-50'" class="p-3 rounded">
            <p :class="expense.status === 'rejected' ? 'text-red-700' : 'text-gray-700'">{{ expense.approval.notes }}</p>
          </div>
        </div>

        <!-- Manager Approval Actions -->
        <div v-if="isManager && expense.status === 'awaiting_approval'" class="border-t pt-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Manager Action</label>
          <div class="space-y-3">
            <textarea v-model="notes" placeholder="Add notes (optional)" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" rows="2"></textarea>
            <div class="flex flex-col sm:flex-row gap-2">
              <button @click="emitApprove" :disabled="processing" class="flex-1 bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium">{{ processing ? 'Processing...' : '✓ Approve' }}</button>
              <button @click="emitReject" :disabled="processing" class="flex-1 bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium">{{ processing ? 'Processing...' : '✗ Reject' }}</button>
            </div>
            <p v-if="approvalError" class="text-red-600 text-sm">{{ approvalError }}</p>
          </div>
        </div>
      </div>

      <div class="mt-6 flex justify-end">
        <button @click="emitClose" class="btn btn-secondary w-full sm:w-auto">Close</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ expense: any | null; visible: boolean; processing?: boolean; approvalError?: string }>()
const emit = defineEmits(['close', 'approve', 'reject'])

const { formatIDR, formatDate } = useFormat()
const authStore = useAuthStore()

const notes = ref('')

const isManager = computed(() => authStore.user?.role === 'manager')
const processing = props.processing || false
const approvalError = props.approvalError || ''

watch(() => props.expense, () => {
  notes.value = ''
})

watch(() => props.visible, (v) => {
  if (!v) notes.value = ''
})

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

const emitClose = () => emit('close')
const emitApprove = () => emit('approve', notes.value.trim())
const emitReject = () => emit('reject', notes.value.trim())

const onImageError = (e: Event) => {
  ;(e.target as HTMLImageElement).src = 'https://placehold.co/400x500/eee/666?text=Receipt+Unavailable'
}
</script>
