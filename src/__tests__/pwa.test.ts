import { describe, it, expect } from 'vitest'
import { readFileSync, existsSync } from 'fs'
import { resolve } from 'path'
import { execSync } from 'child_process'

describe('PWA manifest', () => {
  it('build generates a web manifest with correct metadata', () => {
    execSync('npx vite build', { stdio: 'pipe' })
    const distDir = resolve(__dirname, '../../dist')

    const manifestPath = resolve(distDir, 'manifest.webmanifest')
    expect(existsSync(manifestPath)).toBe(true)

    const manifest = JSON.parse(readFileSync(manifestPath, 'utf-8'))
    expect(manifest.name).toBe('Spacer')
    expect(manifest.display).toBe('standalone')
    expect(manifest.theme_color).toBe('#6366f1')
  })
})
