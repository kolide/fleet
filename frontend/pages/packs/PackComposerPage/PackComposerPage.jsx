import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import { push } from 'react-router-redux';

import Kolide from 'kolide';
import PackForm from 'components/forms/packs/PackForm';
import packActions from 'redux/nodes/entities/packs/actions';
import packsPageActions from 'redux/nodes/components/PacksPages/actions';
import queryActions from 'redux/nodes/entities/queries/actions';
import queryInterface from 'interfaces/query';
import QueriesListWrapper from 'components/queries/QueriesListWrapper';
import { renderFlash } from 'redux/nodes/notifications/actions';
import stateEntityGetter from 'redux/utilities/entityGetter';

const baseClass = 'pack-composer-page';

export class PackComposerPage extends Component {
  static propTypes = {
    allQueries: PropTypes.arrayOf(queryInterface),
    configurations: PropTypes.arrayOf(PropTypes.shape({
      interval: PropTypes.string,
      platform: PropTypes.string,
      logging_type: PropTypes.string,
      query_ids: PropTypes.arrayOf(PropTypes.number),
    })),
    configuredQueryIDs: PropTypes.arrayOf(PropTypes.number),
    dispatch: PropTypes.func,
    stagedQueries: PropTypes.arrayOf(queryInterface),
  };

  static defaultProps = {
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = { selectedTargetsCount: 0 };
  }

  componentDidMount () {
    const { dispatch } = this.props;

    dispatch(queryActions.loadAll());
  }

  onClearStagedQueries = () => {
    const { clearStagedQueries } = packsPageActions;
    const { dispatch } = this.props;

    dispatch(clearStagedQueries);

    return false;
  }

  onConfigureQueries = (formData) => {
    const { configureStagedQueries } = packsPageActions;
    const { dispatch } = this.props;

    dispatch(configureStagedQueries(formData));
  }

  onFetchTargets = (query, targetsResponse) => {
    const {
      selected_targets_count: selectedTargetsCount,
    } = targetsResponse;

    this.setState({ selectedTargetsCount });

    return false;
  }

  onStageQuery = (query) => {
    const { dispatch } = this.props;
    const { stageQuery } = packsPageActions;

    dispatch(stageQuery(query));

    return false;
  }

  onUnstageQuery = (query) => {
    const { dispatch } = this.props;
    const { unstageQuery } = packsPageActions;

    dispatch(unstageQuery(query));

    return false;
  }

  handleSubmit = (formData) => {
    const { configurations, dispatch } = this.props;
    const { load, create } = packActions;

    dispatch(create(formData))
      .then((pack) => {
        const { id: packID } = pack;
        const promises = [];

        configurations.forEach((configuration) => {
          configuration.query_ids.forEach((queryID) => {
            promises.push(Kolide.addQueryToPack({ packID, queryID }));
          });
        });

        promises.push(dispatch(load(packID)));

        Promise.all(promises)
          .then(() => {
            dispatch(push('/packs/all'));
            dispatch(renderFlash('success', 'Pack created!'));
          })
          .catch(() => {
            dispatch(renderFlash('error', 'There was an error creating your pack'));
          });
      });

    console.log('configurations data', configurations);

    return false;
  }

  render () {
    const {
      handleSubmit,
      onClearStagedQueries,
      onConfigureQueries,
      onFetchTargets,
      onStageQuery,
      onUnstageQuery,
    } = this;
    const { selectedTargetsCount } = this.state;
    const { allQueries, configuredQueryIDs, stagedQueries } = this.props;

    return (
      <div className={baseClass}>
        <PackForm
          className={`${baseClass}__pack-form`}
          handleSubmit={handleSubmit}
          onFetchTargets={onFetchTargets}
          selectedTargetsCount={selectedTargetsCount}
        />
        <QueriesListWrapper
          configuredQueryIDs={configuredQueryIDs}
          onClearStagedQueries={onClearStagedQueries}
          onConfigureQueries={onConfigureQueries}
          onDeselectQuery={onUnstageQuery}
          onSelectQuery={onStageQuery}
          queries={allQueries}
          stagedQueries={stagedQueries}
        />
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: queries } = stateEntityGetter(state).get('queries');
  const { configurations, configuredQueryIDs, stagedQueries } = state.components.PacksPages;

  return { allQueries: queries, configurations, configuredQueryIDs, stagedQueries };
};

export default connect(mapStateToProps)(PackComposerPage);
