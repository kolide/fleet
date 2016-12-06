import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

import ConfirmInviteForm from 'components/forms/ConfirmInviteForm';
import AuthenticationFormWrapper from 'components/AuthenticationFormWrapper';

const baseClass = 'confirm-invite-page';

class ConfirmInvitePage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    inviteFormData: PropTypes.shape({
      invite_token: PropTypes.string.isRequired,
    }).isRequired,
  };

  onSubmit = (formData) => {
    console.log(formData);

    return false;
  }

  render () {
    const { inviteFormData } = this.props;
    const { onSubmit } = this;

    return (
      <AuthenticationFormWrapper>
        <div className={`${baseClass}__lead-wrapper`}>
          <p className={`${baseClass}__lead-text`}>
            Welcome to the party!
          </p>
          <p className={`${baseClass}__sub-lead-text`}>
            Please take a moment to fill out the following information before we take you into <b>Kolide</b>
          </p>
        </div>
        <div className={`${baseClass}__form-section-wrapper`}>
          <div className={`${baseClass}__form-section-description`}>
            <h2>SET USERNAME & PASSWORD</h2>
            <p>Password must include 7 characters, at least 1 number (eg. 0-9), and at least 1 symbol (eg. ^&*#)</p>
          </div>
          <ConfirmInviteForm
            className={`${baseClass}__form`}
            formData={inviteFormData}
            handleSubmit={onSubmit}
          />
        </div>
      </AuthenticationFormWrapper>
    );
  }
}

const mapStateToProps = (state, { params }) => {
  const { invite_token: inviteToken } = params;
  const inviteFormData = { invite_token: inviteToken };

  return { inviteFormData };
};

export default connect(mapStateToProps)(ConfirmInvitePage);
