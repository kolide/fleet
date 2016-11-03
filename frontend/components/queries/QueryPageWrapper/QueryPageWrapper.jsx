import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';

import entityGetter from 'redux/utilities/entityGetter';
import queryActions from 'redux/nodes/entities/queries/actions';
import queryInterface from 'interfaces/query';
import { renderFlash } from 'redux/nodes/notifications/actions';

class QueryPageWrapper extends Component {
  static propTypes = {
    children: PropTypes.node,
    dispatch: PropTypes.func,
    query: queryInterface,
    queryID: PropTypes.string,
  };

  componentDidMount () {
    const { dispatch, query, queryID } = this.props;

    if (queryID && !query) {
      dispatch(queryActions.load(queryID))
        .catch((errorResponse) => {
          const { error } = errorResponse;
          dispatch(push('/queries/new'));
          dispatch(renderFlash('error', error));
        });
    }

    return false;
  }

  render () {
    const { children } = this.props;

    return (
      <div>
        {children}
      </div>
    );
  }
}

const mapStateToProps = (state, { params }) => {
  const { id: queryID } = params;
  const query = entityGetter(state).get('queries').findBy({ id: queryID });

  return { query, queryID };
};

export default connect(mapStateToProps)(QueryPageWrapper);
