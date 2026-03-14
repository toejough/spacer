export interface SM2State {
  easeFactor: number
  interval: number
  repetitions: number
}

export function sm2(quality: number, state: SM2State): SM2State {
  if (quality < 0 || quality > 5) {
    throw new RangeError('Quality must be between 0 and 5')
  }

  const newEF = Math.max(
    1.3,
    state.easeFactor + (0.1 - (5 - quality) * (0.08 + (5 - quality) * 0.02)),
  )

  if (quality < 3) {
    return { easeFactor: newEF, interval: 1, repetitions: 0 }
  }

  let interval: number
  if (state.repetitions === 0) {
    interval = 1
  } else if (state.repetitions === 1) {
    interval = 6
  } else {
    interval = Math.round(state.interval * state.easeFactor)
  }

  return {
    easeFactor: newEF,
    interval,
    repetitions: state.repetitions + 1,
  }
}
