// package: protocol
// file: protocol/server.proto

import * as jspb from "google-protobuf";

export class LoginData extends jspb.Message {
  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginData.AsObject;
  static toObject(includeInstance: boolean, msg: LoginData): LoginData.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LoginData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginData;
  static deserializeBinaryFromReader(message: LoginData, reader: jspb.BinaryReader): LoginData;
}

export namespace LoginData {
  export type AsObject = {
    token: string,
  }
}

export class IQData extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): void;

  getStatus(): StatusType;
  setStatus(value: StatusType): void;

  getError(): string;
  setError(value: string): void;

  getType(): ChannelType;
  setType(value: ChannelType): void;

  clearSamplesList(): void;
  getSamplesList(): Array<number>;
  setSamplesList(value: Array<number>): void;
  addSamples(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IQData.AsObject;
  static toObject(includeInstance: boolean, msg: IQData): IQData.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: IQData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IQData;
  static deserializeBinaryFromReader(message: IQData, reader: jspb.BinaryReader): IQData;
}

export namespace IQData {
  export type AsObject = {
    timestamp: number,
    status: StatusType,
    error: string,
    type: ChannelType,
    samplesList: Array<number>,
  }
}

export class ChannelConfig extends jspb.Message {
  hasLogininfo(): boolean;
  clearLogininfo(): void;
  getLogininfo(): LoginData | undefined;
  setLogininfo(value?: LoginData): void;

  getCenterfrequency(): number;
  setCenterfrequency(value: number): void;

  getDecimationstage(): number;
  setDecimationstage(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChannelConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ChannelConfig): ChannelConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChannelConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChannelConfig;
  static deserializeBinaryFromReader(message: ChannelConfig, reader: jspb.BinaryReader): ChannelConfig;
}

export namespace ChannelConfig {
  export type AsObject = {
    logininfo?: LoginData.AsObject,
    centerfrequency: number,
    decimationstage: number,
  }
}

export class FrequencyChannelConfig extends jspb.Message {
  hasLogininfo(): boolean;
  clearLogininfo(): void;
  getLogininfo(): LoginData | undefined;
  setLogininfo(value?: LoginData): void;

  getCenterfrequency(): number;
  setCenterfrequency(value: number): void;

  getDecimationstage(): number;
  setDecimationstage(value: number): void;

  getLength(): number;
  setLength(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FrequencyChannelConfig.AsObject;
  static toObject(includeInstance: boolean, msg: FrequencyChannelConfig): FrequencyChannelConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FrequencyChannelConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FrequencyChannelConfig;
  static deserializeBinaryFromReader(message: FrequencyChannelConfig, reader: jspb.BinaryReader): FrequencyChannelConfig;
}

export namespace FrequencyChannelConfig {
  export type AsObject = {
    logininfo?: LoginData.AsObject,
    centerfrequency: number,
    decimationstage: number,
    length: number,
  }
}

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}

export class PingData extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): void;

  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PingData.AsObject;
  static toObject(includeInstance: boolean, msg: PingData): PingData.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PingData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PingData;
  static deserializeBinaryFromReader(message: PingData, reader: jspb.BinaryReader): PingData;
}

export namespace PingData {
  export type AsObject = {
    timestamp: number,
    token: string,
  }
}

export class HelloData extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getApplication(): string;
  setApplication(value: string): void;

  getUsername(): string;
  setUsername(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HelloData.AsObject;
  static toObject(includeInstance: boolean, msg: HelloData): HelloData.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HelloData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HelloData;
  static deserializeBinaryFromReader(message: HelloData, reader: jspb.BinaryReader): HelloData;
}

export namespace HelloData {
  export type AsObject = {
    name: string,
    application: string,
    username: string,
    password: string,
  }
}

export class HelloReturn extends jspb.Message {
  getStatus(): StatusType;
  setStatus(value: StatusType): void;

  hasLogin(): boolean;
  clearLogin(): void;
  getLogin(): LoginData | undefined;
  setLogin(value?: LoginData): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HelloReturn.AsObject;
  static toObject(includeInstance: boolean, msg: HelloReturn): HelloReturn.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HelloReturn, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HelloReturn;
  static deserializeBinaryFromReader(message: HelloReturn, reader: jspb.BinaryReader): HelloReturn;
}

export namespace HelloReturn {
  export type AsObject = {
    status: StatusType,
    login?: LoginData.AsObject,
  }
}

export class VersionData extends jspb.Message {
  getMajor(): number;
  setMajor(value: number): void;

  getMinor(): number;
  setMinor(value: number): void;

