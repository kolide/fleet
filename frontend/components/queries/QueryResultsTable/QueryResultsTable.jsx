import React, { Component } from 'react';
import { get, keys, omit, values } from 'lodash';

import campaignInterface from 'interfaces/campaign';
import ProgressBar from 'components/ProgressBar';

const baseClass = 'query-results-table';

class QueryResultsTable extends Component {
  static propTypes = {
    campaign: campaignInterface.isRequired,
  };

  renderProgressDetails = () => {
    const { campaign } = this.props;
    const totalHostsCount = get(campaign, 'totals.count', 0);
    const totalHostsReturned = get(campaign, 'hosts.length', 0);
    const totalRowsCount = get(campaign, 'query_results.length', 0);

    return (
      <div className={`${baseClass}__progress-details`}>
        <span>
          <b>{totalHostsReturned}</b>&nbsp;of&nbsp;
          <b>{totalHostsCount} Hosts</b>&nbsp;Returning&nbsp;
          <b>{totalRowsCount} Records</b>
        </span>
        <ProgressBar max={totalHostsCount} value={totalHostsReturned} />
      </div>
    );
  }

  renderTableHeaderRow = () => {
    const { campaign } = this.props;
    const { query_results: queryResults } = campaign;

    const queryAttrs = omit(queryResults[0], ['hostname']);
    const queryResultColumns = keys(queryAttrs);

    return (
      <tr>
        <th>host</th>
        {queryResultColumns.map((column) => {
          return <th key={column}>{column}</th>;
        })}
      </tr>
    );
  }

  renderTableRows = () => {
    const { campaign } = this.props;
    const { query_results: queryResults } = campaign;

    return queryResults.map((row) => {
      const queryAttrs = omit(row, ['hostname']);
      const queryResult = values(queryAttrs);

      return (
        <tr>
          <td>{row.hostname}</td>
          {queryResult.map((attribute, i) => {
            return <td key={`query-results-table-row-${i}`}>{attribute}</td>;
          })}
        </tr>
      );
    });
  }

  render () {
    const { campaign } = this.props;
    const {
      renderProgressDetails,
      renderTableHeaderRow,
      renderTableRows,
    } = this;
    const { query_results: queryResults } = campaign;

    if (!queryResults || !queryResults.length) {
      return false;
    }

    return (
      <div className={baseClass}>
        {renderProgressDetails()}
        <div className={`${baseClass}__table-wrapper`}>
          <table className={`${baseClass}__table`}>
            <thead>
              {renderTableHeaderRow()}
            </thead>
            <tbody>
              {renderTableRows()}
            </tbody>
          </table>
        </div>
      </div>
    );
  }
}

export default QueryResultsTable;
