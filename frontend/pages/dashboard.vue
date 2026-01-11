<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader />

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
        <PageTabs />

        <!-- Expense List -->
        <div class="card mb-6">
          <h2 class="text-2xl font-bold mb-4">{{ authStore.user?.role === 'manager' ? 'All Expenses' : 'My Expenses' }}</h2>

          <div class="mb-4 flex flex-wrap gap-2">
            <button @click="setFilter('')" :class="filterStatus === '' && !filterAutoApproved ? 'btn btn-primary' : 'btn btn-secondary'">All</button>
            <button @click="setFilter('awaiting_approval')" :class="filterStatus === 'awaiting_approval' ? 'btn btn-primary' : 'btn btn-secondary'">Pending</button>
            <button @click="setFilter('completed')" :class="filterStatus === 'completed' && !filterAutoApproved ? 'btn btn-primary' : 'btn btn-secondary'">Approved</button>
            <button @click="setFilter('rejected')" :class="filterStatus === 'rejected' ? 'btn btn-primary' : 'btn btn-secondary'">Rejected</button>
            <button @click="setAutoApprovedFilter()" :class="filterAutoApproved ? 'btn btn-primary' : 'btn btn-secondary'">Auto-Approved</button>
          </div>

          <div v-if="loading" class="text-center py-8">Loading...</div>
          <div v-else-if="expenses.length === 0" class="text-center py-8 text-gray-500">No expenses{{ filterStatus ? ' with this status' : ' yet' }}</div>

          <div v-else class="space-y-4">
            <ExpenseCard v-for="expense in expenses" :key="expense.id" :expense="expense" @select="openExpenseDetail" />
          </div>

          <Pagination v-if="total > limit" :page="page" :total="total" :limit="limit" @prev="page > 1 && page--" @next="page < totalPages && page++" />
        </div>
    </div>

    <ExpenseDetailModal
      :expense="selectedExpense"
      :visible="!!selectedExpense"
      :processing="processingApproval"
      :approval-error="approvalError"
      @close="closeExpenseDetail"
      @approve="(notes) => selectedExpense && handleApproval(selectedExpense.id, 'approve', notes)"
      @reject="(notes) => selectedExpense && handleApproval(selectedExpense.id, 'reject', notes)"
    />
  </div>
</template>

<script setup lang="ts">
import ExpenseCard from '~/components/ExpenseCard.vue'
import ExpenseDetailModal from '~/components/ExpenseDetailModal.vue'
import Pagination from '~/components/Pagination.vue'

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
const { apiFetch } = useApi()

const expenses = ref<Expense[]>([])
const selectedExpense = ref<Expense | null>(null)
const total = ref(0)
const page = ref(1)
const limit = ref(20)
const loading = ref(false)
const filterStatus = ref('')
const filterAutoApproved = ref(false)
const processingApproval = ref(false)
const approvalError = ref('')

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit.value)))

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
  approvalError.value = ''
}

const handleApproval = async (expenseId: number, action: 'approve' | 'reject', notes?: string) => {
  try {
    processingApproval.value = true
    approvalError.value = ''

    const payload: any = {}
    if (notes && notes.trim()) {
      payload.notes = notes.trim()
    } 

    await apiFetch(`/expenses/${expenseId}/${action}`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    })

    await loadExpenses()
    closeExpenseDetail()
  } catch (err: any) {
    approvalError.value = err.message || `Failed to ${action} expense`
  } finally {
    processingApproval.value = false
  }
}

watch([page, filterStatus, filterAutoApproved], () => {
  loadExpenses()
})

onMounted(() => {
  loadExpenses()
})
</script>
