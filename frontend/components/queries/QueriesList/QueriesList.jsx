import React, { Component, PropTypes } from 'react';
import { includes } from 'lodash';

import queryInterface from 'interfaces/query';
import QueriesListItem from 'components/queries/QueriesList/QueriesListItem';

class QueriesList extends Component {
  static propTypes = {
    configuredQueries: PropTypes.arrayOf(queryInterface),
    onSelectQuery: PropTypes.func.isRequired,
    queries: PropTypes.arrayOf(queryInterface).isRequired,
    selectedQueries: PropTypes.arrayOf(queryInterface),
  };

  render () {
    const {
      configuredQueries,
      onSelectQuery,
      queries,
      selectedQueries,
    } = this.props;

    return (
      <table>
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
          {queries.map((query) => {
            return (
              <QueriesListItem
                checked={includes(selectedQueries, query)}
                configured={includes(configuredQueries, query)}
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
