import logo from '../../logo.svg';
import * as React from "react";

type MainAppProps = {
  handleSetSectionName: (name: string) => void,
}

class MainApp extends React.Component<MainAppProps> {
  constructor(props: MainAppProps) {
    super(props);

    if (this.props.handleSetSectionName) {
      this.props.handleSetSectionName(`Main Page`)
    }
  }

  render() {
    return (
      <div className="App">
        <img src={logo} className="App-logo" alt="logo"/>
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </div>
    );
  }
}

export default MainApp;

