import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import AuthenticationFormWrapper from '../../components/AuthenticationFormWrapper';
import LogoutForm from '../../components/forms/LogoutForm';

export class LogoutPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    user: PropTypes.object,
  };

  static defaultProps = {
    dispatch: noop,
  };

  onSubmit = (formData) => {
    console.log('formData', formData);
  }

  render () {
    const { user } = this.props;
    const { onSubmit } = this;

    if (!user) return false;

    return (
      <AuthenticationFormWrapper>
        <LogoutForm onSubmit={onSubmit} user={user} />
      </AuthenticationFormWrapper>
    );
  }
}

const mapStateToProps = (state) => {
  const { user } = state.auth;

  return { user };
};
export default connect(mapStateToProps)(LogoutPage);
