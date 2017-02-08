import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';

import APP_CONSTANTS from 'app_constants';
import { createLicense } from 'redux/nodes/auth/actions';
import EnsureUnauthenticated from 'components/EnsureUnauthenticated';
import Footer from 'components/Footer';
import LicenseForm from 'components/forms/LicenseForm';
import licenseInterface from 'interfaces/license';
import LicenseSuccess from 'components/LicenseSuccess';
import { showBackgroundImage } from 'redux/nodes/app/actions';

import kolideLogo from '../../../assets/images/kolide-logo-condensed.svg';

const baseClass = 'license-page';
const { PATHS: { SETUP } } = APP_CONSTANTS;

class LicensePage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    license: licenseInterface.isRequired,
  };

  componentWillMount () {
    const { dispatch } = this.props;

    dispatch(showBackgroundImage);

    return false;
  }

  onConfirmLicense = (evt) => {
    evt.preventDefault();

    const { dispatch } = this.props;

    dispatch(push(SETUP));

    return false;
  }

  handleSubmit = ({ license }) => {
    const { dispatch } = this.props;

    dispatch(createLicense({ license }))
      .catch(() => false);

    return false;
  }

  render () {
    const { handleSubmit, onConfirmLicense } = this;
    const { license } = this.props;

    if (license.token) {
      return (
        <div className={baseClass}>
          <img
            alt="Kolide"
            src={kolideLogo}
            className={`${baseClass}__logo`}
          />
          <LicenseSuccess license={license} onConfirmLicense={onConfirmLicense} />
          <Footer />
        </div>
      );
    }

    return (
      <div className={baseClass}>
        <img
          alt="Kolide"
          src={kolideLogo}
          className={`${baseClass}__logo`}
        />
        <LicenseForm handleSubmit={handleSubmit} />
        <Footer />
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { license } = state.auth;

  return { license };
};

const ConnectedComponent = connect(mapStateToProps)(LicensePage);
export default EnsureUnauthenticated(ConnectedComponent);
