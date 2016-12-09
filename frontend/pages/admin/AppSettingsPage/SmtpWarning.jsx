import React, { PropTypes } from 'react';

import { baseClass } from 'pages/admin/AppSettingsPage/AppSettingsPage';
import Button from 'components/buttons/Button';
import Icon from 'components/Icon';


const SmtpWarning = ({ onDismiss, shouldShowWarning }) => {
  if (!shouldShowWarning) {
    return false;
  }

  return (
    <div className={`${baseClass}__smtp-warning`}>
      <div>
        <Icon name="warning-filled" />
        <span className={`${baseClass}__smtp-warning-label`}>Warning!</span>
      </div>
      <span>Email is not currently configured in Kolide. Many features rely on email to work.</span>
      <Button onClick={onDismiss} text="DISMISS" variant="unstyled" />
      <Button text="RESOLVE NOW" variant="unstyled" />
    </div>
  );
};

SmtpWarning.propTypes = {
  onDismiss: PropTypes.func.isRequired,
  shouldShowWarning: PropTypes.bool.isRequired,
}

export default SmtpWarning;
