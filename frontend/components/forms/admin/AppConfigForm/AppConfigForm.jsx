import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Checkbox from 'components/forms/fields/Checkbox';
import Dropdown from 'components/forms/fields/Dropdown';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import Icon from 'components/Icon';
import InputField from 'components/forms/fields/InputField';
import Slider from 'components/buttons/Slider';
import validate from 'components/forms/Admin/AppConfigForm/validate';

const authMethodOptions = [
  { label: 'Plain', value: 'plain' },
  { label: 'Login', value: 'login' },
  { label: 'GSS API', value: 'gssapi' },
  { label: 'Digest MD5', value: 'digest_md5' },
  { label: 'MD5', value: 'md5' },
  { label: 'Cram MD5', value: 'cram_md5' },
];
const authTypeOptions = [
  { label: 'Username and Password', value: 'username_and_password' },
  { label: 'None', value: 'none' },
];
const baseClass = 'app-config-form';
const formFields = [
  'auth_method', 'authentication_type', 'domain', 'enable_ssl_tls', 'enable_start_tls',
  'kolide_server_url', 'org_logo_url', 'org_name', 'password', 'port', 'sender_address',
  'server', 'user_name', 'verify_ssl_certs',
];
const Header = ({ showAdvancedOptions }) => {
  const CaratIcon = <Icon name={showAdvancedOptions ? 'downcarat' : 'upcarat'} />;

  return <span>Advanced Options {CaratIcon}</span>;
};

Header.propTypes = { showAdvancedOptions: PropTypes.bool.isRequired };

class AppConfigForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      auth_method: formFieldInterface.isRequired,
      authentication_type: formFieldInterface.isRequired,
      domain: formFieldInterface.isRequired,
      enable_ssl_tls: formFieldInterface.isRequired,
      enable_start_tls: formFieldInterface.isRequired,
      kolide_server_url: formFieldInterface.isRequired,
      org_logo_url: formFieldInterface.isRequired,
      org_name: formFieldInterface.isRequired,
      password: formFieldInterface.isRequired,
      port: formFieldInterface.isRequired,
      sender_address: formFieldInterface.isRequired,
      server: formFieldInterface.isRequired,
      user_name: formFieldInterface.isRequired,
      verify_ssl_certs: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func,
    smtpConfigured: PropTypes.bool,
  };

  constructor (props) {
    super(props);

    this.state = { showAdvancedOptions: false };
  }

  onToggleAdvancedOptions = (evt) => {
    evt.preventDefault();

    const { showAdvancedOptions } = this.state;

    this.setState({ showAdvancedOptions: !showAdvancedOptions });

    return false;
  }

  renderAdvancedOptions = () => {
    const { fields } = this.props;
    const { showAdvancedOptions } = this.state;

    if (!showAdvancedOptions) {
      return false;
    }

    return (
      <div>
        <InputField {...fields.domain} />
        <Slider {...fields.verify_ssl_certs} />
        <Slider {...fields.enable_start_tls} />
      </div>
    );
  }

  render () {
    const { fields, handleSubmit, smtpConfigured } = this.props;
    const { onToggleAdvancedOptions, renderAdvancedOptions } = this;
    const { showAdvancedOptions } = this.state;

    return (
      <form className={baseClass}>
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
          <h2>SMTP Options <small>STATUS: {smtpConfigured ? 'CONFIGURED' : 'NOT CONFIGURED'}</small></h2>
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
            {...fields.enable_ssl_tls}
          >
            User SSL/TLS to connect (recommended)
          </Checkbox>
          <Dropdown
            {...fields.authentication_type}
            options={authTypeOptions}
          />
          <div className={`${baseClass}__smtp-user-section`}>
            <InputField
              {...fields.user_name}
            />
            <InputField
              {...fields.password}
            />
            <Dropdown
              {...fields.auth_method}
              options={authMethodOptions}
              placeholder=""
            />
          </div>
        </div>
        <div className={`${baseClass}__section`}>
          <h2>
            <Button
              onClick={onToggleAdvancedOptions}
              text={<Header showAdvancedOptions={showAdvancedOptions} />}
              variant="unstyled"
            />
          </h2>
          {renderAdvancedOptions()}
        </div>
        <Button
          onClick={handleSubmit}
          text="UPDATE SETTINGS"
          variant="brand"
        />
      </form>
    );
  }
}

export default Form(AppConfigForm, {
  fields: formFields,
  validate,
});

