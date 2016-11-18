import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';

import PackForm from 'components/forms/PackForm';

export class PackComposerPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
  };

  static defaultProps = {
    dispatch: noop,
  };

  handleSubmit = (formData) => {
    console.log(formData);

    return false;
  }

  render () {
    const { handleSubmit } = this;

    return (
      <div>
        <PackForm handleSubmit={handleSubmit} />
      </div>
    );
  }
}

export default connect()(PackComposerPage);
