import {AddServerInfo, DefinedActions} from "./actions";
import {AddFFTInitialState, ServerInfoInitialState} from "./initialStates";
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

function serverInfo(state: any | void | null, action: any) {
  if (action.type === DefinedActions.AddServerInfo) {
    const s = !state ? ServerInfoInitialState : state;
    return {
      ...s,
      ...AddServerInfo(action),
    }
  }

  return state || ServerInfoInitialState;
}

export default combineReducers({
  fftSamples,
  serverInfo,
})
