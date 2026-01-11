<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader />

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
      <PageTabs />

      <div class="card max-w-2xl mx-auto">
        <h2 class="text-2xl font-bold mb-6">Submit New Expense</h2>
        
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Amount (IDR) <span class="text-red-500">*</span>
            </label>
            <input
              v-model="form.amountFormatted"
              @input="formatAmount"
              type="text"
              required
              class="input"
              placeholder="Rp 1.000.000"
            />
            <p class="text-xs text-gray-500 mt-1">Minimum: Rp 10.000 | Maximum: Rp 50.000.000</p>
            <p v-if="form.amount >= 1000000" class="text-sm text-orange-600 mt-2">
              This expense requires manager approval (≥ Rp 1.000.000)
            </p>
            <p v-else-if="form.amount > 0" class="text-sm text-green-600 mt-2">
              ✓ This expense will be auto-approved (< Rp 1.000.000)
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Description <span class="text-red-500">*</span>
            </label>
            <textarea
              v-model="form.description"
              required
              rows="3"
              class="input"
              placeholder="Describe the purpose of this expense..."
            ></textarea>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Receipt Upload (Optional)</label>
            <FileUpload ref="fileUploadRef" accept="image/*,.pdf" :maxSize="5 * 1024 * 1024" hint="Accepted: Images (JPG, PNG) or PDF, max 5MB" @file-selected="onFileSelected" @error="(msg) => error = msg" />
          </div>

          <div v-if="error" class="p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {{ error }}
          </div>

          <div v-if="success" class="p-3 bg-green-100 border border-green-400 text-green-700 rounded">
            Expense submitted successfully! Redirecting...
          </div>

          <div class="flex flex-col sm:flex-row gap-3 sm:gap-4">
            <button type="submit" :disabled="submitting" class="btn btn-primary w-full sm:w-auto">
              {{ submitting ? 'Processing...' : 'Submit Expense' }}
            </button>
            <button type="button" @click="resetForm" class="btn btn-secondary w-full sm:w-auto">
              Reset
            </button>
            <NuxtLink to="/dashboard" class="btn btn-secondary w-full sm:w-auto text-center">
              Cancel
            </NuxtLink>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import FileUpload from '~/components/FileUpload.vue'

definePageMeta({
  middleware: ['auth']
})

const router = useRouter()
const { apiFetch } = useApi()
const { formatIDR, parseIDR } = useFormat()
const fileUploadRef = ref<any>(null)

const form = ref({
  amount: 0,
  amountFormatted: '',
  description: '',
  receiptUrl: '',
  receiptFile: null as File | null,
  receiptPreview: '',
  receiptFileName: ''
})

const submitting = ref(false)
const error = ref('')
const success = ref(false)

const formatAmount = (e: Event) => {
  const input = e.target as HTMLInputElement
  const value = parseIDR(input.value)
  form.value.amount = value
  form.value.amountFormatted = formatIDR(value).replace('IDR', 'Rp')
}

const onFileSelected = (payload: any) => {
  const { file, preview, fileName } = payload || {}

  error.value = ''
  if (!file) {
    form.value.receiptFile = null
    form.value.receiptPreview = ''
    form.value.receiptFileName = ''
    form.value.receiptUrl = ''
    return
  }

  form.value.receiptFile = file
  form.value.receiptFileName = fileName || file.name
  form.value.receiptUrl = '/mock-receipt.jpg'
  form.value.receiptPreview = preview || ''
}

const handleSubmit = async () => {
  try {
    error.value = ''
    submitting.value = true

    if (form.value.amount < 10000) {
      error.value = 'Minimum amount adalah Rp 10.000'
      return
    }

    if (form.value.amount > 50000000) {
      error.value = 'Maximum amount adalah Rp 50.000.000'
      return
    }

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
    receiptUrl: '',
    receiptFile: null,
    receiptPreview: '',
    receiptFileName: ''
  }
  error.value = ''
  success.value = false
  fileUploadRef.value?.clear()
}
</script>
