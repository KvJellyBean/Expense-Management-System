<template>
  <div>
    <input
      ref="fileInput"
      @change="onChange"
      type="file"
      :accept="accept"
      class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
    />
    <p class="text-xs text-gray-500 mt-1">{{ hint }}</p>

    <div v-if="preview" class="mt-3">
      <div v-if="isImage" class="space-y-2">
        <img :src="preview" class="max-w-full sm:max-w-xs max-h-48 rounded border" alt="Receipt preview" />
        <p class="text-xs text-gray-500">✓ {{ fileName }}</p>
      </div>
      <div v-else-if="isPdf" class="space-y-2">
        <div class="flex items-center gap-2 p-3 bg-gray-50 rounded border">
          <svg class="w-8 h-8 text-red-600" fill="currentColor" viewBox="0 0 20 20">
            <path d="M4 18h12V6h-4V2H4v16zm-2 1V0h10l4 4v16H2v-1z"/>
          </svg>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-700">{{ fileName }}</p>
            <p class="text-xs text-gray-500">✓ File ready to submit</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ accept?: string; maxSize?: number; hint?: string }>()
const emit = defineEmits(['file-selected', 'error'])
const fileInput = ref<HTMLInputElement | null>(null)
const preview = ref<string | null>(null)
const fileName = ref('')
const isImage = ref(false)
const isPdf = ref(false)

const clear = () => {
  if (fileInput.value) fileInput.value.value = ''
  preview.value = null
  fileName.value = ''
  isImage.value = false
  isPdf.value = false
  emit('file-selected', { file: null, preview: null, fileName: '' })
}

defineExpose({ clear })

const onChange = (e: Event) => {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) {
    preview.value = null
    fileName.value = ''
    isImage.value = false
    isPdf.value = false
    emit('file-selected', { file: null, preview: null, fileName: '' })
    return
  }

  if (props.maxSize && file.size > props.maxSize) {
    emit('error', 'File size must be less than 5 MB')
    input.value = ''
    return
  }

  fileName.value = file.name
  isImage.value = file.type.startsWith('image/')
  isPdf.value = file.type === 'application/pdf'

  if (isImage.value) {
    const reader = new FileReader()
    reader.onload = (ev) => {
      preview.value = ev.target?.result as string
      emit('file-selected', { file, preview: preview.value, fileName: file.name })
      emit('error', '')
    }
    reader.readAsDataURL(file)
  } else if (isPdf.value) {
    preview.value = 'pdf'
    emit('file-selected', { file, preview: 'pdf', fileName: file.name })
    emit('error', '')
  } else {
    preview.value = null
    emit('file-selected', { file, preview: null, fileName: file.name })
    emit('error', '')
  }
}
</script>
