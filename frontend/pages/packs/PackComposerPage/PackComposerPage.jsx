import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import { push } from 'react-router-redux';

import Kolide from 'kolide';
import packActions from 'redux/nodes/entities/packs/actions';
import PackForm from 'components/forms/packs/PackForm';
import PackInfoSidePanel from 'components/side_panels/PackInfoSidePanel';
import { renderFlash } from 'redux/nodes/notifications/actions';

const baseClass = 'pack-composer';

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
    const { create, load } = packActions;

    return dispatch(create(formData))
      .then((pack) => {
        const { id: packID } = pack;
        const { targets } = formData;

        const promises = targets.map((target) => {
          const { id: targetID } = target;

          if (target.target_type === 'labels') {
            Kolide.addLabelToPack(packID, targetID);
          }

          // TODO: Add host to pack when API is available
          return Promise.resolve();
        });

        Promise.all(promises)
          .then(() => {
            dispatch(load(packID))
              .then(() => {
                dispatch(push(`/packs/${packID}`));
                dispatch(renderFlash('success', 'Pack created!'));
              });
          });

        return false;
      });
  }

  render () {
    const { handleSubmit, onFetchTargets } = this;
    const { selectedTargetsCount } = this.state;

    return (
      <div className="has-sidebar">
        <PackForm
          className={`${baseClass}__pack-form body-wrap`}
          handleSubmit={handleSubmit}
          onFetchTargets={onFetchTargets}
          selectedTargetsCount={selectedTargetsCount}
        />
        <PackInfoSidePanel />
      </div>
    );
  }
}

export default connect()(PackComposerPage);
