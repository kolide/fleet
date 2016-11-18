import React, { Component } from 'react';
import { connect } from 'react-redux';

import PackForm from 'components/forms/PackForm';

export class PackComposerPage extends Component {
  constructor (props) {
    super(props);

    this.state = { selectedTargetsCount: 0 };
  }

  onFetchTargets = (query, targetsResponse) => {
    const {
      selected_targets_count: selectedTargetsCount,
    } = targetsResponse;

    this.setState({ selectedTargetsCount });

    return false;
  }

  handleSubmit = (formData) => {
    console.log(formData);

    return false;
  }

  render () {
    const { handleSubmit, onFetchTargets } = this;
    const { selectedTargetsCount } = this.state;

    return (
      <div>
        <PackForm
          handleSubmit={handleSubmit}
          onFetchTargets={onFetchTargets}
          selectedTargetsCount={selectedTargetsCount}
        />
      </div>
    );
  }
}

export default connect()(PackComposerPage);
