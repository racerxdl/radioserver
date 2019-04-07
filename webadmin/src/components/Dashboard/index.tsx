import * as React from "react";
import FFTBoard from "../FFTBoard";
import {connect} from "react-redux";

type DashboardProps = {
  handleSetSectionName: (name: string) => void,
  samples?: number[] | void,
}

type DashboardState = {}

class Dashboard extends React.Component<DashboardProps, DashboardState> {
  constructor(props: DashboardProps) {
    super(props);
    if (this.props.handleSetSectionName) {
      this.props.handleSetSectionName(`Dashboard`)
    }
  }

  render() {
    return (
      <div className="App">
        <FFTBoard samples={this.props.samples} maxVal={-120} range={90} width={512} height={256}/>
      </div>
    );
  }

}

const mapStateToProps = (state: any) => {
  return ({
    samples: state.fftSamples.samples,
  });
};

export default connect(mapStateToProps)(Dashboard);
