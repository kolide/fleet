import React, { Component } from 'react';
import componentStyles from './styles';
import { loadBackground, resizeBackground } from '../../utilities/backgroundImage';

export class LoginPage extends Component {
  componentWillMount () {
    const { window } = global;

    loadBackground();
    window.onresize = resizeBackground;
  }

  render () {
    const { containerStyles } = componentStyles;

    return (
      <div style={containerStyles}>
        <h1>Login page</h1>
      </div>
    );
  }
}

export default LoginPage;
