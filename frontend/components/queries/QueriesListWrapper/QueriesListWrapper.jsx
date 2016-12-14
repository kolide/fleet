import React, { Component, PropTypes } from 'react';
import { pull } from 'lodash';

import Button from 'components/buttons/Button';
import helpers from 'components/queries/QueriesListWrapper/helpers';
import InputField from 'components/forms/fields/InputField';
import QueriesList from 'components/queries/QueriesList';
import queryInterface from 'interfaces/query';

const baseClass = 'queries-list-wrapper';

class QueriesListWrapper extends Component {
  static propTypes = {
    allQueries: PropTypes.arrayOf(queryInterface),
    scheduledQueries: PropTypes.arrayOf(queryInterface),
  };

  constructor (props) {
    super(props);

    const queryCount = props.scheduledQueries.length;

    this.state = {
      querySearchText: '',
      selectAll: false,
      selectedQueries: [],
      shouldShowPackForm: queryCount === 0,
    };
  }

  onHidePackForm = () => {
    this.setState({ shouldShowPackForm: false });

    return false;
  }

  onSelectQuery = (query) => {
    return (shouldAddQuery) => {
      const { selectedQueries } = this.state;
      const newSelectedQueries = shouldAddQuery ?
        selectedQueries.concat(query) :
        pull(selectedQueries, query);

      this.setState({ selectedQueries: newSelectedQueries });

      return false;
    };
  }

  onShowPackForm = (evt) => {
    evt.preventDefault();

    this.setState({ shouldShowPackForm: true });

    return false;
  }

  onUpdateQuerySearchText = (querySearchText) => {
    this.setState({ querySearchText });
  }

  getQueries = () => {
    const { scheduledQueries } = this.props;
    const { querySearchText } = this.state;

    return helpers.filterQueries(scheduledQueries, querySearchText);
  }

  renderQueryCount = () => {
    const { scheduledQueries } = this.props;
    const queryCount = scheduledQueries.length;
    const queryText = queryCount === 1 ? 'Query' : 'Queries';

    return <p className={`${baseClass}__query-count`}><span>{queryCount}</span> {queryText}</p>;
  }

  renderQueriesList = () => {
    const { getQueries, onHidePackForm, onSelectQuery } = this;
    const { selectedQueries, shouldShowPackForm } = this.state;

    return (
      <div className={`${baseClass}__queries-list-wrapper`}>
        <QueriesList
          onHidePackForm={onHidePackForm}
          onSelectQuery={onSelectQuery}
          queries={getQueries()}
          selectedQueries={selectedQueries}
          shouldShowPackForm={shouldShowPackForm}
        />
      </div>
    );
  }

  render () {
    const { onShowPackForm, onUpdateQuerySearchText, renderQueryCount, renderQueriesList } = this;
    const { querySearchText, shouldShowPackForm } = this.state;

    return (
      <div className={`${baseClass} ${baseClass}__wrapper`}>
        {renderQueryCount()}
        <div>
          <InputField
            inputClassName={`${baseClass}__search-queries-input`}
            name="search-queries"
            onChange={onUpdateQuerySearchText}
            placeholder="Search Queries"
            value={querySearchText}
          />
          <Button
            className={`${baseClass}__add-new-query-btn`}
            disabled={shouldShowPackForm}
            onClick={onShowPackForm}
            text="Add New Query"
            variant="brand"
          />
        </div>
        {renderQueriesList()}
      </div>
    );
  }
}

export default QueriesListWrapper;
