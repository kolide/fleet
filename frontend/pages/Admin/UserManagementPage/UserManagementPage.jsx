import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import entityGetter from '../../../redux/entityGetter';
import userActions from '../../../redux/nodes/entities/users/actions';

class UserManagementPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    users: PropTypes.arrayOf(PropTypes.object),
  };

  componentWillMount () {
    const { dispatch } = this.props;
    const { load } = userActions;

    dispatch(load());
  }

  render () {
    return (
      <div>
        <h1>User Management Page</h1>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: users } = entityGetter(state).get('users');

  return { users };
};

export default connect(mapStateToProps)(UserManagementPage);

