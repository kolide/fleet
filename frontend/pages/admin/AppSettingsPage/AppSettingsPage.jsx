import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

class AppSettingsPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
  };

  render () {
    return (
      <div>
        <h1>App Settings</h1>
      </div>
    );
  }
}

export default connect()(AppSettingsPage);
