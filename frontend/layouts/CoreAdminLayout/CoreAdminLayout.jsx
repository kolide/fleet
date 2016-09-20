import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import paths from '../../router/paths';

export class CoreAdminLayout extends Component {
  static propTypes = {
    children: PropTypes.node,
    dispatch: PropTypes.func,
    user: PropTypes.object,
  };

  componentWillMount () {
    const { dispatch, user: { admin } } = this.props;
    const { HOME } = paths;

    if (!admin) dispatch(push(HOME));
  }

  render () {
    const { children } = this.props;

    return (
      <div>
        <div>SidePanel</div>
        <div>{children}</div>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { user } = state.auth;

  return { user };
};

export default connect(mapStateToProps)(CoreAdminLayout);
