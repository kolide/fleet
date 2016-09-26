import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import Avatar from '../../../components/Avatar';
import GradientButton from '../../../components/buttons/GradientButton';
import componentStyles from './styles';
import Dropdown from '../../../components/forms/fields/Dropdown';
import entityGetter from '../../../redux/entityGetter';
import userActions from '../../../redux/nodes/entities/users/actions';

class UserManagementPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    users: PropTypes.arrayOf(PropTypes.object),
  };

  static userActionOptions = (user) => {
    const userEnableAction = user.enabled
      ? { text: 'Disable Account', value: 'disable_account' }
      : { text: 'Enable Account', value: 'enable_account' };
    const userPromotionAction = user.admin
      ? { text: 'Demote User', value: 'demote_user' }
      : { text: 'Promote User', value: 'promote_user' };

    return [
      { text: 'Actions...', value: '' },
      userEnableAction,
      userPromotionAction,
      { text: 'Require Password Reset', value: 'reset_password' },
      { text: 'Modify Details', value: 'modify_details' },
    ];
  }

  componentWillMount () {
    const { dispatch, users } = this.props;
    const { load } = userActions;

    if (!users.length) dispatch(load());

    return false;
  }

  onUserActionSelect = (user) => {
    return (formData) => {
      const { dispatch } = this.props;

      if (formData.user_actions) {
        switch (formData.user_actions) {
          case 'demote_user':
            return dispatch(userActions.update(user, { admin: false }));
          case 'disable_account':
            return dispatch(userActions.update(user, { enabled: false }));
          case 'enable_account':
            return dispatch(userActions.update(user, { enabled: true }));
          case 'promote_user':
            return dispatch(userActions.update(user, { admin: true }));
          case 'reset_password':
            return dispatch(userActions.update(user, { force_password_reset: true }));
          default:
            return false;
        }
      }

      return false;
    };
  }

  renderUserBlock = (user) => {
    const {
      avatarStyles,
      nameStyles,
      userDetailsStyles,
      userEmailStyles,
      userHeaderStyles,
      userLabelStyles,
      usernameStyles,
      userPositionStyles,
      userStatusStyles,
      userStatusWrapperStyles,
      userWrapperStyles,
    } = componentStyles;
    const {
      admin,
      email,
      enabled,
      name,
      position,
      username,
    } = user;
    const userLabel = admin ? 'Admin' : 'User';
    const activeLabel = enabled ? 'Active' : 'Disabled';
    const userActionOptions = UserManagementPage.userActionOptions(user);

    return (
      <div key={email} style={userWrapperStyles}>
        <div style={userHeaderStyles}>
          <span style={nameStyles}>{name}</span>
        </div>
        <div style={userDetailsStyles}>
          <Avatar user={user} style={avatarStyles} />
          <div style={userStatusWrapperStyles}>
            <span style={userLabelStyles}>{userLabel}</span>
            <span style={userStatusStyles(enabled)}>{activeLabel}</span>
            <div style={{ clear: 'both' }} />
          </div>
          <p style={usernameStyles}>{username}</p>
          <p style={userPositionStyles}>{position}</p>
          <p style={userEmailStyles}>{email}</p>
          <Dropdown
            fieldName="user_actions"
            options={userActionOptions}
            initialOption={{ text: 'Actions...' }}
            onSelect={this.onUserActionSelect(user)}
          />
        </div>
      </div>
    );
  }

  render () {
    const {
      addUserButtonStyles,
      addUserWrapperStyles,
      containerStyles,
      numUsersStyles,
      usersWrapperStyles,
    } = componentStyles;
    const { users } = this.props;

    return (
      <div style={containerStyles}>
        <span style={numUsersStyles}>Listing {users.length} users</span>
        <div style={addUserWrapperStyles}>
          <GradientButton
            style={addUserButtonStyles}
            text="Add User"
          />
        </div>
        <div style={usersWrapperStyles}>
          {users.map(user => {
            return this.renderUserBlock(user);
          })}
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: users } = entityGetter(state).get('users');

  return { users };
};

export default connect(mapStateToProps)(UserManagementPage);

