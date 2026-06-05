<template>
  <div class="request-list">
    <div
      v-for="req in requests"
      :key="req.id"
      class="request-item"
      :class="{ selected: req.id === selectedId }"
      @click="$emit('select', req.id)"
    >
      <div class="item-top">
        <span class="method-badge" :class="methodClass(req.method)">{{ req.method }}</span>
        <span class="item-path">{{ req.path }}</span>
        <button
          class="item-delete btn-ghost"
          @click.stop="$emit('delete', req.id)"
          title="Delete"
        >×</button>
      </div>
      <div class="item-bottom">
        <span class="item-time">{{ timeAgo(req.created_at) }}</span>
        <span v-if="req.content_type" class="item-ct">{{ shortType(req.content_type) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({ requests: Array, selectedId: String })
defineEmits(['select', 'delete'])

function methodClass(method) {
  const map = {
    GET: 'method-get', POST: 'method-post', PUT: 'method-put',
    PATCH: 'method-patch', DELETE: 'method-delete',
    HEAD: 'method-head', OPTIONS: 'method-options',
  }
  return map[method?.toUpperCase()] || 'method-other'
}

function timeAgo(dateStr) {
  const secs = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000)
  if (secs < 60)   return `${secs}s ago`
  if (secs < 3600) return `${Math.floor(secs / 60)}m ago`
  if (secs < 86400) return `${Math.floor(secs / 3600)}h ago`
  return new Date(dateStr).toLocaleDateString()
}

function shortType(ct) {
  if (!ct) return ''
  if (ct.includes('json'))             return 'JSON'
  if (ct.includes('xml'))              return 'XML'
  if (ct.includes('html'))             return 'HTML'
  if (ct.includes('form-urlencoded'))  return 'FORM'
  if (ct.includes('multipart'))        return 'MULTIPART'
  if (ct.includes('text/plain'))       return 'TEXT'
  return (ct.split('/')[1] || ct).split(';')[0].toUpperCase()
}
</script>

<style scoped>
.request-list {
  display: flex;
  flex-direction: column;
}

.request-item {
  padding: 8px 12px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: background 0.1s;
}
.request-item:hover        { background: var(--surface2); }
.request-item.selected     { background: var(--surface2); border-left: 2px solid var(--accent); padding-left: 10px; }

.item-top {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.item-path {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: 'SF Mono', monospace;
  font-size: 11px;
  color: var(--text-bright);
}

.item-delete {
  opacity: 0;
  font-size: 14px;
  line-height: 1;
  padding: 1px 5px;
  color: var(--text-dim);
  flex-shrink: 0;
}
.request-item:hover .item-delete { opacity: 1; }

.item-bottom {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 3px;
  padding-left: 1px;
}

.item-time {
  font-size: 10px;
  color: var(--text-dim);
}

.item-ct {
  font-size: 10px;
  color: var(--text-dim);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 3px;
  padding: 0 4px;
  font-family: monospace;
}
</style>
