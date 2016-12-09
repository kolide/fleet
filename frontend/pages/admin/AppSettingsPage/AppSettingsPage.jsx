import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

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

  render () {
    const { appConfig } = this.props;
    const { onDismissSmtpWarning } = this;
    const { showSmtpWarning } = this.state;
    const { smtp_configured: smtpConfigured } = appConfig;
    const shouldShowWarning = !smtpConfigured && showSmtpWarning;

    return (
      <div className={`${baseClass} body-wrap`}>
        <h1>App Settings</h1>
        <SmtpWarning
          onDismiss={onDismissSmtpWarning}
          shouldShowWarning={shouldShowWarning}
        />
      </div>
    );
  }
}

const mapStateToProps = ({ app }) => {
  const { config: appConfig } = app;

  return { appConfig };
};

export default connect(mapStateToProps)(AppSettingsPage);
