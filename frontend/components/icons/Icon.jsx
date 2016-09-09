import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import Check from './svg/Check';
import KolideLoginBackground from './svg/KolideLoginBackground';
import KolideLogo from './svg/KolideLogo';
import KolideText from './svg/KolideText';
import Lock from './svg/Lock';
import User from './svg/User';

class Icon extends Component {
  static propTypes = {
    alt: PropTypes.string,
    name: PropTypes.string.isRequired,
    style: PropTypes.object,
    variant: PropTypes.string,
  };

  static iconNames = {
    check: Check,
    kolideLoginBackground: KolideLoginBackground,
    kolideLogo: KolideLogo,
    kolideText: KolideText,
    lock: Lock,
    user: User,
  };

  render () {
    const IconComponent = Icon.iconNames[this.props.name];

    return <IconComponent {...this.props} />;
  }
}

export default radium(Icon);
