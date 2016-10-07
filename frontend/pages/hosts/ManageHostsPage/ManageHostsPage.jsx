import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

import entityGetter from '../../../redux/entityGetter';
import hostActions from '../../../redux/nodes/entities/hosts/actions';

class ManageHostsPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    hosts: PropTypes.array,
  };

  componentWillMount () {
    const { dispatch, hosts } = this.props;

    if (!hosts.length) dispatch(hostActions.loadAll());

    return false;
  }

  render () {
    return (
      <h1>Manage Hosts Page</h1>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: hosts } = entityGetter(state).get('hosts');

  return { hosts };
};

export default connect(mapStateToProps)(ManageHostsPage);
