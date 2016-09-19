import React, { Component, PropTypes } from 'react';
import componentStyles from './styles';
import GradientButton from '../../buttons/GradientButton';

class LogoutForm extends Component {
  static propTypes = {
    onSubmit: PropTypes.func,
    user: PropTypes.object,
  };

  render () {
    const { onSubmit, user } = this.props;
    const { containerStyles, formStyles, submitButtonStyles } = componentStyles;
    const { gravatarURL } = user;

    return (
      <form onSubmit={onSubmit} style={formStyles}>
        <div style={containerStyles}>
          <img alt="Avatar" src={gravatarURL} />
          <h1>{user.username}</h1>
        </div>
        <GradientButton
          onClick={onSubmit}
          style={submitButtonStyles}
          text="Logout"
          type="submit"
        />
      </form>
    );
  }
}

export default LogoutForm;
