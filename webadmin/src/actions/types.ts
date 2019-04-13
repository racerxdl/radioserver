import {DeviceType} from "../radioserver/protocol/server_pb";

export type ActionType = {
  type: string
}

export type AddFFTState = {
  samples: number[]
  centerFrequency: number,
  sampleRate: number,
}

export type DeviceInfo = {
  deviceType: DeviceType,
  deviceSerial: string,
  deviceName: string,
  maximumSampleRate: number,
  maximumGain: number,
  maximumDecimation: number,
  minimumFrequency: number,
  maximumFrequency: number,
  adcResolution: number,
}

export type VersionInfo = {
  major: number,
  minor: number,
  hash: number,
}

export type AddServerInfoState = {
  controlAllowed: boolean,
  serverCenterFrequency: number,
  minimumIQCenterFrequency: number,
  maximumIQcenterFrequency: number,
  minimumSmartFrequency: number,
  maximumSmartFrequency: number,
  deviceInfo?: DeviceInfo
  version?: VersionInfo,
}

export type AddFFTAction = ActionType & AddFFTState
export type AddServerInfoAction = ActionType & AddServerInfoState
