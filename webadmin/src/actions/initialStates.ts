import {AddFFTState, AddServerInfoState} from "./types";
import {DeviceType} from "../radioserver/protocol/server_pb";

const AddFFTInitialState: AddFFTState = {
  samples: [],
  centerFrequency: 0,
  sampleRate: 0,
};

const ServerInfoInitialState: AddServerInfoState = {
  controlAllowed: false,
  serverCenterFrequency: 0,
  minimumIQCenterFrequency: 0,
  maximumIQcenterFrequency: 0,
  minimumSmartFrequency: 0,
  maximumSmartFrequency: 0,
  deviceInfo: {
    deviceType: DeviceType.DEVICEINVALID,
    deviceSerial: '',
    deviceName: 'Not connected',
    maximumSampleRate: 0,
    maximumGain: 0,
    maximumDecimation: 0,
    minimumFrequency: 0,
    maximumFrequency: 0,
    adcResolution: 0,
  },
  version: {
    major: 0,
    minor: 0,
    hash: 0,
  },
};

export {
  AddFFTInitialState,
  ServerInfoInitialState,
}
