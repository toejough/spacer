import { describe, it, expect } from 'vitest'
import { sm2, type SM2State } from '../sm2'

const freshCard: SM2State = {
  easeFactor: 2.5,
  interval: 0,
  repetitions: 0,
}

describe('sm2', () => {
  it('resets on quality < 3', () => {
    const result = sm2(2, { easeFactor: 2.5, interval: 10, repetitions: 5 })
    expect(result.repetitions).toBe(0)
    expect(result.interval).toBe(1)
  })

  it('adjusts ease factor even on failed reviews', () => {
    const result = sm2(1, { easeFactor: 2.5, interval: 10, repetitions: 5 })
    expect(result.easeFactor).not.toBe(2.5)
  })

  it('sets interval to 1 on first successful review', () => {
    const result = sm2(4, freshCard)
    expect(result.repetitions).toBe(1)
    expect(result.interval).toBe(1)
  })

  it('sets interval to 6 on second successful review', () => {
    const result = sm2(4, { easeFactor: 2.5, interval: 1, repetitions: 1 })
    expect(result.repetitions).toBe(2)
    expect(result.interval).toBe(6)
  })

  it('multiplies interval by easeFactor on third+ review', () => {
    const result = sm2(4, { easeFactor: 2.5, interval: 6, repetitions: 2 })
    expect(result.repetitions).toBe(3)
    expect(result.interval).toBe(15) // Math.round(6 * 2.5)
  })

  it('adjusts ease factor based on quality', () => {
    const result = sm2(5, freshCard)
    expect(result.easeFactor).toBeCloseTo(2.6, 1)
  })

  it('never lets ease factor drop below 1.3', () => {
    const result = sm2(3, { easeFactor: 1.3, interval: 1, repetitions: 1 })
    expect(result.easeFactor).toBeGreaterThanOrEqual(1.3)
  })
})
