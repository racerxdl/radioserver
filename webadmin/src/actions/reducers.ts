import {DefinedActions} from "./actions";
import {AddFFTInitialState} from "./initialStates";
import {combineReducers} from "redux";


function fftSamples(state: any | void | null, action: any) {
  if (action.type === DefinedActions.AddFFT) {
    const s = !state ? AddFFTInitialState : state;
    return {
      ...s,
      samples: action.samples,
      sampleRate: action.sampleRate,
      centerFrequency: action.centerFrequency,
    };
  }

  return state || AddFFTInitialState;
}

export default combineReducers({
  fftSamples,
})
