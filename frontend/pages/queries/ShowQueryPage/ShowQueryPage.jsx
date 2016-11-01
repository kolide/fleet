import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

import entityGetter from '../../../redux/utilities/entityGetter';
import queryInterface from '../../../interfaces/query';
import osqueryTableInterface from '../../../interfaces/osquery_table';

class ShowQueryPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    query: queryInterface,
    selectedOsqueryTable: osqueryTableInterface,
  };

  render () {
    return (
      <div>
        <h1>Show Query Page</h1>
      </div>
    );
  }
}

const mapStateToProps = (state, { params }) => {
  const { id: queryID } = params;
  const query = entityGetter(state).get('queries').findBy({ id: queryID });
  const { selectedOsqueryTable } = state.app;

  return { query, selectedOsqueryTable };
};

export default connect(mapStateToProps)(ShowQueryPage);
