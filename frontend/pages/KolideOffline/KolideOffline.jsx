import React, { Component } from 'react';

import kolideLogo from '../../../assets/images/kolide-logo-condensed.svg';
import gopher from '../../../assets/images/404.svg';

const baseClass = 'kolide-offline';

class KolideOffline extends Component {

  render () {
    return (
      <div className={baseClass}>
        <header className="primary-header">
          <a href="/">
            <img className="primary-header__logo" src={kolideLogo} alt="Kolide" />
          </a>
        </header>
        <main>
          <h1>Oh noes, you're out of internets!</h1>
          <h2>Offline!</h2>
          <p>You seem to have lost your connection.</p>
          <p>Try deleting system32 directory.</p>
          <div className="gopher-container">
            <img src={gopher} role="presentation" />
            <p>Need immediate assistance? <br />Contact <a href="mailto:support@kolide.co">support@kolide.co</a></p>
          </div>
        </main>
      </div>
    );
  }
}

export default KolideOffline;
