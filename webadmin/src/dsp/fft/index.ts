const supportedSizes = [256, 512, 1024, 2048, 4096];

type FFTCallback = (data: Float32Array) => void;

// type FFTProcessor = (inputBuffer: Float32Array, outputBuffer: Float32Array | FFTCallback, callback?: FFTCallback) => void

class FFT {
  cb?: FFTCallback;
  analyser: AnalyserNode;
  sourceNode: ScriptProcessorNode;
  buffers: Float32Array[];

  constructor(size: number) {
    if (supportedSizes.indexOf(size) === -1) {
      throw new Error("Invalid buffer size")
    }

    const audioContext = new AudioContext();
    this.analyser = audioContext.createAnalyser();
    this.analyser.fftSize = size;
    this.sourceNode = audioContext.createScriptProcessor(size, 0, 1);
    this.buffers = [];
    this.sourceNode.onaudioprocess = this.processAudio;

    this.sourceNode.connect(this.analyser);
  }

  processAudio = (event: AudioProcessingEvent) => {
    if (this.buffers.length) {
      const buff = this.buffers.splice(0, 1)[0];
      event.outputBuffer.getChannelData(0).set(buff);
    }
  };

  PutSamples = (samples: Float32Array) => {
    this.buffers.push(samples);
  };

  getFFT(): Uint8Array {
    const data = new Uint8Array(this.analyser.frequencyBinCount);
    this.analyser.getByteFrequencyData(data);
    return data;
  }
}

export default FFT;
