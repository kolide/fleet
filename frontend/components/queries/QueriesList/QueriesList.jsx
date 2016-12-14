import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
import { includes, size } from 'lodash';

import PackQueryConfigForm from 'components/forms/packs/PackQueryConfigForm';
import queryInterface from 'interfaces/query';
import QueriesListItem from 'components/queries/QueriesList/QueriesListItem';

const baseClass = 'queries-list';

class QueriesList extends Component {
  static propTypes = {
    allQueries: PropTypes.arrayOf(queryInterface).isRequired,
    onHidePackForm: PropTypes.func.isRequired,
    onScheduledQueryFormSubmit: PropTypes.func,
    onSelectQuery: PropTypes.func.isRequired,
    scheduledQueries: PropTypes.arrayOf(queryInterface).isRequired,
    selectedQueries: PropTypes.arrayOf(queryInterface).isRequired,
    shouldShowPackForm: PropTypes.bool,
  };

  renderPackQueryConfigForm = () => {
    const { onHidePackForm, onScheduledQueryFormSubmit, allQueries, shouldShowPackForm } = this.props;

    if (!shouldShowPackForm) {
      return false;
    }

    // TODO: Move this to state
    const queryOptions = allQueries.map((query) => {
      return { label: query.name, value: String(query.id) };
    });

    return (
      <tr>
        <td colSpan={6}>
          <PackQueryConfigForm
            handleSubmit={onScheduledQueryFormSubmit}
            onCancel={onHidePackForm}
            queryOptions={queryOptions}
          />
        </td>
      </tr>
    );
  }

  render () {
    const { onSelectQuery, scheduledQueries, selectedQueries, shouldShowPackForm } = this.props;
    const { renderPackQueryConfigForm } = this;

    const wrapperClassName = classnames(`${baseClass}__wrapper`, {
      [`${baseClass}__wrapper--query-selected`]: size(selectedQueries),
    });

    return (
      <table className={wrapperClassName}>
        <thead>
          <tr>
            <td />
            <td>Query Name</td>
            <td>Description</td>
            <td>Platform</td>
            <td>Author</td>
            <td>Last Modified</td>
          </tr>
        </thead>
        <tbody>
          {renderPackQueryConfigForm()}
          {scheduledQueries.map((query) => {
            return (
              <QueriesListItem
                checked={includes(selectedQueries, query)}
                disabled={shouldShowPackForm}
                key={`query-${query.id}`}
                onSelect={onSelectQuery(query)}
                query={query}
              />
            );
          })}
        </tbody>
      </table>
    );
  }
}

export default QueriesList;
