import React, { Component, PropTypes } from 'react';

import Checkbox from 'components/forms/fields/Checkbox';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';
import validate from 'components/forms/admin/AppConfigForm/validate';

const baseClass = 'app-config-form';
const formFields = [
  'auth_method', 'authentication_type', 'domain', 'enable_sll_tls', 'enable_start_tls',
  'kolide_server_url', 'org_logo_url', 'org_name', 'password', 'port', 'sender_address',
  'server', 'smtp_configured', 'user_name', 'verify_sll_certs',
];

class AppConfigForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      auth_method: formFieldInterface.isRequired,
      authentication_type: formFieldInterface.isRequired,
      domain: formFieldInterface.isRequired,
      enable_sll_tls: formFieldInterface.isRequired,
      enable_start_tls: formFieldInterface.isRequired,
      kolide_server_url: formFieldInterface.isRequired,
      org_logo_url: formFieldInterface.isRequired,
      org_name: formFieldInterface.isRequired,
      password: formFieldInterface.isRequired,
      port: formFieldInterface.isRequired,
      sender_address: formFieldInterface.isRequired,
      server: formFieldInterface.isRequired,
      user_name: formFieldInterface.isRequired,
      verify_sll_certs: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func,
  };

  render () {
    const { fields, formData, handleSubmit } = this.props;

    return (
      <form onSubmit={handleSubmit} className={baseClass}>
        <div className={`${baseClass}__section`}>
          <h2>Organization Info</h2>
          <InputField
            {...fields.org_name}
            label="Organization Name"
          />
          <InputField
            {...fields.org_logo_url}
            label="Organization Avatar"
          />
        </div>
        <div className={`${baseClass}__section`}>
          <h2>Kolide Web Address</h2>
          <InputField
            {...fields.kolide_server_url}
            label="Kolide App URL"
          />
        </div>
        <div className={`${baseClass}__section`}>
          <h2>SMTP Options</h2>
          <InputField
            {...fields.sender_address}
            label="Sender Address"
          />
          <InputField
            {...fields.server}
            label="SMTP Server"
          />
          <InputField {...fields.port} />
          <Checkbox
            {...fields.enable_sll_tls}
          >
            User SSL/TLS to connect (recommended)
          </Checkbox>
        </div>
      </form>
    );
  }
}

export default Form(AppConfigForm, {
  fields: formFields,
  validate,
});

