import React, { Component, PropTypes } from 'react';
import { size } from 'lodash';

import helpers from 'components/queries/QueriesListWrapper/helpers';
import InputField from 'components/forms/fields/InputField';
import PackQueryConfigForm from 'components/forms/packs/PackQueryConfigForm';
import QueriesList from 'components/queries/QueriesList';
import queryInterface from 'interfaces/query';

const baseClass = 'queries-list-wrapper';

class QueriesListWrapper extends Component {
  static propTypes = {
    configuredQueryIDs: PropTypes.arrayOf(PropTypes.number),
    onClearStagedQueries: PropTypes.func,
    onConfigureQueries: PropTypes.func,
    onDeselectQuery: PropTypes.func,
    onSelectQuery: PropTypes.func,
    queries: PropTypes.arrayOf(queryInterface),
    stagedQueries: PropTypes.arrayOf(queryInterface),
  };

  constructor (props) {
    super(props);

    this.state = {
      querySearchText: '',
      selectAll: false,
    };
  }

  onSelectQuery = (query) => {
    return (shouldCheck) => {
      const { onDeselectQuery, onSelectQuery } = this.props;

      if (shouldCheck) {
        onSelectQuery(query);

        return false;
      }

      onDeselectQuery(query);

      return false;
    };
  }

  onUpdateQuerySearchText = (querySearchText) => {
    this.setState({ querySearchText });
  }

  getQueries = () => {
    const { queries } = this.props;
    const { querySearchText } = this.state;

    return helpers.filterQueries(queries, querySearchText);
  }

  renderPackQueryConfigForm = () => {
    const {
      onClearStagedQueries,
      onConfigureQueries,
      stagedQueries,
    } = this.props;

    const formData = { queries: stagedQueries };

    return (
      <PackQueryConfigForm
        formData={formData}
        handleSubmit={onConfigureQueries}
        onCancel={onClearStagedQueries}
      />
    );
  }

  render () {
    const {
      getQueries,
      onSelectQuery,
      onUpdateQuerySearchText,
      renderPackQueryConfigForm,
    } = this;
    const { configuredQueryIDs, queries, stagedQueries } = this.props;
    const { querySearchText } = this.state;
    const queryCount = size(queries);

    if (!queryCount) {
      return false;
    }

    return (
      <div className={`${baseClass} ${baseClass}__wrapper`}>
        <p>Add Queries to Pack</p>
        <InputField
          className={`${baseClass}__search-queries-input`}
          name="search-queries"
          onChange={onUpdateQuerySearchText}
          placeholder="Search Queries"
          value={querySearchText}
        />
        <div style={{ position: 'relative' }}>
          {renderPackQueryConfigForm()}
          <QueriesList
            configuredQueryIDs={configuredQueryIDs}
            onSelectQuery={onSelectQuery}
            queries={getQueries()}
            selectedQueries={stagedQueries}
          />
        </div>
      </div>
    );
  }
}

export default QueriesListWrapper;
