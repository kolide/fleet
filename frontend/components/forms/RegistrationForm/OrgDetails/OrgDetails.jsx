import React from 'react';

import BasePageForm from 'components/forms/RegistrationForm/BasePageForm';
import Button from 'components/buttons/Button';
import InputFieldWithIcon from 'components/forms/fields/InputFieldWithIcon';

class OrgDetails extends BasePageForm {
  valid = () => {
    const { errors } = this.state;
    const {
      formData: {
        org_name: orgName,
        org_web_url: orgWebUrl,
        org_logo_url: orgLogoUrl,
      },
    } = this.props;

    if (orgName && orgWebUrl && orgLogoUrl) {
      return true;
    }

    this.setState({
      errors: {
        ...errors,
        org_name: !orgName ? 'Organization name must be present' : null,
        org_web_url: !orgWebUrl ? 'Organization web URL must be present' : null,
        org_logo_url: !orgLogoUrl ? 'Organization logo URL must be present' : null,
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
          error={errors('org_name')}
          name="organization name"
          onChange={onChange('org_name')}
          placeholder="Organization Name"
          value={formData.org_name}
        />
        <InputFieldWithIcon
          error={errors('org_web_url')}
          name="org web url"
          onChange={onChange('org_web_url')}
          placeholder="Organization Website URL"
          value={formData.org_web_url}
        />
        <InputFieldWithIcon
          error={errors('org_logo_url')}
          name="org logo url"
          onChange={onChange('org_logo_url')}
          placeholder="Organization Logo URL"
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
