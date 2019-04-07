import {AddFFTAction} from "./types";

const DefinedActions = {
  AddFFT: 'ADD_FFT',
};

function AddFFT(samples: number[]): AddFFTAction {
  return {
    type: DefinedActions.AddFFT,
    samples,
  }
}

export {
  DefinedActions,
  AddFFT,
}
