import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import Avatar from '../../../../components/Avatar';
import componentStyles from './styles';
import Dropdown from '../../../../components/forms/fields/Dropdown';

class UserBlock extends Component {
  static propTypes = {
    onSelect: PropTypes.func,
    user: PropTypes.object,
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
  };

  onUserActionSelect = (formData) => {
    const { onSelect, user } = this.props;

    return onSelect(user, formData);
  }

  render () {
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
    const { user } = this.props;
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
    const userActionOptions = UserBlock.userActionOptions(user);

    return (
      <div style={userWrapperStyles}>
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
}

export default radium(UserBlock);
