import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

class UserSettingsPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
  };

  render () {
    return (
      <div className="body-wrap">
        <h1>User Settings Page</h1>
      </div>
    );
  }
}

export default connect()(UserSettingsPage);
