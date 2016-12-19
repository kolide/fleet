import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
import { includes, isEqual, size } from 'lodash';

import Icon from 'components/Icon';
import PackQueryConfigForm from 'components/forms/packs/PackQueryConfigForm';
import queryInterface from 'interfaces/query';
import QueriesListItem from 'components/queries/QueriesList/QueriesListItem';
import Checkbox from 'components/forms/fields/Checkbox';

const baseClass = 'queries-list';

class QueriesList extends Component {
  static propTypes = {
    allQueries: PropTypes.arrayOf(queryInterface).isRequired,
    isScheduledQueriesAvailable: PropTypes.bool,
    onHidePackForm: PropTypes.func.isRequired,
    onScheduledQueryFormSubmit: PropTypes.func,
    onSelectQuery: PropTypes.func.isRequired,
    scheduledQueries: PropTypes.arrayOf(queryInterface).isRequired,
    selectedScheduledQueryIDs: PropTypes.arrayOf(PropTypes.number).isRequired,
    shouldShowPackForm: PropTypes.bool,
  };

  constructor (props) {
    super(props);

    this.state = {
      queryDropdownOptions: [],
      allQueriesSelected: false,
    };
  }

  componentWillMount () {
    const { allQueries } = this.props;

    const queryDropdownOptions = allQueries.map((query) => {
      return { label: query.name, value: String(query.id) };
    });

    this.setState({ queryDropdownOptions });
  }

  componentWillReceiveProps (nextProps) {
    const { allQueries } = nextProps;

    if (!isEqual(allQueries, this.props.allQueries)) {
      const queryDropdownOptions = allQueries.map((query) => {
        return { label: query.name, value: String(query.id) };
      });

      this.setState({ queryDropdownOptions });
    }
  }

  handleSubmit = (formData) => {
    const { onHidePackForm, onScheduledQueryFormSubmit } = this.props;

    onHidePackForm();

    return onScheduledQueryFormSubmit(formData);
  }

  isChecked = (scheduledQuery) => {
    const { selectedScheduledQueryIDs } = this.props;

    return includes(selectedScheduledQueryIDs, scheduledQuery.id);
  }

  renderPackQueryConfigForm = () => {
    const { onHidePackForm, shouldShowPackForm } = this.props;
    const { queryDropdownOptions } = this.state;

    if (!shouldShowPackForm) {
      return false;
    }

    const { handleSubmit } = this;

    return (
      <tr>
        <td colSpan={6}>
          <PackQueryConfigForm
            handleSubmit={handleSubmit}
            onCancel={onHidePackForm}
            queryOptions={queryDropdownOptions}
          />
        </td>
      </tr>
    );
  }

  handleSelectAllQueries = () => {
    const { allQueriesSelected } = this.state;

    this.setState({
      allQueriesSelected: !allQueriesSelected,
    })
  }

  renderHelpText = () => {
    const { isScheduledQueriesAvailable, scheduledQueries } = this.props;

    if (scheduledQueries.length) {
      return false;
    }

    if (isScheduledQueriesAvailable) {
      return (
        <tr>
          <td colSpan={6}>
            <p>No queries matched your search criteria.</p>
          </td>
        </tr>
      );
    }

    return (
      <tr>
        <td colSpan={6}>
          <div className={`${baseClass}__first-query`}>
            <h1>First let's <span>add a query</span>.</h1>
            <h2>Then we'll set the following:</h2>
            <p><strong>interval:</strong> the amount of time the query waits before running</p>
            <p><strong>minimum <Icon name="osquery" /> version:</strong> the minimum required <strong>osqueryd</strong> version installed on a host</p>
            <p><strong>logging type:</strong></p>
            <ul>
              <li><strong><Icon name="plus-minus" /> differential:</strong> show only what's different from last run</li>
              <li><strong><Icon name="camera" /> snapshot:</strong> show everything in its current state</li>
            </ul>
          </div>
        </td>
      </tr>
    );
  }

  render () {
    const { onSelectQuery, scheduledQueries, selectedScheduledQueryIDs, shouldShowPackForm } = this.props;
    const { allQueriesSelected } = this.state;
    const { renderHelpText, renderPackQueryConfigForm, handleSelectAllQueries } = this;

    const wrapperClassName = classnames(`${baseClass}__table`, {
      [`${baseClass}__table--query-selected`]: size(selectedScheduledQueryIDs),
    });

    return (
      <div className={baseClass}>
        <table className={wrapperClassName}>
          <thead>
            <tr>
              <th><Checkbox
                name="select-all-queries"
                onChange={handleSelectAllQueries}
                value={allQueriesSelected}
              /></th>
              <th>Query Name</th>
              <th>Interval [s]</th>
              <th>Platform</th>
              <th><Icon name="osquery" /> Ver.</th>
              <th>Log</th>
            </tr>
          </thead>
          <tbody>
            {renderPackQueryConfigForm()}
            {renderHelpText()}
            {!!scheduledQueries.length && scheduledQueries.map((scheduledQuery) => {
              return (
                <QueriesListItem
                  checked={this.isChecked(scheduledQuery)}
                  disabled={shouldShowPackForm}
                  key={`scheduled-query-${scheduledQuery.id}`}
                  onSelect={onSelectQuery}
                  scheduledQuery={scheduledQuery}
                />
              );
            })}
          </tbody>
        </table>
      </div>
    );
  }
}

export default QueriesList;
