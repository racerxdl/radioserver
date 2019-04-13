import {AddFFTAction, AddServerInfoAction} from "./types";
import {ServerInfoData} from "../radioserver/protocol/server_pb";

const DefinedActions = {
  AddFFT: 'ADD_FFT',
  AddServerInfo: 'ADD_SERVER_INFO'
};

function AddFFT(centerFrequency: number, sampleRate: number, samples: number[]): AddFFTAction {
  return {
    type: DefinedActions.AddFFT,
    samples,
    centerFrequency,
    sampleRate,
  }
}

function AddServerInfo(serverInfo: ServerInfoData.AsObject): AddServerInfoAction {
  const deviceInfo = serverInfo.deviceinfo ? {
    deviceType: serverInfo.deviceinfo.devicetype,
    deviceSerial: serverInfo.deviceinfo.deviceserial,
    deviceName: serverInfo.deviceinfo.devicename,
    maximumSampleRate: serverInfo.deviceinfo.maximumsamplerate,
    maximumGain: serverInfo.deviceinfo.maximumgain,
    maximumDecimation: serverInfo.deviceinfo.maximumdecimation,
    minimumFrequency: serverInfo.deviceinfo.minimumfrequency,
    maximumFrequency: serverInfo.deviceinfo.maximumfrequency,
    adcResolution: serverInfo.deviceinfo.adcresolution,
  } : undefined;

  const version = serverInfo.version ? {
    major: serverInfo.version.major,
    minor: serverInfo.version.minor,
    hash: serverInfo.version.hash,
  } : undefined;

  return {
    type: DefinedActions.AddServerInfo,
    controlAllowed: serverInfo.controlallowed,
    serverCenterFrequency: serverInfo.servercenterfrequency,
    minimumIQCenterFrequency: serverInfo.minimumiqcenterfrequency,
    maximumIQcenterFrequency: serverInfo.maximumiqcenterfrequency,
    minimumSmartFrequency: serverInfo.minimumsmartfrequency,
    maximumSmartFrequency: serverInfo.maximumsmartfrequency,
    deviceInfo,
    version,
  }
}

export {
  DefinedActions,
  AddFFT,
  AddServerInfo,
}
