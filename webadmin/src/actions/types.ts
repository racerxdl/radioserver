export type ActionType = {
  type: string
}

export type AddFFTState = {
  samples: number[]
  centerFrequency: number,
  sampleRate: number,
}

export type AddFFTAction = ActionType & AddFFTState
