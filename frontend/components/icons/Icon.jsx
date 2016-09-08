import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import KolideLoginBackground from './svg/KolideLoginBackground';
import KolideText from './svg/KolideText';
import User from './svg/User';

class Icon extends Component {
  static propTypes = {
    alt: PropTypes.string,
    name: PropTypes.string,
    style: PropTypes.object,
    variant: PropTypes.string,
  };

  static iconNames = {
    kolideLoginBackground: KolideLoginBackground,
    kolideText: KolideText,
    user: User,
  };

  render () {
    const { alt, name, style, variant } = this.props;
    const IconComponent = Icon.iconNames[name];

    return <IconComponent alt={alt} style={style} variant={variant} />;
  }
}

export default radium(Icon);
