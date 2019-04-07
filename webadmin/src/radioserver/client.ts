import {grpc} from '@improbable-eng/grpc-web';
import {RadioServer} from './protocol/server_pb_service';
import {
  ChannelConfig,
  FrequencyChannelConfig,
  HelloData,
  HelloReturn,
  IQData,
  LoginData,
  PingData,
  StatusType
} from "./protocol/server_pb";

type IQCallback = (samples: number[] | null, error: string | null) => void;

class RadioClient {
  url: string;
  smartIQClient: grpc.Client<ChannelConfig, IQData>;
  fftClient: grpc.Client<FrequencyChannelConfig, IQData>;
  token: string | null;
  onSmartIQ?: IQCallback;
  onFFT?: IQCallback;

  constructor(url: string) {
    this.url = url;
    this.smartIQClient = grpc.client(RadioServer.SmartIQ, {
      host: url,
    });
    this.fftClient = grpc.client(RadioServer.FrequencyData, {
      host: url,
    });

    this.smartIQClient.onMessage((message: IQData) => {
      if (this.onSmartIQ) {
        if (message.getStatus() === StatusType.OK) {
          this.onSmartIQ(message.getSamplesList(), null)
        } else {
          this.onSmartIQ(null, message.getError());
        }
      }
    });
    this.fftClient.onMessage((message: IQData) => {
      if (this.onFFT) {
        if (message.getStatus() === StatusType.OK) {
          this.onFFT(message.getSamplesList(), null)
        } else {
          this.onFFT(null, message.getError());
        }
      }
    });
    this.token = null;
  }

  async Ping() {
    return new Promise((resolve, reject) => {
      const req = new PingData();
      req.setTimestamp(Date.now());
      grpc.unary(RadioServer.Ping, {
        request: req,
        host: this.url,
        onEnd: (res) => {
          const {status, statusMessage, message} = res;
          if (status === grpc.Code.OK && message) {
            resolve(message.toObject());
          } else {
            reject(statusMessage);
          }
        },
      })
    });
  }

  SetOnSmartIQ(cb: (samples: number[] | null, error: string | null) => void) {
    this.onSmartIQ = cb
  }

  SetOnFFT(cb: (samples: number[] | null, error: string | null) => void) {
    this.onFFT = cb
  }

  async Login() {
    this.token = await this._login();
  }

  async Logout() {
    return new Promise((resolve, reject) => {
      if (this.token) {
        const loginData = new LoginData();
        loginData.setToken(this.token);
        grpc.unary(RadioServer.Bye, {
          request: loginData,
          host: this.url,
          onEnd: (res) => {
            if (res.status === grpc.Code.OK) {
              resolve();
            } else {
              reject(`Error: ${res.statusMessage}`);
            }
          },
        });
        return;
      }
      reject('Not logged in');
    });
  }

  async _login(): Promise<string> {
    return new Promise((resolve, reject) => {
      const req = new HelloData();
      req.setApplication("RadioServer WebAdmin");
      req.setName("None");
      grpc.unary(RadioServer.Hello, {
        request: req,
        host: this.url,
        onEnd: (res) => {
          const {status, statusMessage, message} = res;
          if (status === grpc.Code.OK && message) {
            const login = (message as HelloReturn).getLogin();
            if (login) {
              return resolve(login.getToken());
            }
            reject('login came null');
          } else {
            reject(statusMessage);
          }
        }
      })
    });
  }

  StartSmartIQ = async (centerFrequency: number, decimationStage: number) => {
    if (this.token === null) {
      throw new Error('not logged in');
    }
    const loginData = new LoginData();
    loginData.setToken(this.token);
    this.smartIQClient.start();
    const cc = new ChannelConfig();
    cc.setCenterfrequency(centerFrequency);
    cc.setDecimationstage(decimationStage);
    cc.setLogininfo(loginData);

    this.smartIQClient.send(cc);
  };

  StopSmartIQ = () => {
    this.smartIQClient.close();
  };

  StartFFT = async (centerFrequency: number, decimationStage: number, length: number) => {
    if (this.token === null) {
      throw new Error('not logged in');
    }
    const loginData = new LoginData();
    loginData.setToken(this.token);
    this.fftClient.start();
    const cc = new FrequencyChannelConfig();
    cc.setCenterfrequency(centerFrequency);
    cc.setDecimationstage(decimationStage);
    cc.setLogininfo(loginData);
    cc.setLength(length);

    this.fftClient.send(cc);
  };

  StopFFT = () => {
    this.fftClient.close();
  };
}

export default RadioClient;
