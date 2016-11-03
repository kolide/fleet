import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

import entityGetter from 'redux/utilities/entityGetter';
import queryActions from 'redux/nodes/entities/queries/actions';

class QueryPageWrapper extends Component {
  static propTypes = {
    children: PropTypes.node,
  };

  componentWillMount () {
    const { dispatch, query, queryID } = this.props;

    console.log('here', queryID);
    if (queryID && !query) {
      dispatch(queryActions.loadOne(queryID));
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

export default QueryPageWrapper;
