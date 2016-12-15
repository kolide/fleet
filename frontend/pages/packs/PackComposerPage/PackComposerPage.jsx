import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import { push } from 'react-router-redux';

import packActions from 'redux/nodes/entities/packs/actions';
import PackForm from 'components/forms/packs/PackForm';
import PackInfoSidePanel from 'components/side_panels/PackInfoSidePanel';
import { renderFlash } from 'redux/nodes/notifications/actions';
import ShowSidePanel from 'components/side_panels/ShowSidePanel';

const baseClass = 'pack-composer-page';

export class PackComposerPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
  };

  static defaultProps = {
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = { selectedTargetsCount: 0 };
  }

  onFetchTargets = (query, targetsResponse) => {
    const { targets_count: selectedTargetsCount } = targetsResponse;

    this.setState({ selectedTargetsCount });

    return false;
  }

  handleSubmit = (formData) => {
    const { dispatch } = this.props;
    const { create } = packActions;

    return dispatch(create(formData))
      .then((pack) => {
        const { id: packID } = pack;

        return dispatch(push(`/packs/${packID}`));
      });
  }

  render () {
    const { handleSubmit, onFetchTargets } = this;
    const { selectedTargetsCount } = this.state;

    return (
      <div className={`${baseClass} body-wrap`}>
        <PackForm
          className={`${baseClass}__pack-form`}
          handleSubmit={handleSubmit}
          onFetchTargets={onFetchTargets}
          selectedTargetsCount={selectedTargetsCount}
        />
        <PackInfoSidePanel />
      </div>
    );
  }
}

const ConnectedComponent = connect()(PackComposerPage);
export default ShowSidePanel(ConnectedComponent);
