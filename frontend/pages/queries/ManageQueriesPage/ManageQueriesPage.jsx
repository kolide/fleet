import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { filter, get, includes, isEqual, noop, pull } from 'lodash';
import { push } from 'react-router-redux';

import Button from 'components/buttons/Button';
import entityGetter from 'redux/utilities/entityGetter';
import Icon from 'components/icons/Icon';
import InputField from 'components/forms/fields/InputField';
import NumberPill from 'components/NumberPill';
import packInterface from 'interfaces/pack';
import PacksList from 'components/packs/PacksList';
import paths from 'router/paths';
import QueriesList from 'components/queries/QueriesList';
import queryActions from 'redux/nodes/entities/queries/actions';
import queryInterface from 'interfaces/query';
import { renderFlash } from 'redux/nodes/notifications/actions';

const baseClass = 'manage-queries-page';

export class ManageQueriesPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    queries: PropTypes.arrayOf(queryInterface),
    selectedQuery: queryInterface,
    selectedPacks: PropTypes.arrayOf(packInterface),
  }

  static defaultProps = {
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = {
      allQueriesChecked: false,
      checkedQueryIDs: [],
      queriesFilter: '',
    };
  }

  componentWillMount() {
    const { dispatch, queries, selectedQuery } = this.props;

    if (!queries.length) {
      dispatch(queryActions.loadAll());
    }

    if (selectedQuery) {
      // TODO: Get packs for the query
    }

    return false;
  }

  componentWillReceiveProps ({ selectedQuery }) {
    if (!isEqual(this.props.selectedQuery, selectedQuery)) {
      // TODO: Get packs for the query
    }

    return false;
  }

  onDeleteQueries = (evt) => {
    evt.preventDefault();

    const { checkedQueryIDs } = this.state;
    const { dispatch } = this.props;
    const { destroy } = queryActions;

    const promises = checkedQueryIDs.map((queryID) => {
      return dispatch(destroy({ id: queryID }));
    });

    return Promise.all(promises)
      .then(() => {
        dispatch(renderFlash('success', 'Queries successfully deleted.'));

        return false;
      })
      .catch(() => dispatch(renderFlash('error', 'Something went wrong.')));
  }

  onCheckAllQueries = (shouldCheck) => {
    if (shouldCheck) {
      const queries = this.getQueries();
      const checkedQueryIDs = queries.map(query => query.id);

      this.setState({ allQueriesChecked: true, checkedQueryIDs });

      return false;
    }

    this.setState({ allQueriesChecked: false, checkedQueryIDs: [] });

    return false;
  }

  onCheckQuery = (checked, id) => {
    const { checkedQueryIDs } = this.state;
    const newCheckedQueryIDs = checked ? checkedQueryIDs.concat(id) : pull(checkedQueryIDs, id);

    this.setState({ allQueriesChecked: false, checkedQueryIDs: newCheckedQueryIDs });

    return false;
  }

  onFilterQueries = (queriesFilter) => {
    this.setState({ queriesFilter });

    return false;
  }

  onSelectQuery = (selectedQuery) => {
    const { dispatch } = this.props;
    const locationObject = {
      pathname: '/queries/manage',
      query: { selectedQuery: selectedQuery.id },
    };

    dispatch(push(locationObject));

    return false;
  }

  getQueries = () => {
    const { queriesFilter } = this.state;
    const { queries } = this.props;

    if (!queriesFilter) {
      return queries;
    }

    const lowerQueryFilter = queriesFilter.toLowerCase();

    return filter(queries, (query) => {
      if (!query.name) {
        return false;
      }

      const lowerQueryName = query.name.toLowerCase();

      return includes(lowerQueryName, lowerQueryFilter);
    });
  }

  goToNewQueryPage = () => {
    const { dispatch } = this.props;
    const { NEW_QUERY } = paths;

    dispatch(push(NEW_QUERY));

    return false;
  }

  renderCTAs = () => {
    const { goToNewQueryPage, onDeleteQueries } = this;
    const btnClass = `${baseClass}__delete-queries-btn`;
    const checkedQueryCount = this.state.checkedQueryIDs.length;

    if (checkedQueryCount) {
      const queryText = checkedQueryCount === 1 ? 'Query' : 'Queries';

      return (
        <div>
          <p className={`${baseClass}__query-count`}>{checkedQueryCount} {queryText} Selected</p>
          <Button
            className={btnClass}
            onClick={onDeleteQueries}
            variant="alert"
          >
            <Icon name="offline" /> Disable
          </Button>
        </div>
      );
    }

    return (
      <Button variant="brand" onClick={goToNewQueryPage}>CREATE NEW QUERY</Button>
    );
  }

  renderSidePanel = () => {
    const { selectedQuery, selectedPacks } = this.props;

    // TODO: Render QueryDetailsSidePanel

    return false;
  }

  render () {
    const { allQueriesChecked, checkedQueryIDs, queriesFilter } = this.state;
    const {
      getQueries,
      onCheckAllQueries,
      onCheckQuery,
      onSelectQuery,
      onFilterQueries,
      renderCTAs,
      renderSidePanel,
    } = this;
    const { queries: allQueries, selectedQuery } = this.props;
    const queries = getQueries();
    const queriesCount = queries.length;
    const isQueriesAvailable = allQueries.length > 0;

    return (
      <div className={`${baseClass} has-sidebar`}>
        <div className={`${baseClass}__wrapper body-wrap`}>
          <p className={`${baseClass}__title`}>
            <NumberPill number={queriesCount} /> Queries
          </p>
          <div className={`${baseClass}__search-create-section`}>
            <InputField
              name="query-filter"
              onChange={onFilterQueries}
              placeholder="Search Queries"
              value={queriesFilter}
            />
            {renderCTAs()}
          </div>
          <QueriesList
            checkedQueryIDs={checkedQueryIDs}
            className={`${baseClass}__table`}
            isQueriesAvailable={isQueriesAvailable}
            onCheckAll={onCheckAllQueries}
            onCheckQuery={onCheckQuery}
            onSelectQuery={onSelectQuery}
            queries={queries}
            selectedQuery={selectedQuery}
          />
        </div>
        {renderSidePanel()}
      </div>
    );
  }
}

const mapStateToProps = (state, { location }) => {
  const queryEntities = entityGetter(state).get('queries');
  const packEntities = entityGetter(state).get('packs');
  const { entities: queries } = queryEntities;
  const selectedQueryID = get(location, 'query.selectedQuery');
  const selectedQuery = selectedQueryID && queryEntities.findBy({ id: selectedQueryID });
  const selectedPacks = selectedQuery && packEntities.where({ query_id: selectedQuery.id });

  return { queries, selectedQuery, selectedPacks };
};

export default connect(mapStateToProps)(ManageQueriesPage);

