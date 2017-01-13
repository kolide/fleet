import React, { Component, PropTypes } from 'react';

import convertSeconds from './convertSeconds';

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
    };
  }

  componentWillReceiveProps ({ running }) {
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
