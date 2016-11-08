import React from 'react';
import { size, startsWith } from 'lodash';

import BasePageForm from 'components/forms/RegistrationForm/BasePageForm';
import Button from 'components/buttons/Button';
import InputFieldWithIcon from 'components/forms/fields/InputFieldWithIcon';

class KolideDetails extends BasePageForm {
  valid = () => {
    const clientErrors = {};
    const {
      errors,
      formData: {
        kolide_web_address: kolideWebAddress,
      },
    } = this.props;

    if (!kolideWebAddress) {
      clientErrors.kolide_web_address = 'Kolide web address must be completed';
    }

    if (kolideWebAddress && !startsWith(kolideWebAddress, 'https://')) {
      clientErrors.kolide_web_address = 'Kolide web address must start with https://';
    }

    if (size(clientErrors)) {
      this.setState({
        errors: {
          ...errors,
          ...clientErrors,
        },
      });

      return false;
    }

    return true;
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

