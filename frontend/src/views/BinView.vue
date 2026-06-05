<template>
  <div v-if="notFound" class="not-found">
    <h2>Bin not found</h2>
    <p>This bin may have expired or been deleted.</p>
    <router-link to="/" class="btn-primary" style="display:inline-block;padding:8px 16px;border-radius:6px;">
      Create a new bin
    </router-link>
  </div>

  <div v-else class="bin-view">
    <header class="bin-header">
      <router-link to="/" class="home-link">
        <svg width="18" height="18" viewBox="0 0 18 18" fill="none">
          <rect width="18" height="18" rx="4" fill="#7c3aed" fill-opacity="0.2"/>
          <path d="M4 7h10M4 9.5h6M4 12h8" stroke="#7c3aed" stroke-width="1.5" stroke-linecap="round"/>
        </svg>
        <span>RequestBin</span>
      </router-link>

      <div class="bin-url-row">
        <span class="url-label">Bin URL</span>
        <code class="url-value">{{ captureUrl }}</code>
        <button class="btn-ghost copy-btn" @click="copyUrl" :title="copied ? 'Copied!' : 'Copy'">
          {{ copied ? '✓ Copied' : '⎘ Copy' }}
        </button>
      </div>

      <div class="header-actions">
        <span class="req-count">{{ requests.length }} request{{ requests.length !== 1 ? 's' : '' }}</span>
        <button class="btn-ghost" @click="clearAll" :disabled="!requests.length">Clear all</button>
        <button class="btn-ghost delete-bin-btn" @click="deleteBin">Delete bin</button>
      </div>
    </header>

    <div class="content">
      <aside class="sidebar">
        <div v-if="!requests.length" class="empty-list">
          <svg width="32" height="32" viewBox="0 0 32 32" fill="none" style="margin-bottom:8px;opacity:0.3">
            <path d="M6 8h20M6 14h14M6 20h10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <p>Waiting for requests…</p>
          <p class="empty-hint">Send a request to:</p>
          <code class="empty-url">{{ captureUrl }}</code>
        </div>
        <RequestList
          v-else
          :requests="requests"
          :selected-id="selectedId"
          @select="selectedId = $event"
          @delete="deleteRequest"
        />
      </aside>

      <main class="detail-pane">
        <RequestDetail v-if="selectedRequest" :request="selectedRequest" />
        <div v-else class="empty-detail">
          <svg width="40" height="40" viewBox="0 0 40 40" fill="none" style="opacity:0.2;margin-bottom:10px">
            <rect x="6" y="8" width="28" height="24" rx="3" stroke="currentColor" stroke-width="2"/>
            <path d="M11 14h18M11 19h12M11 24h15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <p>Select a request to inspect it</p>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import RequestList from '../components/RequestList.vue'
import RequestDetail from '../components/RequestDetail.vue'

const route = useRoute()
const router = useRouter()

const binId = computed(() => route.params.id)
const captureUrl = computed(() => `${window.location.origin}/r/${binId.value}`)

const requests = ref([])
const selectedId = ref(null)
const notFound = ref(false)
const copied = ref(false)

const selectedRequest = computed(() =>
  requests.value.find(r => r.id === selectedId.value) ?? null
)

let eventSource = null

async function loadRequests() {
  try {
    const res = await fetch(`/api/bins/${binId.value}/requests`)
    if (res.status === 404) { notFound.value = true; return }
    if (!res.ok) return
    const data = await res.json()
    requests.value = data
    if (data.length && !selectedId.value) selectedId.value = data[0].id
  } catch {}
}

function setupSSE() {
  eventSource?.close()
  eventSource = new EventSource(`/api/bins/${binId.value}/sse`)
  eventSource.addEventListener('request', (e) => {
    const req = JSON.parse(e.data)
    requests.value.unshift(req)
    if (!selectedId.value) selectedId.value = req.id
  })
}

async function copyUrl() {
  await navigator.clipboard.writeText(captureUrl.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}

async function clearAll() {
  if (!requests.value.length) return
  await fetch(`/api/bins/${binId.value}/requests`, { method: 'DELETE' })
  requests.value = []
  selectedId.value = null
}

async function deleteBin() {
  if (!confirm('Delete this bin and all its requests?')) return
  await fetch(`/api/bins/${binId.value}`, { method: 'DELETE' })
  router.push('/')
}

async function deleteRequest(id) {
  await fetch(`/api/bins/${binId.value}/requests/${id}`, { method: 'DELETE' })
  requests.value = requests.value.filter(r => r.id !== id)
  if (selectedId.value === id) {
    selectedId.value = requests.value[0]?.id ?? null
  }
}

onMounted(async () => {
  await loadRequests()
  if (!notFound.value) setupSSE()
})

onUnmounted(() => {
  eventSource?.close()
})
</script>

<style scoped>
.not-found {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--text-dim);
}
.not-found h2 { color: var(--text-bright); font-size: 20px; }

.bin-view {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.bin-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 0 16px;
  height: 48px;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.home-link {
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
  color: var(--text-bright);
  font-weight: 600;
  font-size: 13px;
  flex-shrink: 0;
}
.home-link:hover { text-decoration: none; opacity: 0.8; }

.bin-url-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 4px 10px;
}

.url-label {
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--text-dim);
  letter-spacing: 0.5px;
  flex-shrink: 0;
}

.url-value {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--text-bright);
  font-size: 12px;
}

.copy-btn {
  flex-shrink: 0;
  font-size: 11px;
  padding: 2px 8px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.req-count {
  font-size: 11px;
  color: var(--text-dim);
  white-space: nowrap;
}

.delete-bin-btn {
  color: var(--red);
}
.delete-bin-btn:hover {
  background: rgba(239, 68, 68, 0.1);
}

.content {
  flex: 1;
  display: flex;
  min-height: 0;
}

.sidebar {
  width: 280px;
  flex-shrink: 0;
  border-right: 1px solid var(--border);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.detail-pane {
  flex: 1;
  overflow-y: auto;
  min-width: 0;
}

.empty-list {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  text-align: center;
  color: var(--text-dim);
  gap: 4px;
}

.empty-hint {
  margin-top: 8px;
  font-size: 11px;
}

.empty-url {
  font-size: 11px;
  color: var(--accent);
  word-break: break-all;
}

.empty-detail {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text-dim);
  gap: 4px;
  font-size: 13px;
}
</style>
