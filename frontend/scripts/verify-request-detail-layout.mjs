import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, resolve } from 'node:path'
import assert from 'node:assert/strict'

const __dirname = dirname(fileURLToPath(import.meta.url))
const componentPath = resolve(__dirname, '../src/components/RequestDetail.vue')
const source = readFileSync(componentPath, 'utf8')

function cssRule(selector) {
  const escaped = selector.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const match = source.match(new RegExp(`${escaped}\\s*\\{([^}]*)\\}`, 'm'))
  assert.ok(match, `${selector} rule should exist`)
  return match[1]
}

function assertDeclaration(selector, property, value) {
  const declarations = cssRule(selector)
  assert.match(
    declarations,
    new RegExp(`${property}\\s*:\\s*${value}\\s*;`),
    `${selector} should declare ${property}: ${value}`
  )
}

assertDeclaration('.request-detail', 'min-height', '0')
assertDeclaration('.sections', 'min-height', '0')
assertDeclaration('section', 'flex-shrink', '0')
