import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import AuthenticationFormWrapper from '../../components/AuthenticationFormWrapper';
import debounce from '../../utilities/debounce';
import LogoutForm from '../../components/forms/LogoutForm';
import { logoutUser } from '../../redux/nodes/auth/actions';

export class LogoutPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    user: PropTypes.object,
  };

  static defaultProps = {
    dispatch: noop,
  };

  onSubmit = debounce(() => {
    const { dispatch } = this.props;

    return dispatch(logoutUser());
  })

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
