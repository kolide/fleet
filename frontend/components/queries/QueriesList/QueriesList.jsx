import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
import { includes, isEqual, size } from 'lodash';

import Icon from 'components/Icon';
import PackQueryConfigForm from 'components/forms/packs/PackQueryConfigForm';
import queryInterface from 'interfaces/query';
import QueriesListItem from 'components/queries/QueriesList/QueriesListItem';

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

    this.state = { queryDropdownOptions: [] };
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
        <td colSpan={8}>
          <PackQueryConfigForm
            handleSubmit={handleSubmit}
            onCancel={onHidePackForm}
            queryOptions={queryDropdownOptions}
          />
        </td>
      </tr>
    );
  }

  renderHelpText = () => {
    const { isScheduledQueriesAvailable, scheduledQueries } = this.props;

    if (scheduledQueries.length) {
      return false;
    }

    if (isScheduledQueriesAvailable) {
      return (
        <tr>
          <td colSpan={8}>
            <p>No queries matched your search criteria.</p>
          </td>
        </tr>
      );
    }

    return (
      <tr>
        <td colSpan={8}>
          <h1>First let&apos;s <span>add a query</span></h1>
          <h2>Then we&apos;ll add the following:</h2>
          <p><b>interval:</b> the amount of time the query waits before running</p>
          <p><b>minimum <Icon name="osquery" /> version:</b> the minimum required <b>osqueryd</b> version installed on a host</p>
          <p><b><Icon name="add-plus" /> differential:</b> show only what&apos; different from last run</p>
          <p><b><Icon name="camera" /> snapshot:</b> show everything in its current state</p>
        </td>
      </tr>
    );
  }

  render () {
    const { onSelectQuery, scheduledQueries, selectedScheduledQueryIDs, shouldShowPackForm } = this.props;
    const { renderHelpText, renderPackQueryConfigForm } = this;

    const wrapperClassName = classnames(`${baseClass}__wrapper`, {
      [`${baseClass}__wrapper--query-selected`]: size(selectedScheduledQueryIDs),
    });

    return (
      <table className={wrapperClassName}>
        <thead>
          <tr>
            <td />
            <td>Query Name</td>
            <td>Interval</td>
            <td>Platform</td>
            <td>Version</td>
            <td>Logging Type</td>
            <td>Author</td>
            <td>Last Modified</td>
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
    );
  }
}

export default QueriesList;
