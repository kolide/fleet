import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';
import validate from 'components/forms/LicenseForm/validate';

import freeTrial from '../../../../assets/images/sign-up-pencil.svg';

const fields = ['license'];
const baseClass = 'license-form';

class LicenseForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      license: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func.isRequired,
  };

  render () {
    const { fields: formFields, handleSubmit } = this.props;

    return (
      <form className={baseClass} onSubmit={handleSubmit}>
        <div className={`${baseClass}__container`}>
          <h2>Kolide License</h2>
          <InputField
            {...formFields.license}
            hint={<p className={`${baseClass}__help-text`}>Found under <span>Account Settings</span> at Kolide.co</p>}
            inputClassName={`${baseClass}__input`}
            label="Enter License File"
            type="textarea"
          />
          <Button block className={`${baseClass}__upload-btn`} type="submit">
            UPLOAD LICENSE
          </Button>
          <p className="form-field__label">Don&apos;t have a license?</p>
          <p className={`${baseClass}__free-trial-text`}>Start a free trial of Kolide today!</p>
          <Button
            block
            className={`${baseClass}__free-trial-btn`}
            onClick={() => false}
            variant="unstyled"
          >
            <img
              alt="Free trial"
              src={freeTrial}
              className={`${baseClass}__free-trial-img`}
            />
            <span>Sign up for Free Kolide Trial</span>
          </Button>
        </div>
      </form>
    );
  }
}

export default Form(LicenseForm, { fields, validate });
