<template>
  <div class="request-detail">
    <div class="detail-top">
      <span class="method-badge" :class="methodClass(request.method)">{{ request.method }}</span>
      <span class="detail-path">{{ request.path }}</span>
      <span class="detail-time">{{ formatDate(request.created_at) }}</span>
    </div>

    <div class="sections">

      <section>
        <h3>Overview</h3>
        <table class="kv-table">
          <tbody>
            <tr>
              <td>Remote Address</td>
              <td><code>{{ request.remote_addr }}</code></td>
            </tr>
            <tr v-if="request.content_type">
              <td>Content-Type</td>
              <td><code>{{ request.content_type }}</code></td>
            </tr>
            <tr v-if="request.content_length > 0">
              <td>Content-Length</td>
              <td><code>{{ request.content_length }} bytes</code></td>
            </tr>
          </tbody>
        </table>
      </section>

      <section v-if="hasQueryParams">
        <h3>Query Parameters</h3>
        <table class="kv-table">
          <tbody>
            <tr v-for="(val, key) in request.query_params" :key="key">
              <td>{{ key }}</td>
              <td><code>{{ val }}</code></td>
            </tr>
          </tbody>
        </table>
      </section>

      <section v-if="hasHeaders">
        <h3>Headers</h3>
        <table class="kv-table">
          <tbody>
            <tr v-for="(val, key) in sortedHeaders" :key="key">
              <td>{{ key }}</td>
              <td><code>{{ val }}</code></td>
            </tr>
          </tbody>
        </table>
      </section>

      <section v-if="request.body">
        <h3>
          Body
          <button class="btn-small" @click="copyBody">
            {{ bodyCopied ? '✓ Copied' : 'Copy' }}
          </button>
        </h3>
        <pre class="body-pre"><code>{{ formattedBody }}</code></pre>
      </section>

    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({ request: Object })

const bodyCopied = ref(false)

const hasHeaders = computed(() =>
  props.request.headers && Object.keys(props.request.headers).length > 0
)

const hasQueryParams = computed(() =>
  props.request.query_params && Object.keys(props.request.query_params).length > 0
)

const sortedHeaders = computed(() => {
  const h = props.request.headers || {}
  return Object.fromEntries(
    Object.entries(h).sort(([a], [b]) => a.localeCompare(b))
  )
})

const formattedBody = computed(() => {
  const body = props.request.body
  if (!body) return ''
  const ct = (props.request.content_type || '').toLowerCase()
  if (ct.includes('json')) {
    try { return JSON.stringify(JSON.parse(body), null, 2) } catch {}
  }
  return body
})

function methodClass(method) {
  const map = {
    GET: 'method-get', POST: 'method-post', PUT: 'method-put',
    PATCH: 'method-patch', DELETE: 'method-delete',
    HEAD: 'method-head', OPTIONS: 'method-options',
  }
  return map[method?.toUpperCase()] || 'method-other'
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleString()
}

async function copyBody() {
  await navigator.clipboard.writeText(props.request.body)
  bodyCopied.value = true
  setTimeout(() => { bodyCopied.value = false }, 2000)
}
</script>

<style scoped>
.request-detail {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.detail-top {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
  background: var(--surface);
  flex-shrink: 0;
}

.detail-path {
  flex: 1;
  font-family: 'SF Mono', monospace;
  font-size: 12px;
  color: var(--text-bright);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.detail-time {
  font-size: 11px;
  color: var(--text-dim);
  flex-shrink: 0;
}

.sections {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

section {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
}

h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.6px;
  color: var(--text-dim);
  padding: 8px 12px;
  border-bottom: 1px solid var(--border);
  background: var(--surface2);
}

.kv-table {
  width: 100%;
  border-collapse: collapse;
}

.kv-table tr {
  border-bottom: 1px solid var(--border);
}
.kv-table tr:last-child {
  border-bottom: none;
}

.kv-table td {
  padding: 6px 12px;
  vertical-align: top;
}

.kv-table td:first-child {
  width: 35%;
  color: var(--text-dim);
  font-size: 12px;
  white-space: nowrap;
}

.kv-table td:last-child {
  color: var(--text-bright);
}

.kv-table code {
  color: var(--text-bright);
  font-size: 11.5px;
  word-break: break-all;
}

.body-pre {
  padding: 12px;
  margin: 0;
  overflow-x: auto;
  white-space: pre;
  background: var(--bg);
}

.body-pre code {
  font-size: 12px;
  color: var(--text-bright);
  line-height: 1.6;
}
</style>
