import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';

import Button from 'components/buttons/Button';
import formDataInterface from 'interfaces/registration_form_data';
import Icon from 'components/Icon';

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
        name,
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
                <td>{name}</td>
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
                <td><span className={`${baseClass}__table-url`} title={kolideWebAddress}>{kolideWebAddress}</span></td>
              </tr>
            </tbody>
          </table>

          <div className={`${baseClass}__import`}>
            <label htmlFor="import-install">
              <input type="checkbox" name="import-install" id="import-install" className="kolide-checkbox" />
              <p>I am migrating an existing <strong>osquery</strong> installation.</p>
              <p>Take me to the <strong>Import Configuration</strong> page.</p>
            </label>
          </div>
        </div>

        <Button
          onClick={onSubmit}
          text="Finish"
          variant="gradient"
          className={`${baseClass}__submit`}
        />
      </div>
    );
  }
}

export default ConfirmationPage;

