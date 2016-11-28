import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
import { includes, size } from 'lodash';

import queryInterface from 'interfaces/query';
import QueriesListItem from 'components/queries/QueriesList/QueriesListItem';

const baseClass = 'queries-list';

class QueriesList extends Component {
  static propTypes = {
    configuredQueryIDs: PropTypes.arrayOf(PropTypes.number),
    onSelectQuery: PropTypes.func.isRequired,
    queries: PropTypes.arrayOf(queryInterface).isRequired,
    selectedQueries: PropTypes.arrayOf(queryInterface),
  };

  render () {
    const {
      configuredQueryIDs,
      onSelectQuery,
      queries,
      selectedQueries,
    } = this.props;

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
          {queries.map((query) => {
            return (
              <QueriesListItem
                checked={includes(selectedQueries, query)}
                configured={includes(configuredQueryIDs, query.id)}
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
