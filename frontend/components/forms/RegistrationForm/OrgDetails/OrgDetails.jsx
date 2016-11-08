import React from 'react';
import { size, startsWith } from 'lodash';

import BasePageForm from 'components/forms/RegistrationForm/BasePageForm';
import Button from 'components/buttons/Button';
import InputFieldWithIcon from 'components/forms/fields/InputFieldWithIcon';

class OrgDetails extends BasePageForm {
  valid = () => {
    const clientErrors = {};
    const { errors } = this.state;
    const {
      formData: {
        org_name: orgName,
        org_logo_url: orgLogoUrl,
      },
    } = this.props;

    if (!orgName) {
      clientErrors.org_name = 'Organization name must be present';
    }

    if (!orgLogoUrl) {
      clientErrors.org_logo_url = 'Organization logo URL must be present';
    }

    if (orgLogoUrl && !startsWith(orgLogoUrl, 'https://')) {
      clientErrors.org_logo_url = 'Organization logo URL must start with https://';
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
          error={errors('org_name')}
          name="organization name"
          onChange={onChange('org_name')}
          placeholder="Organization Name"
          value={formData.org_name}
        />
        <InputFieldWithIcon
          error={errors('org_logo_url')}
          name="org logo url"
          onChange={onChange('org_logo_url')}
          placeholder="Organization Logo URL (must start with https://)"
          value={formData.org_logo_url}
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

export default OrgDetails;
