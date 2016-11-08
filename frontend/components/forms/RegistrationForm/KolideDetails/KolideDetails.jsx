import React from 'react';

import BasePageForm from 'components/forms/RegistrationForm/BasePageForm';
import Button from 'components/buttons/Button';
import InputFieldWithIcon from 'components/forms/fields/InputFieldWithIcon';

class KolideDetails extends BasePageForm {
  valid = () => {
    const {
      errors,
      formData: {
        kolide_web_address: kolideWebAddress,
      },
    } = this.props;

    if (kolideWebAddress) {
      return true;
    }

    this.setState({
      errors: {
        ...errors,
        kolide_web_address: 'Kolide web address must be completed',
      },
    });

    return false;
  }

  render () {
    const { formData } = this.props;
    const { errors, onChange, onSubmit } = this;

    return (
      <div>
        <InputFieldWithIcon
          error={errors('kolide_web_address')}
          name="kolide web address"
          onChange={onChange('kolide_web_address')}
          placeholder="Kolide Web Address"
          value={formData.kolide_web_address}
        />
        <Button
          onClick={onSubmit}
          text="Submit"
          variant="gradient"
        />
      </div>
    );
  }
}

export default KolideDetails;

