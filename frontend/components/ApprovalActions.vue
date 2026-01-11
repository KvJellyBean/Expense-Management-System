<template>
  <div class="border-t pt-4">
    <label class="block text-sm font-medium text-gray-700 mb-2">Notes (Optional)</label>
    <textarea
      :value="notes"
      @input="$emit('update:notes', $event.target.value)"
      rows="2"
      class="input mb-3"
      placeholder="Add approval notes..."
    ></textarea>

    <div class="flex flex-col sm:flex-row gap-2 sm:gap-3">
      <button @click="$emit('approve')" :disabled="processing || disabled" class="btn btn-success w-full sm:w-auto">
        {{ processing ? 'Processing...' : 'Approve' }}
      </button>
      <button @click="$emit('reject')" :disabled="processing || disabled" class="btn btn-danger w-full sm:w-auto">
        {{ processing ? 'Processing...' : 'Reject' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ notes?: string; processing?: boolean; disabled?: boolean }>()
defineEmits(['update:notes', 'approve', 'reject'])
</script>
