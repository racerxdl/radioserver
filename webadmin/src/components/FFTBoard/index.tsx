import {Component, createRef, default as React} from "react";
import PropTypes from 'prop-types';
import withStyles from "@material-ui/core/styles/withStyles";
import createStyles from "@material-ui/core/styles/createStyles";
import {WithStyles} from "@material-ui/core";


const margin = 50;
const marginTop = 28;
const spacing = 60;

const styles = () => createStyles({
  fftCanvas: {
    borderRadius: '8px',
  },
});

interface FFTCanvasProps extends WithStyles<'root'> {
  samples: number[] | void
  maxVal: number
  range: number
  width: number
  height: number
  fps: number
  centerFrequency: number,
  sampleRate: number,
  classes: any,
}

function CalcDiv(size: number): number {
  return (size / spacing) >> 0;
}

class FFTCanvas extends Component<FFTCanvasProps> {
  canvasRef: any;
  classes: any;

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

    // region Pre Compute Helper Variables
    const fftScreenWidth = width - margin * 1.5;
    const fftScreenHeight = height - marginTop * 2.2;
    const fftScreenX = margin;
    const fftScreenY = marginTop;

    const bufferLength = data.length;
    const bufferScale = fftScreenWidth / bufferLength;
    const max = this.props.maxVal;
    const dbPerPixel = this.props.range / fftScreenHeight;

    const hDivs = CalcDiv(fftScreenHeight) * 2;
    const vDivs = CalcDiv(fftScreenWidth);
    const hDivDelta = fftScreenHeight / hDivs;
    const vDivDelta = fftScreenWidth / vDivs;
    const delta = this.props.sampleRate / fftScreenWidth;
    const startFreq = this.props.centerFrequency - (this.props.sampleRate / 2);

    const minVal = this.props.maxVal - this.props.range;
    const dbToPixel = (dB: number): number => {
      return fftScreenY + fftScreenHeight - ((dB - minVal) / this.props.range) * fftScreenHeight;
    };
    // endregion
    // region Draw Frequency Labels
    ctx.beginPath();
    ctx.lineWidth = 1;
    ctx.strokeStyle = '#444444';
    for (let i = 0; i < fftScreenWidth + 1; i += vDivDelta) {
      ctx.moveTo(margin + i, marginTop);
      ctx.lineTo(margin + i, marginTop + fftScreenHeight + 5);
      const freq = (startFreq + i * delta) / 1e6;
      const freqStr = freq.toLocaleString();
      ctx.save();
      ctx.fillStyle = '#FFFFFF';
      const freqX = margin + i - ctx.measureText(freqStr).width / 2;
      ctx.fillText(freqStr, freqX, marginTop + fftScreenHeight + 25);
      ctx.restore();
    }

    ctx.save();
    ctx.fillStyle = '#FFFFFF';
    const MHzText = 'MHz';
    ctx.fillText(MHzText, fftScreenX + fftScreenWidth - 5, fftScreenY + fftScreenHeight + 15);
    ctx.restore();
    // endregion
    // region Draw dB Label

    for (let i = 0; i < fftScreenHeight + 1; i += hDivDelta) {
      ctx.moveTo(fftScreenX - 10, fftScreenY + i);
      ctx.lineTo(fftScreenWidth + margin, marginTop + i);
      const dbLvl = max - (i * dbPerPixel);
      const dbLvlStr = dbLvl.toFixed(0);
      ctx.save();
      ctx.fillStyle = '#FFFFFF';
      const dbLvlX = fftScreenX - ctx.measureText(dbLvlStr).width - 15;
      ctx.fillText(dbLvlStr, dbLvlX, marginTop + i + 5);
      ctx.restore();
    }
    ctx.stroke();
    ctx.closePath();
    // endregion
    // region Draw FFT
    const firstV = dbToPixel(data[0]);
    ctx.beginPath();
    ctx.lineWidth = 2;
    ctx.strokeStyle = '#AAAAAA';
    ctx.moveTo(fftScreenX, firstV);
    let avg = 0;
    for (let i = 0; i < bufferLength; i++) {
      const v = dbToPixel(data[i]);
      if (v > fftScreenHeight + fftScreenY || v < fftScreenY) {
        ctx.moveTo(fftScreenX + i * bufferScale, v);
      } else {
        ctx.lineTo(fftScreenX + i * bufferScale, v);
      }
      avg += data[i];
    }
    avg /= bufferLength;
    ctx.stroke();
    // endregion
    // region Stats Label
    const avgH = dbToPixel(avg);

    ctx.lineWidth = 1;
    ctx.strokeStyle = '#AA0000';
    ctx.beginPath();
    ctx.moveTo(fftScreenX, avgH);
    ctx.lineTo(fftScreenX + fftScreenWidth, avgH);
    ctx.stroke();

    ctx.fillStyle = 'rgb(255,255,255)';
    ctx.fillText(`${this.props.fps} FPS | Avg: ${Math.round(avg)} dB`, 8, 15);
    // endregion
  }

  render() {
    return (
      <canvas width={this.props.width} height={this.props.height} ref={this.canvasRef}
              className={this.props.classes.fftCanvas}/>
    )
  }
}

interface FFTBoardProps extends FFTCanvasProps {
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


// @ts-ignore
export default withStyles(styles)(FFTBoard);
