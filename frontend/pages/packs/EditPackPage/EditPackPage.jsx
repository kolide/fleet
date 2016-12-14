import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop, size } from 'lodash';

import packActions from 'redux/nodes/entities/packs/actions';
import PackForm from 'components/forms/packs/PackForm';
import PackInfoSidePanel from 'components/side_panels/PackInfoSidePanel';
import packInterface from 'interfaces/pack';
import queryActions from 'redux/nodes/entities/queries/actions';
import queryInterface from 'interfaces/query';
import QueriesListWrapper from 'components/queries/QueriesListWrapper';
import ShowSidePanel from 'components/side_panels/ShowSidePanel';
import stateEntityGetter from 'redux/utilities/entityGetter';

const baseClass = 'edit-pack-page';

export class EditPackPage extends Component {
  static propTypes = {
    allQueries: PropTypes.arrayOf(queryInterface),
    dispatch: PropTypes.func,
    loadingPack: PropTypes.bool,
    pack: packInterface,
    packID: PropTypes.string,
    scheduledQueries: PropTypes.arrayOf(queryInterface),
  };

  static defaultProps = {
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = { selectedTargetsCount: 0 };
  }

  componentWillMount () {
    const { dispatch, loadingPack, pack, packID, allQueries } = this.props;
    const { load } = packActions;
    const { loadAll } = queryActions;

    if (!pack && !loadingPack) {
      dispatch(load(packID));
    }

    if (!size(allQueries)) {
      dispatch(loadAll());
    }

    return false;
  }

  onFetchTargets = (query, targetsResponse) => {
    const { selected_targets_count: selectedTargetsCount } = targetsResponse;

    this.setState({ selectedTargetsCount });

    return false;
  }

  handlePackFormSubmit = (formData) => {
    const { dispatch } = this.props;
    const { update } = packActions;

    return dispatch(update(formData));
  }

  render () {
    const { handlePackFormSubmit, onFetchTargets } = this;
    const { selectedTargetsCount } = this.state;
    const { allQueries, pack, scheduledQueries } = this.props;

    if (!pack) {
      return false;
    }

    return (
      <div className={`${baseClass}`}>
        <PackForm
          className={`${baseClass}__pack-form body-wrap`}
          handleSubmit={handlePackFormSubmit}
          formData={pack}
          onFetchTargets={onFetchTargets}
          selectedTargetsCount={selectedTargetsCount}
        />
        <QueriesListWrapper
          allQueries={allQueries}
          scheduledQueries={scheduledQueries}
        />
        <PackInfoSidePanel />
      </div>
    );
  }
}

const mapStateToProps = (state, { params }) => {
  const entityGetter = stateEntityGetter(state);
  const loadingPack = state.entities.packs.loading;
  const { id: packID } = params;
  const pack = entityGetter.get('packs').findBy({ id: packID });
  const { entities: allQueries } = entityGetter.get('queries');
  const scheduledQueries = []

  return { loadingPack, pack, packID, allQueries, scheduledQueries };
};

const ConnectedComponent = connect(mapStateToProps)(EditPackPage);
export default ShowSidePanel(ConnectedComponent);

