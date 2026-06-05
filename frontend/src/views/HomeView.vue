<template>
  <div class="home">
    <div class="hero">
      <div class="logo-mark">
        <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
          <rect width="40" height="40" rx="10" fill="#7c3aed" fill-opacity="0.15"/>
          <path d="M10 14h20M10 20h12M10 26h16" stroke="#7c3aed" stroke-width="2.5" stroke-linecap="round"/>
          <circle cx="30" cy="26" r="5" fill="#7c3aed" fill-opacity="0.3" stroke="#7c3aed" stroke-width="1.5"/>
          <path d="M28.5 26l1 1 2-2" stroke="#7c3aed" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </div>
      <h1>RequestBin</h1>
      <p class="tagline">Inspect and debug HTTP requests in real-time</p>
      <button class="btn-primary create-btn" @click="createBin" :disabled="loading">
        {{ loading ? 'Creating…' : '+ Create a new bin' }}
      </button>
      <p v-if="error" class="error-msg">{{ error }}</p>
    </div>

    <div v-if="recentBins.length" class="recent">
      <h2>Recent bins</h2>
      <div class="bin-list">
        <router-link
          v-for="bin in recentBins"
          :key="bin.id"
          :to="`/b/${bin.id}`"
          class="bin-card"
        >
          <span class="bin-id">{{ bin.id }}</span>
          <span class="bin-time">{{ formatDate(bin.createdAt) }}</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const recentBins = ref([])

onMounted(() => {
  const stored = localStorage.getItem('rb_recent_bins')
  if (stored) {
    try { recentBins.value = JSON.parse(stored) } catch {}
  }
})

async function createBin() {
  loading.value = true
  error.value = ''
  try {
    const res = await fetch('/api/bins', { method: 'POST' })
    if (!res.ok) throw new Error('Failed to create bin')
    const bin = await res.json()

    const recent = JSON.parse(localStorage.getItem('rb_recent_bins') || '[]')
    recent.unshift({ id: bin.id, createdAt: bin.created_at })
    localStorage.setItem('rb_recent_bins', JSON.stringify(recent.slice(0, 10)))

    router.push(`/b/${bin.id}`)
  } catch (e) {
    error.value = e.message
    loading.value = false
  }
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleString()
}
</script>

<style scoped>
.home {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 48px;
  padding: 24px;
}

.hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  text-align: center;
}

.logo-mark {
  margin-bottom: 4px;
}

h1 {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-bright);
  letter-spacing: -0.5px;
}

.tagline {
  color: var(--text-dim);
  font-size: 15px;
  margin-bottom: 8px;
}

.create-btn {
  padding: 10px 24px;
  font-size: 14px;
  border-radius: 8px;
}

.error-msg {
  color: var(--red);
  font-size: 12px;
}

.recent {
  width: 100%;
  max-width: 480px;
}

.recent h2 {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  color: var(--text-dim);
  margin-bottom: 10px;
}

.bin-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.bin-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  text-decoration: none;
  transition: border-color 0.15s, background 0.15s;
}
.bin-card:hover {
  border-color: var(--accent);
  background: var(--surface2);
  text-decoration: none;
}

.bin-id {
  font-family: 'SF Mono', monospace;
  font-size: 12px;
  color: var(--text-bright);
}

.bin-time {
  font-size: 11px;
  color: var(--text-dim);
}
</style>
