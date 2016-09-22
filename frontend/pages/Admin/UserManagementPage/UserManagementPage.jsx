import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import Avatar from '../../../components/Avatar';
import componentStyles from './styles';
import Dropdown from '../../../components/forms/fields/Dropdown';
import entityGetter from '../../../redux/entityGetter';
import userActions from '../../../redux/nodes/entities/users/actions';

class UserManagementPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    users: PropTypes.arrayOf(PropTypes.object),
  };

  static userActionOptions = [
    { text: 'Disable Account', value: 'disable_account' },
    { text: 'Demote User', value: 'demote_user' },
    { text: 'Change Password', value: 'change_password' },
    { text: 'Require Password Reset', value: 'reset_password' },
    { text: 'Modify Details', value: 'modify_details' },
  ];

  componentWillMount () {
    const { dispatch, users } = this.props;
    const { load } = userActions;

    if (!users.length) dispatch(load());

    return false;
  }

  onUserActionSelect = (value) => {
    console.log(value);
    return false;
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
    const userActionOptions = UserManagementPage.userActionOptions;

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
            onSelect={this.onUserActionSelect}
          />
        </div>
      </div>
    );
  }

  render () {
    const {
      containerStyles,
      numUsersStyles,
      usersWrapperStyles,
    } = componentStyles;
    const { users } = this.props;

    return (
      <div style={containerStyles}>
        <p style={numUsersStyles}>Listing {users.length} users</p>
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

