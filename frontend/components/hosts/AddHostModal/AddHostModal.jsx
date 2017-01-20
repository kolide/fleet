import React, { PropTypes } from 'react';
import { noop } from 'lodash';

import Button from 'components/buttons/Button';

const baseClass = 'add-host-modal';

const AddHostModal = ({ onReturnToApp = noop }) => {
  return (
    <div className={baseClass}>
      <p>Follow the instructions below to add hosts to your Kolide Instance.</p>

      <div className={`${baseClass}__manual-install-header`}>
        <h2>Manual Install</h2>
        <h3>Fully Customize Your <strong>Osquery</strong> Installation</h3>
      </div>

      <div className={`${baseClass}__manual-install-content`}>
        <ol className={`${baseClass}__install-steps`}>
          <li>
            <h4><a href="#linkToInstallDocs">Kolide / Osquery - Install Docs</a></h4>
            <p>In order to install <strong>osquery</strong> on a client you will need the items below:</p>
          </li>
          <li>
            <h4>Download Osquery Package and Certificate</h4>
            <p>Osquery requires the same TLS certificate that Kolide is using in order to authenticate. You can fetch the certificate below:</p>
          </li>
          <li>
            <h4>Retrieve Osquery Enroll Secret</h4>
            <p>When prompted, enter the provided secret code into <strong>osqueryd</strong>:</p>
          </li>
        </ol>
      </div>

      <div className={`${baseClass}__button-wrap`}>
        <Button onClick={onReturnToApp} variant="success">
          Return To App
        </Button>
      </div>
    </div>
  );
};

AddHostModal.propTypes = {
  onReturnToApp: PropTypes.func,
};

export default AddHostModal;
