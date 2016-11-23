import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';

import PackForm from 'components/forms/packs/PackForm';
import packsPageActions from 'redux/nodes/components/PacksPages/actions';
import queryActions from 'redux/nodes/entities/queries/actions';
import queryInterface from 'interfaces/query';
import QueriesListWrapper from 'components/queries/QueriesListWrapper';
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
    configuredQueries: PropTypes.arrayOf(queryInterface),
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

  onConfigureQueries = (formData) => {
    console.log('configure queries', formData);
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
    const { configurations } = this.props;
    console.log('pack form data', formData);
    console.log('configurations', configurations);

    return false;
  }

  render () {
    const {
      handleSubmit,
      onConfigureQueries,
      onFetchTargets,
      onStageQuery,
      onUnstageQuery,
    } = this;
    const { selectedTargetsCount } = this.state;
    const { allQueries, configuredQueries, stagedQueries } = this.props;

    return (
      <div className={baseClass}>
        <PackForm
          className={`${baseClass}__pack-form`}
          handleSubmit={handleSubmit}
          onFetchTargets={onFetchTargets}
          selectedTargetsCount={selectedTargetsCount}
        />
        <QueriesListWrapper
          configuredQueries={configuredQueries}
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
  const { configurations, configuredQueries, stagedQueries } = state.components.PacksPages;

  return { allQueries: queries, configurations, configuredQueries, stagedQueries };
};

export default connect(mapStateToProps)(PackComposerPage);
