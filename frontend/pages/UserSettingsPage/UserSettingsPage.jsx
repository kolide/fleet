import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { goBack } from 'react-router-redux';
import moment from 'moment';
import local from 'utilities/local';

import Avatar from 'components/Avatar';
import Button from 'components/buttons/Button';
import ChangeEmailForm from 'components/forms/ChangeEmailForm';
import ChangePasswordForm from 'components/forms/ChangePasswordForm';
import deepDifference from 'utilities/deep_difference';
import Icon from 'components/icons/Icon';
import InputField from 'components/forms/fields/InputField';
import { logoutUser, updateUser } from 'redux/nodes/auth/actions';
import Modal from 'components/modals/Modal';
import { renderFlash } from 'redux/nodes/notifications/actions';
import userActions from 'redux/nodes/entities/users/actions';
import userInterface from 'interfaces/user';
import UserSettingsForm from 'components/forms/UserSettingsForm';

const baseClass = 'user-settings';

export class UserSettingsPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func.isRequired,
    errors: PropTypes.shape({
      username: PropTypes.string,
      base: PropTypes.string,
    }),
    user: userInterface,
    userErrors: PropTypes.shape({
      base: PropTypes.string,
      new_password: PropTypes.string,
      old_password: PropTypes.string,
    }),
  };

  constructor (props) {
    super(props);

    this.state = {
      pendingEmail: undefined,
      showEmailModal: false,
      showPasswordModal: false,
      updatedUser: {},
    };
  }

  onCancel = (evt) => {
    evt.preventDefault();

    const { dispatch } = this.props;

    dispatch(goBack());

    return false;
  }

  onLogout = (evt) => {
    evt.preventDefault();

    const { dispatch } = this.props;

    dispatch(logoutUser());

    return false;
  }

  onShowModal = (evt) => {
    evt.preventDefault();

    this.setState({ showPasswordModal: true });

    return false;
  }

  onShowApiTokenModal = (evt) => {
    evt.preventDefault();

    this.setState({ showApiTokenModal: true });

    return false;
  }

  onToggleEmailModal = (updatedUser = {}) => {
    const { showEmailModal } = this.state;

    this.setState({
      showEmailModal: !showEmailModal,
      updatedUser,
    });

    return false;
  }

  onTogglePasswordModal = (evt) => {
    evt.preventDefault();

    const { showPasswordModal } = this.state;

    this.setState({ showPasswordModal: !showPasswordModal });

    return false;
  }

  onToggleApiTokenModal = (evt) => {
    evt.preventDefault();

    const { showApiTokenModal } = this.state;

    this.setState({ showApiTokenModal: !showApiTokenModal });

    return false;
  }

  handleSubmit = (formData) => {
    const { dispatch, user } = this.props;
    const updatedUser = deepDifference(formData, user);

    if (updatedUser.email && !updatedUser.password) {
      return this.onToggleEmailModal(updatedUser);
    }

    return dispatch(updateUser(user, updatedUser))
      .then(() => {
        if (updatedUser.email) {
          this.setState({ pendingEmail: updatedUser.email });
        }

        dispatch(renderFlash('success', 'Account updated!'));

        return true;
      })
      .catch(() => false);
  }

  handleSubmitPasswordForm = (formData) => {
    const { dispatch, user } = this.props;

    return dispatch(userActions.changePassword(user, formData))
      .then(() => {
        dispatch(renderFlash('success', 'Password changed successfully'));
        this.setState({ showPasswordModal: false });

        return false;
      });
  }

  renderEmailModal = () => {
    const { errors } = this.props;
    const { updatedUser, showEmailModal } = this.state;
    const { handleSubmit, onToggleEmailModal } = this;

    const emailSubmit = (formData) => {
      handleSubmit(formData)
        .then((r) => {
          return r ? onToggleEmailModal() : false;
        });
    };

    if (!showEmailModal) {
      return false;
    }

    return (
      <Modal
        title="To change your email you must supply your password"
        onExit={onToggleEmailModal}
      >
        <ChangeEmailForm
          formData={updatedUser}
          handleSubmit={emailSubmit}
          onCancel={onToggleEmailModal}
          serverErrors={errors}
        />
      </Modal>
    );
  }

  renderPasswordModal = () => {
    const { userErrors } = this.props;
    const { showPasswordModal } = this.state;
    const { handleSubmitPasswordForm, onTogglePasswordModal } = this;

    if (!showPasswordModal) {
      return false;
    }

    return (
      <Modal
        title="Change Password"
        onExit={onTogglePasswordModal}
      >
        <ChangePasswordForm
          handleSubmit={handleSubmitPasswordForm}
          onCancel={onTogglePasswordModal}
          serverErrors={userErrors}
        />
      </Modal>
    );
  }

  toggleSecret = (evt) => {
    const { revealSecret } = this.state;
    evt.preventDefault();

    this.setState({ revealSecret: !revealSecret });
    return false;
  }

  onCopySecret = (elementClass) => {
    return (evt) => {
      evt.preventDefault();

      const { dispatch } = this.props;

      if (copyText(elementClass)) {
        dispatch(renderFlash('success', 'Text copied to clipboard'));
      } else {
        this.setState({ revealSecret: true });
        dispatch(renderFlash('error', 'Text not copied. Please copy manually.'));
      }
    };
  }

  renderApiTokenModal = () => {
    const { showApiTokenModal, revealSecret } = this.state;
    const { onToggleApiTokenModal, onCopySecret, toggleSecret } = this;

    if (!showApiTokenModal) {
      return false;
    }

    return (
      <Modal
        title="Get API Token"
        onExit={onToggleApiTokenModal}
      >
          The following is your API Token:
          <a href="#revealSecret" onClick={toggleSecret} className={`${baseClass}__reveal-secret`}>{revealSecret ? 'Hide' : 'Reveal'} Secret</a>
              <div className={`${baseClass}__secret-wrapper`}>
                <InputField
                  disabled
                  inputWrapperClass={`${baseClass}__secret-input`}
                  name="osqueryd-secret"
                  type={revealSecret ? 'text' : 'password'}
                  value={local.getItem('auth_token')}
                />
                <Button variant="unstyled" className={`${baseClass}__secret-copy-icon`} onClick={onCopySecret(`.${baseClass}__secret-input`)}>
                  <Icon name="clipboard" />
                </Button>
              </div>

        <div className={`${baseClass}__button-wrap`}>
          <Button onClick={onToggleApiTokenModal} variant="success">
            Return To App
          </Button>
        </div>
      </Modal>
    );
  }

  render () {
    const {
      handleSubmit,
      onCancel,
      onLogout,
      onShowModal,
      onShowApiTokenModal,
      renderEmailModal,
      renderPasswordModal,
      renderApiTokenModal,
    } = this;
    const { errors, user } = this.props;
    const { pendingEmail } = this.state;

    if (!user) {
      return false;
    }

    const { admin, updated_at: updatedAt, sso_enabled: ssoEnabled } = user;
    const roleText = admin ? 'ADMIN' : 'USER';
    const lastUpdatedAt = moment(updatedAt).fromNow();

    return (
      <div className={baseClass}>
        <div className={`${baseClass}__manage body-wrap`}>
          <h1>Manage User Settings</h1>
          <UserSettingsForm
            formData={user}
            handleSubmit={handleSubmit}
            onCancel={onCancel}
            pendingEmail={pendingEmail}
            serverErrors={errors}
          />
        </div>
        <div className={`${baseClass}__additional body-wrap`}>
          <h1>Additional Info</h1>

          <div className={`${baseClass}__change-avatar`}>
            <Avatar user={user} className={`${baseClass}__avatar`} />
            <a href="http://en.gravatar.com/emails/">Change Photo at Gravatar</a>
          </div>

          <Button onClick={onShowApiTokenModal} variant="brand" className={`${baseClass}__button`}>
            GET API TOKEN
          </Button>

          <div className={`${baseClass}__more-info-detail`}>
            <Icon name="username" />
            <strong>Role</strong> - {roleText}
          </div>
          <div className={`${baseClass}__more-info-detail`}>
            <Icon name="lock-big" />
            <strong>Password</strong>
          </div>
          <Button onClick={onShowModal} variant="brand" disabled={ssoEnabled} className={`${baseClass}__button`}>
            CHANGE PASSWORD
          </Button>
          <p className={`${baseClass}__last-updated`}>Last changed: {lastUpdatedAt}</p>
          <Button onClick={onLogout} variant="alert" className={`${baseClass}__button`}>
            LOGOUT
          </Button>
        </div>
        {renderEmailModal()}
        {renderPasswordModal()}
        {renderApiTokenModal()}
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { errors, user } = state.auth;
  const { errors: userErrors } = state.entities.users;

  return { errors, user, userErrors };
};

export default connect(mapStateToProps)(UserSettingsPage);
