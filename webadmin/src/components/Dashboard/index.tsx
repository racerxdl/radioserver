import * as React from "react";
import FFTBoard from "../FFTBoard";
import {connect} from "react-redux";

type DashboardProps = {
  handleSetSectionName: (name: string) => void,
  samples?: number[] | void,
  sampleRate: number,
  centerFrequency: number,
}

type DashboardState = {
  samples?: number[] | void,
  delta: number,
  lastUpdate: number,
  fps: number,
}

class Dashboard extends React.Component<DashboardProps, DashboardState> {
  rafId?: number;
  state = {
    samples: [],
    delta: 0,
    lastUpdate: Date.now(),
    fps: 0,
  };

  constructor(props: DashboardProps) {
    super(props);
    if (this.props.handleSetSectionName) {
      this.props.handleSetSectionName(`Dashboard`)
    }
  }

  componentDidMount() {
    this.rafId = requestAnimationFrame(this.tick);
  }

  tick = () => {
    const samples = this.props.samples || [];
    let changed = false;
    for (let i = 0; i < samples.length; i++) {
      if (samples[i] !== this.state.samples[i]) {
        changed = true;
        break;
      }
    }

    if (changed && samples.length) {
      const delta = Date.now() - this.state.lastUpdate;
      const lastUpdate = Date.now();
      this.setState({
        samples: this.props.samples,
        lastUpdate,
        delta,
        fps: Math.round(1000 / delta),
      });
    }

    this.rafId = requestAnimationFrame(this.tick);
  };

  render() {
    return (
      <div className="App">
        <FFTBoard
          samples={this.state.samples}
          maxVal={-70}
          range={40}
          width={512}
          height={256}
          fps={this.state.fps}
          centerFrequency={this.props.centerFrequency}
          sampleRate={this.props.sampleRate}
        />
      </div>
    );
  }

}

const mapStateToProps = (state: any) => {
  return ({
    samples: state.fftSamples.samples,
    sampleRate: state.fftSamples.sampleRate,
    centerFrequency: state.fftSamples.centerFrequency,
  });
};

export default connect(mapStateToProps)(Dashboard);
