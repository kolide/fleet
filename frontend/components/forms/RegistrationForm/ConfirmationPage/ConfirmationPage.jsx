import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';

import Button from 'components/buttons/Button';
import formDataInterface from 'interfaces/registration_form_data';

const baseClass = 'confirm-user-reg';

class ConfirmationPage extends Component {
  static propTypes = {
    className: PropTypes.string,
    formData: formDataInterface,
    handleSubmit: PropTypes.func,
  };

  onSubmit = (evt) => {
    evt.preventDefault();

    const { handleSubmit } = this.props;

    return handleSubmit();
  }

  render () {
    const {
      className,
      formData: {
        email,
        full_name: fullName,
        kolide_server_url: kolideWebAddress,
        org_name: orgName,
        username,
      },
    } = this.props;
    const { onSubmit } = this;

    const confirmRegClasses = classnames(className, baseClass);
    const confirmIconClasses = classnames('kolidecon', 'kolidecon-success-check', `${baseClass}__icon`);

    return (
      <div className={confirmRegClasses}>
        <div className={`${baseClass}__wrapper`}>
          <i className={confirmIconClasses} />
          <table className={`${baseClass}__table`}>
            <caption>Administrator Configuration</caption>
            <tbody>
              <tr>
                <th>Full Name:</th>
                <td>{fullName}</td>
              </tr>
              <tr>
                <th>Username:</th>
                <td>{username}</td>
              </tr>
              <tr>
                <th>Email:</th>
                <td>{email}</td>
              </tr>
              <tr>
                <th>Organization:</th>
                <td>{orgName}</td>
              </tr>
              <tr>
                <th>Kolide URL:</th>
                <td>{kolideWebAddress}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <Button
          onClick={onSubmit}
          text="Submit"
          variant="gradient"
          className={`${baseClass}__submit`}
        />
      </div>
    );
  }
}

export default ConfirmationPage;