  getHash(): number;
  setHash(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VersionData.AsObject;
  static toObject(includeInstance: boolean, msg: VersionData): VersionData.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VersionData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VersionData;
  static deserializeBinaryFromReader(message: VersionData, reader: jspb.BinaryReader): VersionData;
}

export namespace VersionData {
  export type AsObject = {
    major: number,
    minor: number,
    hash: number,
  }
}

export class DeviceInfo extends jspb.Message {
  getDevicetype(): DeviceType;
  setDevicetype(value: DeviceType): void;

  getDeviceserial(): string;
  setDeviceserial(value: string): void;

  getDevicename(): string;
  setDevicename(value: string): void;

  getMaximumsamplerate(): number;
  setMaximumsamplerate(value: number): void;

  getMaximumgain(): number;
  setMaximumgain(value: number): void;

  getMaximumdecimation(): number;
  setMaximumdecimation(value: number): void;

  getMinimumfrequency(): number;
  setMinimumfrequency(value: number): void;

  getMaximumfrequency(): number;
  setMaximumfrequency(value: number): void;

  getAdcresolution(): number;
  setAdcresolution(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeviceInfo.AsObject;
  static toObject(includeInstance: boolean, msg: DeviceInfo): DeviceInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeviceInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeviceInfo;
  static deserializeBinaryFromReader(message: DeviceInfo, reader: jspb.BinaryReader): DeviceInfo;
}

export namespace DeviceInfo {
  export type AsObject = {
    devicetype: DeviceType,
    deviceserial: string,
    devicename: string,
    maximumsamplerate: number,
    maximumgain: number,
    maximumdecimation: number,
    minimumfrequency: number,
    maximumfrequency: number,
    adcresolution: number,
  }
}

export class ServerInfoData extends jspb.Message {
  getControlallowed(): boolean;
  setControlallowed(value: boolean): void;

  getServercenterfrequency(): number;
  setServercenterfrequency(value: number): void;

  getMinimumiqcenterfrequency(): number;
  setMinimumiqcenterfrequency(value: number): void;

  getMaximumiqcenterfrequency(): number;
  setMaximumiqcenterfrequency(value: number): void;

  getMinimumsmartfrequency(): number;
  setMinimumsmartfrequency(value: number): void;

  getMaximumsmartfrequency(): number;
  setMaximumsmartfrequency(value: number): void;

  hasDeviceinfo(): boolean;
  clearDeviceinfo(): void;
  getDeviceinfo(): DeviceInfo | undefined;
  setDeviceinfo(value?: DeviceInfo): void;

  hasVersion(): boolean;
  clearVersion(): void;
  getVersion(): VersionData | undefined;
  setVersion(value?: VersionData): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerInfoData.AsObject;
  static toObject(includeInstance: boolean, msg: ServerInfoData): ServerInfoData.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServerInfoData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerInfoData;
  static deserializeBinaryFromReader(message: ServerInfoData, reader: jspb.BinaryReader): ServerInfoData;
}

export namespace ServerInfoData {
  export type AsObject = {
    controlallowed: boolean,
    servercenterfrequency: number,
    minimumiqcenterfrequency: number,
    maximumiqcenterfrequency: number,
    minimumsmartfrequency: number,
    maximumsmartfrequency: number,
    deviceinfo?: DeviceInfo.AsObject,
    version?: VersionData.AsObject,
  }
}

export class MeData extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getApplication(): string;
  setApplication(value: string): void;

  getUsername(): string;
  setUsername(value: string): void;

  getAddress(): string;
  setAddress(value: string): void;

  getControlallowed(): boolean;
  setControlallowed(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MeData.AsObject;
  static toObject(includeInstance: boolean, msg: MeData): MeData.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MeData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MeData;
  static deserializeBinaryFromReader(message: MeData, reader: jspb.BinaryReader): MeData;
}

export namespace MeData {
  export type AsObject = {
    name: string,
    application: string,
    username: string,
    address: string,
    controlallowed: boolean,
  }
}

export class ByeReturn extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ByeReturn.AsObject;
  static toObject(includeInstance: boolean, msg: ByeReturn): ByeReturn.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ByeReturn, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ByeReturn;
  static deserializeBinaryFromReader(message: ByeReturn, reader: jspb.BinaryReader): ByeReturn;
}

export namespace ByeReturn {
  export type AsObject = {
    message: string,
  }
}

export enum StatusType {
  INVALID = 0,
  OK = 1,
  ERROR = 2,
}

export enum ChannelType {
  IQ = 0,
  SMARTIQ = 1,
  FREQUENCY = 2,
}

export enum DeviceType {
  DEVICEINVALID = 0,
  DEVICETESTSIGNAL = 1,
  DEVICEAIRSPYONE = 2,
  DEVICEAIRSPYHF = 3,
  DEVICERTLSDR = 4,
  DEVICELIMESDRUSB = 5,
  DEVICELIMESDRMINI = 6,
  DEVICESPYSERVER = 7,
  DEVICEHACKRF = 8,
}

