<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader />

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
      <PageTabs />

      <div class="card">
        <h2 class="text-2xl font-bold mb-6">Pending Approvals</h2>
        
        <div v-if="loading" class="text-center py-8">Loading...</div>
        <div v-else-if="expenses.length === 0" class="text-center py-8 text-gray-500">
          No pending approvals
        </div>
        <div v-else class="space-y-4 sm:space-y-6">
          <div v-for="expense in expenses" :key="expense.id">
            <ExpenseCard :expense="expense" :clickable="false">
              <ApprovalActions
                :notes="notes[expense.id]"
                :processing="processing[expense.id]"
                @update:notes="val => notes[expense.id] = val"
                @approve="() => handleApprove(expense.id)"
                @reject="() => handleReject(expense.id)"
              />
              <div v-if="errors[expense.id]" class="mt-3 p-3 bg-red-100 border border-red-400 text-red-700 rounded text-sm">{{ errors[expense.id] }}</div>
              <div v-if="successes[expense.id]" class="mt-3 p-3 bg-green-100 border border-green-400 text-green-700 rounded text-sm">{{ successes[expense.id] }}</div>
            </ExpenseCard>
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
import AppHeader from '~/components/AppHeader.vue'
import PageTabs from '~/components/PageTabs.vue'
import ExpenseCard from '~/components/ExpenseCard.vue'
import ApprovalActions from '~/components/ApprovalActions.vue'

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

const { apiFetch } = useApi()
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

watch(page, () => {
  loadPendingApprovals()
})

onMounted(() => {
  loadPendingApprovals()
})
</script>
