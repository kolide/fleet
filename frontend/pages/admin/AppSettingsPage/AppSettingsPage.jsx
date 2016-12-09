import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { size } from 'lodash';

import AppConfigForm from 'components/forms/admin/AppConfigForm';
import configInterface from 'interfaces/config';
import SmtpWarning from 'pages/admin/AppSettingsPage/SmtpWarning';

export const baseClass = 'app-settings-page';

class AppSettingsPage extends Component {
  static propTypes = {
    appConfig: configInterface,
    dispatch: PropTypes.func,
  };

  constructor (props) {
    super(props);

    this.state = { showSmtpWarning: true };
  }

  onDismissSmtpWarning = () => {
    this.setState({ showSmtpWarning: false });

    return false;
  }

  onFormSubmit = (formData) => {
    console.log(formData);

    return false;
  }

  render () {
    const { appConfig } = this.props;
    const { onDismissSmtpWarning, onFormSubmit } = this;
    const { showSmtpWarning } = this.state;
    const { smtp_configured: smtpConfigured } = appConfig;
    const shouldShowWarning = !smtpConfigured && showSmtpWarning;

    if (!size(appConfig)) {
      return false;
    }

    return (
      <div className={`${baseClass} body-wrap`}>
        <h1>App Settings</h1>
        <SmtpWarning
          onDismiss={onDismissSmtpWarning}
          shouldShowWarning={shouldShowWarning}
        />
        <AppConfigForm formData={appConfig} handleSubmit={onFormSubmit} />
      </div>
    );
  }
}

const mapStateToProps = ({ app }) => {
  const { config: appConfig } = app;

  return { appConfig };
};

export default connect(mapStateToProps)(AppSettingsPage);
