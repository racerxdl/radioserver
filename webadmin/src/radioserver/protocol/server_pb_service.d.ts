// package: protocol
// file: protocol/server.proto

import * as protocol_server_pb from "../protocol/server_pb";
import {grpc} from "@improbable-eng/grpc-web";

type RadioServerMe = {
  readonly methodName: string;
  readonly service: typeof RadioServer;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof protocol_server_pb.LoginData;
  readonly responseType: typeof protocol_server_pb.MeData;
};

type RadioServerHello = {
  readonly methodName: string;
  readonly service: typeof RadioServer;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof protocol_server_pb.HelloData;
  readonly responseType: typeof protocol_server_pb.HelloReturn;
};

type RadioServerBye = {
  readonly methodName: string;
  readonly service: typeof RadioServer;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof protocol_server_pb.LoginData;
  readonly responseType: typeof protocol_server_pb.ByeReturn;
};

type RadioServerServerInfo = {
  readonly methodName: string;
  readonly service: typeof RadioServer;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof protocol_server_pb.Empty;
  readonly responseType: typeof protocol_server_pb.ServerInfoData;
};

type RadioServerSmartIQ = {
  readonly methodName: string;
  readonly service: typeof RadioServer;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof protocol_server_pb.ChannelConfig;
  readonly responseType: typeof protocol_server_pb.IQData;
};

type RadioServerFrequencyData = {
  readonly methodName: string;
  readonly service: typeof RadioServer;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof protocol_server_pb.FrequencyChannelConfig;
  readonly responseType: typeof protocol_server_pb.IQData;
};

type RadioServerIQ = {
  readonly methodName: string;
  readonly service: typeof RadioServer;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof protocol_server_pb.ChannelConfig;
  readonly responseType: typeof protocol_server_pb.IQData;
};

type RadioServerPing = {
  readonly methodName: string;
  readonly service: typeof RadioServer;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof protocol_server_pb.PingData;
  readonly responseType: typeof protocol_server_pb.PingData;
};

export class RadioServer {
  static readonly serviceName: string;
  static readonly Me: RadioServerMe;
  static readonly Hello: RadioServerHello;
  static readonly Bye: RadioServerBye;
  static readonly ServerInfo: RadioServerServerInfo;
  static readonly SmartIQ: RadioServerSmartIQ;
  static readonly FrequencyData: RadioServerFrequencyData;
  static readonly IQ: RadioServerIQ;
  static readonly Ping: RadioServerPing;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: () => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: () => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: () => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class RadioServerClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  me(
    requestMessage: protocol_server_pb.LoginData,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.MeData|null) => void
  ): UnaryResponse;
  me(
    requestMessage: protocol_server_pb.LoginData,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.MeData|null) => void
  ): UnaryResponse;
  hello(
    requestMessage: protocol_server_pb.HelloData,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.HelloReturn|null) => void
  ): UnaryResponse;
  hello(
    requestMessage: protocol_server_pb.HelloData,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.HelloReturn|null) => void
  ): UnaryResponse;
  bye(
    requestMessage: protocol_server_pb.LoginData,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.ByeReturn|null) => void
  ): UnaryResponse;
  bye(
    requestMessage: protocol_server_pb.LoginData,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.ByeReturn|null) => void
  ): UnaryResponse;
  serverInfo(
    requestMessage: protocol_server_pb.Empty,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.ServerInfoData|null) => void
  ): UnaryResponse;
  serverInfo(
    requestMessage: protocol_server_pb.Empty,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.ServerInfoData|null) => void
  ): UnaryResponse;
  smartIQ(requestMessage: protocol_server_pb.ChannelConfig, metadata?: grpc.Metadata): ResponseStream<protocol_server_pb.IQData>;
  frequencyData(requestMessage: protocol_server_pb.FrequencyChannelConfig, metadata?: grpc.Metadata): ResponseStream<protocol_server_pb.IQData>;
  iQ(requestMessage: protocol_server_pb.ChannelConfig, metadata?: grpc.Metadata): ResponseStream<protocol_server_pb.IQData>;
  ping(
    requestMessage: protocol_server_pb.PingData,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.PingData|null) => void
  ): UnaryResponse;
  ping(
    requestMessage: protocol_server_pb.PingData,
    callback: (error: ServiceError|null, responseMessage: protocol_server_pb.PingData|null) => void
  ): UnaryResponse;
}

