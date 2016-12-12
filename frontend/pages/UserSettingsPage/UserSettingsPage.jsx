import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { goBack } from 'react-router-redux';

import { renderFlash } from 'redux/nodes/notifications/actions';
import userActions from 'redux/nodes/entities/users/actions';
import userInterface from 'interfaces/user';
import UserSettingsForm from 'components/forms/UserSettingsForm';

class UserSettingsPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func.isRequired,
    user: userInterface,
  };

  onCancel = (evt) => {
    evt.preventDefault();

    const { dispatch } = this.props;

    dispatch(goBack());

    return false;
  }

  handleSubmit = (formData) => {
    const { dispatch, user } = this.props;
    const { update } = userActions;

    return dispatch(update(user, formData))
      .then(() => {
        return dispatch(renderFlash('success', 'Account updated!'));
      });
  }

  render () {
    const { handleSubmit, onCancel } = this;
    const { user } = this.props;

    if (!user) {
      return false;
    }

    return (
      <div>
        <div className="body-wrap">
          <h1>Manage User Settings</h1>
          <UserSettingsForm formData={user} handleSubmit={handleSubmit} onCancel={onCancel} />
        </div>
        <div className="body-wrap">
          <h3>Something else</h3>
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { user } = state.auth;

  return { user };
};

export default connect(mapStateToProps)(UserSettingsPage);
