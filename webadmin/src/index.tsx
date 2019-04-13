import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';
import RadioClient from "./radioserver/client";
import {AddFFT} from "./actions/actions";
import appReducers from "./actions/reducers";
import {createStore} from "redux";
import {Provider} from "react-redux";

const store = createStore(appReducers);
const client = new RadioClient('https://127.0.0.1:8001', store);

const centerFrequency = 97700000;
const sampleRate = 2500000;

function onFFT(samples: number[] | null, error: string | null) {
  if (error) {
    console.log(`Error: ${error}`);
  } else if (samples !== null) {
    store.dispatch(AddFFT(centerFrequency, sampleRate, samples));
  }
}

client.SetOnFFT(onFFT);

(async () => {
  console.log(`Logging into RadioServer`);
  await client.Login();
  client.StartLoop();
  console.log(`Starting FFT`);
  await client.StartFFT(centerFrequency, 0, 1024);
  // draw();
  // await sleep(5000);
  // await client.StopSmartIQ();
})();

window.onbeforeunload = () => {
  console.log(`Logging out`);
  client.StopLoop();
  client.StopSmartIQ();
  client.StopFFT();
  client.Logout();
};

ReactDOM.render((
  <Provider store={store}>
    <App/>
  </Provider>
), document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();

