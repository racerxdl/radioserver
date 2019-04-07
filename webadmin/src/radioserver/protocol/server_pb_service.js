/* eslint-disable */
// package: protocol
// file: protocol/server.proto

var protocol_server_pb = require("../protocol/server_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var RadioServer = (function () {
  function RadioServer() {}
  RadioServer.serviceName = "protocol.RadioServer";
  return RadioServer;
}());

RadioServer.Me = {
  methodName: "Me",
  service: RadioServer,
  requestStream: false,
  responseStream: false,
  requestType: protocol_server_pb.LoginData,
  responseType: protocol_server_pb.MeData
};

RadioServer.Hello = {
  methodName: "Hello",
  service: RadioServer,
  requestStream: false,
  responseStream: false,
  requestType: protocol_server_pb.HelloData,
  responseType: protocol_server_pb.HelloReturn
};

RadioServer.Bye = {
  methodName: "Bye",
  service: RadioServer,
  requestStream: false,
  responseStream: false,
  requestType: protocol_server_pb.LoginData,
  responseType: protocol_server_pb.ByeReturn
};

RadioServer.ServerInfo = {
  methodName: "ServerInfo",
  service: RadioServer,
  requestStream: false,
  responseStream: false,
  requestType: protocol_server_pb.Empty,
  responseType: protocol_server_pb.ServerInfoData
};

RadioServer.SmartIQ = {
  methodName: "SmartIQ",
  service: RadioServer,
  requestStream: false,
  responseStream: true,
  requestType: protocol_server_pb.ChannelConfig,
  responseType: protocol_server_pb.IQData
};

RadioServer.FrequencyData = {
  methodName: "FrequencyData",
  service: RadioServer,
  requestStream: false,
  responseStream: true,
  requestType: protocol_server_pb.FrequencyChannelConfig,
  responseType: protocol_server_pb.IQData
};

RadioServer.IQ = {
  methodName: "IQ",
  service: RadioServer,
  requestStream: false,
  responseStream: true,
  requestType: protocol_server_pb.ChannelConfig,
  responseType: protocol_server_pb.IQData
};

RadioServer.Ping = {
  methodName: "Ping",
  service: RadioServer,
  requestStream: false,
  responseStream: false,
  requestType: protocol_server_pb.PingData,
  responseType: protocol_server_pb.PingData
};

exports.RadioServer = RadioServer;

function RadioServerClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

RadioServerClient.prototype.me = function me(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RadioServer.Me, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

RadioServerClient.prototype.hello = function hello(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RadioServer.Hello, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

RadioServerClient.prototype.bye = function bye(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RadioServer.Bye, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

RadioServerClient.prototype.serverInfo = function serverInfo(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RadioServer.ServerInfo, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

RadioServerClient.prototype.smartIQ = function smartIQ(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(RadioServer.SmartIQ, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.end.forEach(function (handler) {
        handler();
      });
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

RadioServerClient.prototype.frequencyData = function frequencyData(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(RadioServer.FrequencyData, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.end.forEach(function (handler) {
        handler();
      });
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

RadioServerClient.prototype.iQ = function iQ(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(RadioServer.IQ, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.end.forEach(function (handler) {
        handler();
      });
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

RadioServerClient.prototype.ping = function ping(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RadioServer.Ping, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.RadioServerClient = RadioServerClient;

