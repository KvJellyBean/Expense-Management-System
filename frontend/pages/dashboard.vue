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
        <NuxtLink to="/dashboard" class="px-3 py-2 sm:px-4 rounded-lg bg-blue-600 text-white text-sm sm:text-base">
          {{ authStore.user?.role === 'manager' ? 'All Expenses' : 'My Expenses' }}
        </NuxtLink>
        <NuxtLink to="/expenses/new" class="px-3 py-2 sm:px-4 rounded-lg bg-white text-gray-700 hover:bg-gray-100 text-sm sm:text-base">
          Submit New
        </NuxtLink>
        <NuxtLink v-if="authStore.user?.role === 'manager'" to="/approvals" class="px-3 py-2 sm:px-4 rounded-lg bg-white text-gray-700 hover:bg-gray-100 text-sm sm:text-base">
          Pending Approvals
        </NuxtLink>
      </div>

      <!-- Expense List -->
      <div class="card mb-6">
        <h2 class="text-2xl font-bold mb-4">{{ authStore.user?.role === 'manager' ? 'All Expenses' : 'My Expenses' }}</h2>
        
        <div class="mb-4 flex flex-wrap gap-2">
          <button @click="setFilter('')" :class="filterStatus === '' && !filterAutoApproved ? 'btn btn-primary' : 'btn btn-secondary'">
            All
          </button>
          <button @click="setFilter('awaiting_approval')" :class="filterStatus === 'awaiting_approval' ? 'btn btn-primary' : 'btn btn-secondary'">
            Pending
          </button>
          <button @click="setFilter('completed')" :class="filterStatus === 'completed' && !filterAutoApproved ? 'btn btn-primary' : 'btn btn-secondary'">
            Approved
          </button>
          <button @click="setFilter('rejected')" :class="filterStatus === 'rejected' ? 'btn btn-primary' : 'btn btn-secondary'">
            Rejected
          </button>
          <button @click="setAutoApprovedFilter()" :class="filterAutoApproved ? 'btn btn-primary' : 'btn btn-secondary'">
            Auto-Approved
          </button>
        </div>

        <div v-if="loading" class="text-center py-8">Loading...</div>
        <div v-else-if="expenses.length === 0" class="text-center py-8 text-gray-500">
          No expenses{{ filterStatus ? ' with this status' : ' yet' }}
        </div>
        <div v-else class="space-y-4">
          <div 
            v-for="expense in expenses" 
            :key="expense.id" 
            @click="openExpenseDetail(expense.id)"
            class="border rounded-lg p-3 sm:p-4 hover:shadow-md transition cursor-pointer"
          >
            <div class="flex justify-between items-start gap-2">
              <div class="flex-1 min-w-0">
                <div class="flex flex-wrap items-center gap-1 sm:gap-2 mb-2">
                  <span class="font-semibold text-base sm:text-lg break-all">{{ formatIDR(expense.amount_idr) }}</span>
                  <span :class="getStatusClass(expense.status)" class="text-xs px-2 py-1 rounded whitespace-nowrap">
                    {{ getStatusLabel(expense.status) }}
                  </span>
                  <span v-if="expense.auto_approved" class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded whitespace-nowrap">
                    Auto
                  </span>
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
          </div>
        </div>

        <div v-if="total > limit" class="mt-4 flex flex-col sm:flex-row justify-center items-center gap-2 sm:gap-0 sm:space-x-2">
          <button @click="page > 1 && page--" :disabled="page === 1" class="btn btn-secondary w-full sm:w-auto">Previous</button>
          <span class="px-4 py-2 text-sm sm:text-base">Page {{ page }} of {{ Math.ceil(total / limit) }}</span>
          <button @click="page < Math.ceil(total / limit) && page++" :disabled="page >= Math.ceil(total / limit)" class="btn btn-secondary w-full sm:w-auto">Next</button>
        </div>
      </div>
    </div>

    <!-- Expense Detail Modal -->
    <div v-if="selectedExpense" @click.self="closeExpenseDetail" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-lg max-w-2xl w-full p-4 sm:p-6 max-h-screen overflow-y-auto">
        <div class="flex justify-between items-start mb-4">
          <h3 class="text-xl sm:text-2xl font-bold">Expense Detail</h3>
          <button @click="closeExpenseDetail" class="text-gray-400 hover:text-gray-600 p-1">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-500">Amount</label>
            <p class="text-xl sm:text-2xl font-bold text-blue-600 break-all">{{ formatIDR(selectedExpense.amount_idr) }}</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-500">Status</label>
            <div class="mt-1">
              <span :class="getStatusClass(selectedExpense.status)" class="text-sm px-3 py-1 rounded">
                {{ getStatusLabel(selectedExpense.status) }}
              </span>
              <span v-if="selectedExpense.auto_approved" class="ml-2 text-sm bg-blue-100 text-blue-800 px-3 py-1 rounded">
                Auto
              </span>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-500">Description</label>
            <p class="mt-1 text-gray-700">{{ selectedExpense.description }}</p>
          </div>

          <div v-if="selectedExpense.receipt_url">
            <label class="block text-sm font-medium text-gray-500 mb-2">Receipt</label>
            <!-- Mock Receipt Preview (Task.txt: "mock the receipt upload") -->
            <div class="space-y-2">
              <div class="border rounded-lg overflow-hidden bg-gray-50">
                <img 
                  :src="selectedExpense.receipt_url" 
                  :alt="'Receipt for ' + selectedExpense.description"
                  class="max-w-full sm:max-w-md max-h-96 mx-auto"
                  @error="(e) => (e.target as HTMLImageElement).src = 'https://placehold.co/400x500/eee/666?text=Receipt+Unavailable'"
                />
              </div>
              <a 
                :href="selectedExpense.receipt_url" 
                target="_blank" 
                class="inline-block text-sm text-blue-600 hover:underline"
              >
                🔗 Open Receipt in New Tab
              </a>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-500">Submitted At</label>
            <p class="mt-1 text-gray-700">{{ formatDate(selectedExpense.submitted_at) }}</p>
          </div>

          <div v-if="selectedExpense.approval && selectedExpense.approval.notes" class="border-t pt-4">
            <label class="block text-sm font-medium text-gray-500 mb-2">
              {{ selectedExpense.status === 'rejected' ? 'Rejection Notes' : 'Approval Notes' }}
            </label>
            <div :class="selectedExpense.status === 'rejected' ? 'bg-red-50 border border-red-200' : 'bg-gray-50'" class="p-3 rounded">
              <p :class="selectedExpense.status === 'rejected' ? 'text-red-700' : 'text-gray-700'">{{ selectedExpense.approval.notes }}</p>
            </div>
          </div>

          <!-- Manager Approval Actions (if expense is pending and user is manager) -->
          <div v-if="authStore.user?.role === 'manager' && selectedExpense.status === 'awaiting_approval'" class="border-t pt-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Manager Action</label>
            <div class="space-y-3">
              <textarea
                v-model="approvalNotes"
                placeholder="Add notes (optional)"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                rows="2"
              ></textarea>
              <div class="flex flex-col sm:flex-row gap-2">
                <button
                  @click="handleApproval(selectedExpense.id, 'approve')"
                  :disabled="processingApproval"
                  class="flex-1 bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium"
                >
                  {{ processingApproval ? 'Processing...' : '✓ Approve' }}
                </button>
                <button
                  @click="handleApproval(selectedExpense.id, 'reject')"
                  :disabled="processingApproval"
                  class="flex-1 bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium"
                >
                  {{ processingApproval ? 'Processing...' : '✗ Reject' }}
                </button>
              </div>
              <p v-if="approvalError" class="text-red-600 text-sm">{{ approvalError }}</p>
            </div>
          </div>
        </div>

        <div class="mt-6 flex justify-end">
          <button @click="closeExpenseDetail" class="btn btn-secondary w-full sm:w-auto">
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Approval {
  id: number
  expense_id: number
  approver_id: number
  status: string
  notes?: string
  created_at: string
}

