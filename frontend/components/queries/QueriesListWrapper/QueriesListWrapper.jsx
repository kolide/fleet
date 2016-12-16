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
    onRemoveScheduledQueries: PropTypes.func,
    onScheduledQueryFormSubmit: PropTypes.func,
    scheduledQueries: PropTypes.arrayOf(queryInterface),
  };

  constructor (props) {
    super(props);

    this.state = {
      querySearchText: '',
      selectAll: false,
      selectedScheduledQueryIDs: [],
      shouldShowPackForm: false,
    };
  }

  onHidePackForm = () => {
    this.setState({ shouldShowPackForm: false });

    return false;
  }

  onRemoveScheduledQueries = (evt) => {
    evt.preventDefault();

    const { onRemoveScheduledQueries: handleRemoveScheduledQueries } = this.props;
    const { selectedScheduledQueryIDs } = this.state;

    this.setState({ selectedScheduledQueryIDs: [] });

    return handleRemoveScheduledQueries(selectedScheduledQueryIDs);
  }

  onSelectQuery = (shouldAddQuery, scheduledQueryID) => {
    const { selectedScheduledQueryIDs } = this.state;
    const newSelectedScheduledQueryIDs = shouldAddQuery ?
      selectedScheduledQueryIDs.concat(scheduledQueryID) :
      pull(selectedScheduledQueryIDs, scheduledQueryID);

    this.setState({ selectedScheduledQueryIDs: newSelectedScheduledQueryIDs });

    return false;
  }

  onShowPackForm = (evt) => {
    evt.preventDefault();

    this.setState({
      selectedScheduledQueryIDs: [],
      shouldShowPackForm: true,
    });

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

  renderButton = () => {
    const { onRemoveScheduledQueries, onShowPackForm } = this;
    const { selectedScheduledQueryIDs, shouldShowPackForm } = this.state;

    const scheduledQueryCount = selectedScheduledQueryIDs.length;

    if (scheduledQueryCount) {
      const queryText = scheduledQueryCount === 1 ? 'Query' : 'Queries';

      return (
        <Button
          className={`${baseClass}__query-btn`}
          onClick={onRemoveScheduledQueries}
          text={`Remove ${queryText}`}
          variant="alert"
        />
      );
    }

    return (
      <Button
        className={`${baseClass}__query-btn`}
        disabled={shouldShowPackForm}
        onClick={onShowPackForm}
        text="Add New Query"
        variant="brand"
      />
    );
  }

  renderQueryCount = () => {
    const { scheduledQueries } = this.props;
    const queryCount = scheduledQueries.length;
    const queryText = queryCount === 1 ? 'Query' : 'Queries';

    return <h1 className={`${baseClass}__query-count`}><span>{queryCount}</span> {queryText}</h1>;
  }

  renderQueriesList = () => {
    const { getQueries, onHidePackForm, onSelectQuery } = this;
    const { allQueries, onScheduledQueryFormSubmit, scheduledQueries } = this.props;
    const { selectedScheduledQueryIDs, shouldShowPackForm } = this.state;

    return (
      <div className={`${baseClass}__queries-list-wrapper`}>
        <QueriesList
          allQueries={allQueries}
          onHidePackForm={onHidePackForm}
          onScheduledQueryFormSubmit={onScheduledQueryFormSubmit}
          onSelectQuery={onSelectQuery}
          scheduledQueries={getQueries()}
          selectedScheduledQueryIDs={selectedScheduledQueryIDs}
          shouldShowPackForm={shouldShowPackForm}
          isScheduledQueriesAvailable={!!scheduledQueries.length}
        />
      </div>
    );
  }

  render () {
    const { onUpdateQuerySearchText, renderButton, renderQueryCount, renderQueriesList } = this;
    const { querySearchText } = this.state;

    return (
      <div className={baseClass}>
        {renderQueryCount()}
        <InputField
          inputClassName={`${baseClass}__search-queries-input`}
          inputWrapperClass={`${baseClass}__search-queries`}
          name="search-queries"
          onChange={onUpdateQuerySearchText}
          placeholder="Search Queries"
          value={querySearchText}
        />
        {renderButton()}
        {renderQueriesList()}
      </div>
    );
  }
}

export default QueriesListWrapper;
