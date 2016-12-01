import React, { Component } from 'react';
import { keys, some, values } from 'lodash';

import campaignInterface from 'interfaces/campaign';

const baseClass = 'query-results-table';

class QueryResultsTable extends Component {
  static propTypes = {
    campaign: campaignInterface.isRequired,
  };

  renderTableHeaderRow = () => {
    const { campaign } = this.props;
    const { query_results: queryResults } = campaign;

    const { rows } = queryResults[0];
    const queryResultColumns = keys(rows[0]);

    return (
      <tr>
        <th>hostname</th>
        {queryResultColumns.map((column) => {
          return <th key={column}>{column}</th>;
        })}
      </tr>
    );
  }

  renderTableRows = () => {
    const { campaign } = this.props;
    const { query_results: queryResults } = campaign;

    if (!queryResults) {
      return false;
    }


    return queryResults.map((result) => {
      const { host, rows } = result;

      return rows.map((row) => {
        const rowResults = values(row);

        return (
          <tr>
            <td>{host.hostname}</td>
            {rowResults.map((rowData, i) => {
              return <td key={`query-results-table-row-${i}`}>{rowData}</td>;
            })}
          </tr>
        );
      });
    });
  }

  render () {
    const { campaign } = this.props;
    const { renderTableRows, renderTableHeaderRow } = this;
    const { query_results: queryResults } = campaign;

    if (!queryResults) {
      return false;
    }

    const rowsPresent = some(queryResults, (result) => {
      return result.rows.length;
    });

    if (!rowsPresent) {
      return false;
    }

    return (
      <div className={`${baseClass} ${baseClass}__wrapper`}>
        <table className={`${baseClass}__table`}>
          <thead>
            {renderTableHeaderRow()}
          </thead>
          <tbody>
            {renderTableRows()}
          </tbody>
        </table>
      </div>
    );
  }
}

export default QueryResultsTable;
