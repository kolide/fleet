import { PropTypes } from 'react';
import { connect } from 'react-redux';

import entityGetter from '../../../redux/utilities/entityGetter';
import queryInterface from '../../../interfaces/query';
import osqueryTableInterface from '../../../interfaces/osquery_table';
import NewQueryPage from '../NewQueryPage';

class ShowQueryPage extends NewQueryPage {
  static propTypes = {
    dispatch: PropTypes.func,
    query: queryInterface,
    selectedOsqueryTable: osqueryTableInterface,
  };
}

const mapStateToProps = (state, { params }) => {
  const { id: queryID } = params;
  const query = entityGetter(state).get('queries').findBy({ id: queryID });
  const { selectedOsqueryTable } = state.app;

  return { query, selectedOsqueryTable };
};

export default connect(mapStateToProps)(ShowQueryPage);
