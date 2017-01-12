import React, { Component, PropTypes } from 'react';

const baseClass = 'kolide-timer';
let offset = null;
let interval = null;

class Timer extends Component {
  static propTypes = {
    running: PropTypes.bool,
  }

  constructor (props) {
    super(props);

    this.state = {
      totalSeconds: 0,
      currrentTimer: '',
      running: props.running,
    };
  }

  componentWillReceiveProps ({ running }) {
    this.setState({ running });

    if (running) {
      this.play();
    } else {
      this.pause();
    }
  }

  play = () => {
    if (!interval) {
      offset = Date.now();
      interval = setInterval(this.update, 1000);
    }
  }

  pause = () => {
    if (interval) {
      clearInterval(interval);
      interval = null;
    }
  }

  update = () => {
    let { totalSeconds } = this.state;

    totalSeconds += this.calculateOffset();
    this.convertSeconds();

    this.setState({
      totalSeconds,
    });
  }

  convertSeconds = () => {
    const { totalSeconds } = this.state;
    const currentSeconds = totalSeconds / 1000;

    const hours = Math.floor(currentSeconds / 3600);
    const minutes = Math.floor((currentSeconds - (hours * 3600)) / 60);
    let seconds = currentSeconds - (hours * 3600) - (minutes * 60);
    seconds = Math.round((seconds * 100) / 100);

    const resultHrs = hours < 10 ? `0${hours}` : hours;
    const resultMins = minutes < 10 ? `0${minutes}` : minutes;
    const resultSecs = seconds < 10 ? `0${seconds}` : seconds;

    this.setState({ currrentTimer: `${resultHrs}:${resultMins}:${resultSecs}` });
  }

  calculateOffset = () => {
    const now = Date.now();
    const newOffset = now - offset;
    offset = now;

    return newOffset;
  }

  render () {
    const { currrentTimer } = this.state;

    return (
      <span className={baseClass}>{currrentTimer}</span>
    );
  }
}

export default Timer;
