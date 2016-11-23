import React, { Component, PropTypes } from 'react';
import { size } from 'lodash';

import InputField from 'components/forms/fields/InputField';
import PackQueryConfigForm from 'components/forms/packs/PackQueryConfigForm';
import QueriesList from 'components/queries/QueriesList';
import queryInterface from 'interfaces/query';

const baseClass = 'queries-list-wrapper';

class QueriesListWrapper extends Component {
  static propTypes = {
    configuredQueryIDs: PropTypes.arrayOf(PropTypes.number),
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

  renderPackQueryConfigForm = () => {
    const { onConfigureQueries, stagedQueries } = this.props;

    if (!size(stagedQueries)) {
      return false;
    }

    const formData = { queries: stagedQueries };

    return (
      <PackQueryConfigForm
        formData={formData}
        handleSubmit={onConfigureQueries}
      />
    );
  }

  render () {
    const {
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
        {renderPackQueryConfigForm()}
        <QueriesList
          configuredQueryIDs={configuredQueryIDs}
          onSelectQuery={onSelectQuery}
          queries={queries}
          selectedQueries={stagedQueries}
        />
      </div>
    );
  }
}

export default QueriesListWrapper;
