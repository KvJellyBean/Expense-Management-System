<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-bold text-gray-800">Expense Management</h1>
          </div>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-600">{{ authStore.user?.name }}</span>
            <span class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">{{ authStore.user?.role }}</span>
            <button @click="handleLogout" class="text-sm text-red-600 hover:text-red-800">
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
      <div class="flex flex-wrap gap-2 sm:gap-4 mb-4 sm:mb-6">
        <NuxtLink to="/dashboard" class="px-3 py-2 sm:px-4 rounded-lg bg-white text-gray-700 text-sm sm:text-base">
          My Expenses
        </NuxtLink>
        <NuxtLink to="/expenses/new" class="px-3 py-2 sm:px-4 rounded-lg bg-white text-gray-700 text-sm sm:text-base">
          Submit New
        </NuxtLink>
        <NuxtLink to="/approvals" class="px-3 py-2 sm:px-4 rounded-lg bg-blue-600 text-white text-sm sm:text-base">
          Pending Approvals
        </NuxtLink>
      </div>

      <div class="card">
        <h2 class="text-2xl font-bold mb-6">Pending Approvals</h2>
        
        <div v-if="loading" class="text-center py-8">Loading...</div>
        <div v-else-if="expenses.length === 0" class="text-center py-8 text-gray-500">
          No pending approvals
        </div>
        <div v-else class="space-y-4 sm:space-y-6">
          <div v-for="expense in expenses" :key="expense.id" class="border rounded-lg p-4 sm:p-6 hover:shadow-md transition">
            <div class="mb-4">
              <div>
                <div class="flex flex-wrap items-center gap-2 mb-2">
                  <span class="font-bold text-xl sm:text-2xl text-blue-600 break-all">{{ formatIDR(expense.amount_idr) }}</span>
                  <span class="text-xs bg-yellow-100 text-yellow-800 px-2 py-1 rounded whitespace-nowrap">
                    Pending Approval
                  </span>
                </div>
                <p class="text-gray-700 mb-2 text-sm sm:text-base">{{ expense.description }}</p>
                <p class="text-xs sm:text-sm text-gray-500">Submitted: {{ formatDate(expense.submitted_at) }}</p>
                
                <!-- Receipt Preview/Link -->
                <div v-if="expense.receipt_url" class="mt-2">
                  <p class="text-xs text-gray-600 mb-1">Receipt:</p>
                  <img 
                    :src="expense.receipt_url" 
                    alt="Receipt thumbnail"
                    class="max-w-xs max-h-24 rounded border bg-gray-50"
                    @error="(e) => (e.target as HTMLImageElement).src = 'https://placehold.co/200x150/eee/666?text=Receipt'"
                  />
                  <a :href="expense.receipt_url" target="_blank" class="block text-xs text-blue-600 hover:underline mt-1">
                    🔍 View Full Receipt
                  </a>
                </div>
              </div>
            </div>

            <div class="border-t pt-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">Notes (Optional)</label>
              <textarea
                v-model="notes[expense.id]"
                rows="2"
                class="input mb-3"
                placeholder="Add approval notes..."
              ></textarea>
              
              <div class="flex flex-col sm:flex-row gap-2 sm:gap-3">
                <button
                  @click="handleApprove(expense.id)"
                  :disabled="processing[expense.id]"
                  class="btn btn-success w-full sm:w-auto"
                >
                  {{ processing[expense.id] ? 'Processing...' : 'Approve' }}
                </button>
                <button
                  @click="handleReject(expense.id)"
                  :disabled="processing[expense.id]"
                  class="btn btn-danger w-full sm:w-auto"
                >
                  {{ processing[expense.id] ? 'Processing...' : 'Reject' }}
                </button>
              </div>
            </div>

            <div v-if="errors[expense.id]" class="mt-3 p-3 bg-red-100 border border-red-400 text-red-700 rounded text-sm">
              {{ errors[expense.id] }}
            </div>
            <div v-if="successes[expense.id]" class="mt-3 p-3 bg-green-100 border border-green-400 text-green-700 rounded text-sm">
              {{ successes[expense.id] }}
            </div>
          </div>
        </div>

        <div v-if="total > limit" class="mt-6 flex flex-col sm:flex-row justify-center items-center gap-2 sm:gap-0 sm:space-x-2">
          <button @click="page > 1 && page--" :disabled="page === 1" class="btn btn-secondary w-full sm:w-auto">Previous</button>
          <span class="px-4 py-2 text-sm sm:text-base">Page {{ page }} of {{ Math.ceil(total / limit) }}</span>
          <button @click="page < Math.ceil(total / limit) && page++" :disabled="page >= Math.ceil(total / limit)" class="btn btn-secondary w-full sm:w-auto">Next</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Expense {
  id: number
  amount_idr: number
  description: string
  status: string
  submitted_at: string
  receipt_url?: string
}

definePageMeta({
  middleware: ['auth', 'manager']
})

const authStore = useAuthStore()
const router = useRouter()
const { apiFetch } = useApi()
const { formatIDR, formatDate } = useFormat()

const expenses = ref<Expense[]>([])
const total = ref(0)
const page = ref(1)
const limit = ref(20)
const loading = ref(false)

const notes = ref<Record<number, string>>({})
const processing = ref<Record<number, boolean>>({})
const errors = ref<Record<number, string>>({})
const successes = ref<Record<number, string>>({})

const loadPendingApprovals = async () => {
  try {
    loading.value = true
    const query = new URLSearchParams({
      page: page.value.toString(),
      limit: limit.value.toString()
    })

    const data = await apiFetch(`/expenses/pending?${query.toString()}`)
    expenses.value = data.expenses || []
    total.value = data.total || 0
  } catch (err) {
    console.error('Failed to load pending approvals:', err)
  } finally {
    loading.value = false
  }
}

const handleApprove = async (expenseId: number) => {
  try {
    processing.value[expenseId] = true
    errors.value[expenseId] = ''
    successes.value[expenseId] = ''

    const payload: any = {}
    if (notes.value[expenseId]) {
      payload.notes = notes.value[expenseId]
    }

    await apiFetch(`/expenses/${expenseId}/approve`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    })

    successes.value[expenseId] = 'Expense approved successfully!'
    
    setTimeout(() => {
      loadPendingApprovals()
    }, 1000)
  } catch (err: any) {
    errors.value[expenseId] = err.message || 'Failed to approve expense'
  } finally {
    processing.value[expenseId] = false
  }
}

const handleReject = async (expenseId: number) => {
  try {
    processing.value[expenseId] = true
    errors.value[expenseId] = ''
    successes.value[expenseId] = ''

    const payload: any = {}
    if (notes.value[expenseId]) {
      payload.notes = notes.value[expenseId]
    }

    await apiFetch(`/expenses/${expenseId}/reject`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    })

    successes.value[expenseId] = 'Expense rejected successfully!'
    
    setTimeout(() => {
      loadPendingApprovals()
    }, 1000)
  } catch (err: any) {
    errors.value[expenseId] = err.message || 'Failed to reject expense'
  } finally {
    processing.value[expenseId] = false
  }
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

watch(page, () => {
  loadPendingApprovals()
})

onMounted(() => {
  loadPendingApprovals()
})
</script>
