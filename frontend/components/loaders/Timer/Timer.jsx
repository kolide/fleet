import React, { Component, PropTypes } from 'react';

import convertSeconds from './convertSeconds';

const baseClass = 'kolide-timer';
let offset = null;
let interval = null;

class Timer extends Component {
  static propTypes = {
    running: PropTypes.bool,
    reset: PropTypes.bool,
  }

  constructor (props) {
    super(props);

    this.state = {
      totalSeconds: 0,
      currrentTimer: '',
    };
  }

  componentWillReceiveProps ({ running, reset }) {
    if (reset) {
      this.reset();
    }

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

  reset = () => {
    this.setState({
      totalSeconds: 0,
      currrentTimer: convertSeconds(0),
    });
  }

  update = () => {
    let { totalSeconds } = this.state;

    totalSeconds += this.calculateOffset();

    this.setState({
      currrentTimer: convertSeconds(totalSeconds),
      totalSeconds,
    });
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
