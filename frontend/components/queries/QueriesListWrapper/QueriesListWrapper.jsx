import React, { Component, PropTypes } from 'react';
import { remove, size } from 'lodash';

import Button from 'components/buttons/Button';
import InputField from 'components/forms/fields/InputField';
import QueriesList from 'components/queries/QueriesList';
import queryInterface from 'interfaces/query';

const baseClass = 'queries-list-wrapper';

class QueriesListWrapper extends Component {
  static propTypes = {
    queries: PropTypes.arrayOf(queryInterface),
  };

  constructor (props) {
    super(props);

    this.state = {
      querySearchText: '',
      selectAll: false,
      selectedQueries: [],
    };
  }

  onSelectQuery = (query) => {
    return (shouldCheck) => {
      const { selectedQueries } = this.state;

      if (shouldCheck) {
        this.setState({
          selectedQueries: selectedQueries.concat(query),
        });

        return false;
      }

      remove(selectedQueries, query);

      this.setState({ selectedQueries });

      return false;
    };
  }

  onUpdateQuerySearchText = (querySearchText) => {
    this.setState({ querySearchText });
  }

  render () {
    const { onSelectQuery, onUpdateQuerySearchText } = this;
    const { queries } = this.props;
    const { querySearchText, selectedQueries } = this.state;
    const queryCount = size(queries);

    if (!queryCount) {
      return false;
    }

    const addQueryBtnText = (
      <span className={`${baseClass}__add-query-btn-text`}>
        <i className={`${baseClass}__add-query-btn-icon kolidecon-add-button`} />
        Add New Query
      </span>
    );

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
        <Button
          text={addQueryBtnText}
          variant="brand"
        />
        <QueriesList
          onSelectQuery={onSelectQuery}
          queries={queries}
          selectedQueries={selectedQueries}
        />
      </div>
    );
  }
}

export default QueriesListWrapper;
