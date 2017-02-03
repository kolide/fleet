import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';

import APP_CONSTANTS from 'app_constants';
import { createLicense } from 'redux/nodes/auth/actions';
import EnsureUnauthenticated from 'components/EnsureUnauthenticated';
import Footer from 'components/Footer';
import LicenseForm from 'components/forms/LicenseForm';
import { showBackgroundImage } from 'redux/nodes/app/actions';

import kolideLogo from '../../../assets/images/kolide-logo-condensed.svg';

const baseClass = 'license-page';
const { PATHS: { SETUP } } = APP_CONSTANTS;

class LicensePage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
  };

  componentWillMount () {
    const { dispatch } = this.props;

    dispatch(showBackgroundImage);

    return false;
  }

  handleSubmit = ({ license }) => {
    const { dispatch } = this.props;

    dispatch(createLicense({ license }))
      .then(() => {
        dispatch(push(SETUP));
      })
      .catch(() => false);

    return false;
  }

  render () {
    const { handleSubmit } = this;

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

const ConnectedComponent = connect()(LicensePage);
export default EnsureUnauthenticated(ConnectedComponent);
