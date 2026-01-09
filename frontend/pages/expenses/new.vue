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
        <NuxtLink to="/dashboard" class="px-3 py-2 sm:px-4 rounded-lg bg-white text-gray-700 hover:bg-gray-100 text-sm sm:text-base">
          My Expenses
        </NuxtLink>
        <NuxtLink to="/expenses/new" class="px-3 py-2 sm:px-4 rounded-lg bg-blue-600 text-white text-sm sm:text-base">
          Submit New
        </NuxtLink>
        <NuxtLink v-if="authStore.user?.role === 'manager'" to="/approvals" class="px-3 py-2 sm:px-4 rounded-lg bg-white text-gray-700 hover:bg-gray-100 text-sm sm:text-base">
          Pending Approvals
        </NuxtLink>
      </div>

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
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Receipt Upload (Optional)
            </label>
            <input
              @change="handleFileUpload"
              type="file"
              accept="image/*,.pdf"
              class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
            />
            <p class="text-xs text-gray-500 mt-1">Accepted: Images (JPG, PNG) or PDF, max 5MB</p>
            <div v-if="form.receiptPreview" class="mt-3">
              <p class="text-xs text-gray-600 mb-2">File uploaded:</p>
              <div v-if="form.receiptPreview.startsWith('data:image')" class="space-y-2">
                <img :src="form.receiptPreview" class="max-w-full sm:max-w-xs max-h-48 rounded border" alt="Receipt preview" />
                <p class="text-xs text-gray-500">✓ {{ form.receiptFileName }}</p>
                <p class="text-xs text-blue-600">Mock URL: {{ form.receiptUrl }}</p>
              </div>
              <div v-else-if="form.receiptPreview === 'pdf'" class="space-y-2">
                <div class="flex items-center gap-2 p-3 bg-gray-50 rounded border">
                  <svg class="w-8 h-8 text-red-600" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M4 18h12V6h-4V2H4v16zm-2 1V0h10l4 4v16H2v-1z"/>
                    <text x="6" y="14" font-size="6" fill="currentColor">PDF</text>
                  </svg>
                  <div class="flex-1">
                    <p class="text-sm font-medium text-gray-700">{{ form.receiptFileName }}</p>
                    <p class="text-xs text-gray-500">✓ File ready to submit</p>
                  </div>
                </div>
              </div>
            </div>
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
definePageMeta({
  middleware: ['auth']
})

const authStore = useAuthStore()
const router = useRouter()
const { apiFetch } = useApi()
const { formatIDR, parseIDR } = useFormat()

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

const handleFileUpload = (e: Event) => {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  
  if (!file) {
    form.value.receiptFile = null
    form.value.receiptPreview = ''
    form.value.receiptFileName = ''
    form.value.receiptUrl = ''
    return
  }

  // Check file size (max 5MB)
  if (file.size > 5 * 1024 * 1024) {
    error.value = 'File size must be less than 5MB'
    input.value = ''
    return
  }

  form.value.receiptFile = file
  form.value.receiptFileName = file.name

  // MOCK file upload: Use static mock file (Task.txt: "mock the receipt upload")
  // In production: Upload to S3/CDN and get real URL
  // Using static mock file instead of fake URL to avoid "can't be reached" errors
  const fileExtension = file.name.split('.').pop()?.toLowerCase()
  const mockUrl = '/mock-receipt.jpg'
  form.value.receiptUrl = mockUrl

  // Create preview using FileReader for UI display only
  const reader = new FileReader()
  reader.onload = (e) => {
    const dataUrl = e.target?.result as string
    
    // Set preview based on file type
    if (file.type === 'application/pdf') {
      form.value.receiptPreview = 'pdf'
    } else {
      form.value.receiptPreview = dataUrl
    }
  }
  
  // Read file as data URL (base64) for preview only
  reader.readAsDataURL(file)
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
  // Clear file input
  const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement
  if (fileInput) fileInput.value = ''
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
