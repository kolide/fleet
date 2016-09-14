import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import componentStyles from './styles';
import { forgotPasswordAction } from '../../redux/nodes/components/ForgotPasswordPage/actions';
import ForgotPasswordForm from '../../components/forms/ForgotPasswordForm';

class ForgotPasswordPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
  };

  static defaultProps = {
    dispatch: noop,
  };

  onSubmit = (formData) => {
    const { dispatch } = this.props;

    return dispatch(forgotPasswordAction(formData));
  }

  render () {
    const {
      containerStyles,
      forgotPasswordStyles,
      headerStyles,
      smallWhiteTabStyles,
      textStyles,
      whiteTabStyles,
    } = componentStyles;

    return (
      <div style={containerStyles}>
        <div style={smallWhiteTabStyles} />
        <div style={whiteTabStyles} />
        <div style={forgotPasswordStyles}>
          <p style={headerStyles}>Forgot Password</p>
          <p style={textStyles}>If you’ve forgotten your password enter your email below and we will email you a link so that you can reset your password.</p>
          <ForgotPasswordForm onSubmit={this.onSubmit} />
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  return state.components.ForgotPasswordPage;
};

export default connect(mapStateToProps)(ForgotPasswordPage);
