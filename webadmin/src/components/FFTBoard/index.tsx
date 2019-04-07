import {Component, createRef, default as React} from "react";
import PropTypes from 'prop-types';


type FFTCanvasProps = {
  samples: number[] | void
  maxVal: number
  range: number
  width: number
  height: number
}

class FFTCanvas extends Component<FFTCanvasProps> {
  canvasRef: any;

  constructor(props: FFTCanvasProps) {
    super(props);
    this.canvasRef = createRef();
  }

  componentDidUpdate() {
    const canvas = this.canvasRef.current;
    const ctx = canvas.getContext('2d');
    const width = canvas.width;
    const height = canvas.height;

    const data = this.props.samples;
    ctx.fillStyle = 'rgb(0, 0, 0)';
    ctx.fillRect(0, 0, width, height);

    if (!data || data.length === 0) {
      return;
    }

    const bufferLength = data.length;
    const widthScale = width / bufferLength;

    ctx.moveTo(0, 0);
    ctx.beginPath();
    ctx.lineWidth = 2;
    ctx.strokeStyle = '#AAAAAA';
    ctx.moveTo(0, 0);
    for (let i = 0; i < bufferLength; i++) {
      const v = height - (((data[i] - this.props.maxVal) / (this.props.range)) * height);
      ctx.lineTo(i * widthScale, v);
    }
    ctx.stroke();
  }

  render() {
    return (
      <canvas width={this.props.width} height={this.props.height} ref={this.canvasRef}/>
    )
  }
}

type FFTBoardProps = FFTCanvasProps & {
  // TODO
}

class FFTBoard extends Component<FFTBoardProps> {
  constructor(props: FFTBoardProps) {
    super(props);
  }

  render() {
    return (
      <FFTCanvas {...this.props} />
    )
  }

  static propTypes = {
    samples: PropTypes.arrayOf(PropTypes.number),
  }
}


export default FFTBoard;
