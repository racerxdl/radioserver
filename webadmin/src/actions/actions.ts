import {AddFFTAction} from "./types";

const DefinedActions = {
  AddFFT: 'ADD_FFT',
};

function AddFFT(centerFrequency: number, sampleRate: number, samples: number[]): AddFFTAction {
  return {
    type: DefinedActions.AddFFT,
    samples,
    centerFrequency,
    sampleRate,
  }
}

export {
  DefinedActions,
  AddFFT,
}