interface Expense {
  id: number
  amount_idr: number
  description: string
  status: string
  auto_approved: boolean
  submitted_at: string
  receipt_url?: string
  approval?: Approval
}

definePageMeta({
  middleware: ['auth']
})

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const { apiFetch } = useApi()
const { formatIDR, parseIDR, formatDate } = useFormat()

const expenses = ref<Expense[]>([])
const selectedExpense = ref<Expense | null>(null)
const total = ref(0)
const page = ref(1)
const limit = ref(20)
const loading = ref(false)
const filterStatus = ref('')
const filterAutoApproved = ref(false)
const approvalNotes = ref('')
const processingApproval = ref(false)
const approvalError = ref('')

const form = ref({
  amount: 0,
  amountFormatted: '',
  description: '',
  receiptUrl: ''
})

const submitting = ref(false)
const error = ref('')
const success = ref(false)

const formatAmount = (event: any) => {
  const value = parseIDR(event.target.value)
  form.value.amount = value
  form.value.amountFormatted = value ? formatIDR(value).replace('IDR', 'Rp').trim() : ''
}

const handleSubmit = async () => {
  try {
    submitting.value = true
    error.value = ''
    success.value = false

    const payload: any = {
      amount_idr: form.value.amount,
      description: form.value.description
    }

    if (form.value.receiptUrl) {
      payload.receipt_url = form.value.receiptUrl
    }

    await apiFetch('/expenses', {
      method: 'POST',
      body: JSON.stringify(payload)
    })

    success.value = true
    resetForm()
    
    setTimeout(() => {
      router.push('/dashboard')
    }, 1500)
  } catch (err: any) {
    error.value = err.message || 'Failed to submit expense'
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  form.value = {
    amount: 0,
    amountFormatted: '',
    description: '',
    receiptUrl: ''
  }
  error.value = ''
  success.value = false
}

const loadExpenses = async () => {
  try {
    loading.value = true
    const query = new URLSearchParams({
      page: page.value.toString(),
      limit: limit.value.toString()
    })
    
    if (filterStatus.value) {
      query.append('status', filterStatus.value)
    }

    const data = await apiFetch(`/expenses?${query.toString()}`)
    let allExpenses = data.expenses || []
    
    // Client-side filter for auto-approved
    if (filterAutoApproved.value) {
      allExpenses = allExpenses.filter((exp: Expense) => exp.auto_approved === true)
    }
    
    expenses.value = allExpenses
    total.value = filterAutoApproved.value ? allExpenses.length : (data.total || 0)
  } catch (err) {
    console.error('Failed to load expenses:', err)
  } finally {
    loading.value = false
  }
}

const setFilter = (status: string) => {
  filterStatus.value = status
  filterAutoApproved.value = false
  page.value = 1
}

const setAutoApprovedFilter = () => {
  filterStatus.value = 'completed'
  filterAutoApproved.value = true
  page.value = 1
}

const openExpenseDetail = async (expenseId: number) => {
  try {
    const expense = await apiFetch(`/expenses/${expenseId}`)
    selectedExpense.value = expense
  } catch (err) {
    console.error('Failed to load expense detail:', err)
  }
}

const closeExpenseDetail = () => {
  selectedExpense.value = null
  approvalNotes.value = ''
  approvalError.value = ''
}

const handleApproval = async (expenseId: number, action: 'approve' | 'reject') => {
  try {
    processingApproval.value = true
    approvalError.value = ''

    const payload: any = {}
    if (approvalNotes.value.trim()) {
      payload.notes = approvalNotes.value.trim()
    }

    await apiFetch(`/expenses/${expenseId}/${action}`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    })

    // Refresh expense list and close modal
    await loadExpenses()
    closeExpenseDetail()
  } catch (err: any) {
    approvalError.value = err.message || `Failed to ${action} expense`
  } finally {
    processingApproval.value = false
  }
}

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

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

watch([page, filterStatus, filterAutoApproved], () => {
  loadExpenses()
})

onMounted(() => {
  loadExpenses()
})
</script>
