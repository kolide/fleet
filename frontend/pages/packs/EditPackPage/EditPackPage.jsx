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
import { renderFlash } from 'redux/nodes/notifications/actions';
import scheduledQueryActions from 'redux/nodes/entities/scheduled_queries/actions';
import ShowSidePanel from 'components/side_panels/ShowSidePanel';
import stateEntityGetter from 'redux/utilities/entityGetter';

const baseClass = 'edit-pack-page';

export class EditPackPage extends Component {
  static propTypes = {
    allQueries: PropTypes.arrayOf(queryInterface),
    dispatch: PropTypes.func,
    isLoadingPack: PropTypes.bool,
    isLoadingScheduledQueries: PropTypes.bool,
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

  componentDidMount () {
    const { allQueries, dispatch, isLoadingPack, pack, packID, scheduledQueries } = this.props;
    const { load } = packActions;
    const { loadAll } = queryActions;

    if (!pack && !isLoadingPack) {
      dispatch(load(packID));
    }

    if (!size(scheduledQueries)) {
      dispatch(scheduledQueryActions.loadAll({ id: packID }));
    }

    if (!size(allQueries)) {
      dispatch(loadAll());
    }

    return false;
  }

  onFetchTargets = (query, targetsResponse) => {
    const { targets_count: selectedTargetsCount } = targetsResponse;

    this.setState({ selectedTargetsCount });

    return false;
  }

  handlePackFormSubmit = (formData) => {
    const { dispatch } = this.props;
    const { update } = packActions;

    return dispatch(update(formData));
  }

  handleScheduledQueryFormSubmit = (formData) => {
    const { create } = scheduledQueryActions;
    const { dispatch, packID } = this.props;
    const scheduledQueryData = {
      ...formData,
      snapshot: formData.logging_type === 'snapshot',
      pack_id: packID,
    };

    dispatch(create(scheduledQueryData))
      .then(() => {
        dispatch(renderFlash('success', 'Query scheduled!'));
      })
      .catch(() => {
        dispatch(renderFlash('error', 'Unable to schedule your query.'));
      });

    return false;
  }

  render () {
    const { handlePackFormSubmit, handleScheduledQueryFormSubmit, onFetchTargets } = this;
    const { selectedTargetsCount } = this.state;
    const { allQueries, isLoadingScheduledQueries, pack, scheduledQueries } = this.props;

    if (!pack || isLoadingScheduledQueries) {
      return false;
    }

    return (
      <div className={`${baseClass} has-sidebar`}>
        <div className="body-wrap body-wrap--unstyled">
          <PackForm
            className={`${baseClass}__pack-form body-wrap`}
            handleSubmit={handlePackFormSubmit}
            formData={pack}
            onFetchTargets={onFetchTargets}
            selectedTargetsCount={selectedTargetsCount}
          />
          <QueriesListWrapper
            allQueries={allQueries}
            onScheduledQueryFormSubmit={handleScheduledQueryFormSubmit}
            scheduledQueries={scheduledQueries}
          />
        </div>
        <PackInfoSidePanel />
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const entityGetter = stateEntityGetter(state);
  const isLoadingPack = state.entities.packs.loading;
  const { id: packID } = ownProps.params;
  const pack = entityGetter.get('packs').findBy({ id: packID });
  const { entities: allQueries } = entityGetter.get('queries');
  const scheduledQueries = entityGetter.get('scheduled_queries').where({ pack_id: packID });
  const isLoadingScheduledQueries = state.entities.scheduled_queries.loading;

  return {
    allQueries,
    isLoadingPack,
    isLoadingScheduledQueries,
    pack,
    packID,
    scheduledQueries,
  };
};

const ConnectedComponent = connect(mapStateToProps)(EditPackPage);
export default ShowSidePanel(ConnectedComponent);

