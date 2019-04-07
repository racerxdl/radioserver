export type ActionType = {
  type: string
}

export type AddFFTState = {
  samples: number[]
}

export type AddFFTAction = ActionType & AddFFTState
