/**
 * @flow
 */

import { Store } from 'nuclear-js';
import * as Immutable from 'immutable';

var UserStore = new Store({
  getInitialState() {
    return Immutable.Map({
      logged_in: false,
      is_authenticating: false,
    });
  },

  initialize() {
    this.on('RECEIVE_USER_INFO', receiveUserInfo);
    this.on('IS_AUTHENTICATING', setAuthenticatingState);
  }

});

function receiveUserInfo(state: Immutable.Map<string, any>, user: Immutable.Map<string, any>) {
  return state.merge(Immutable.Map({
    id: user.get("id"),
    username: user.get("username"),
    email: user.get("email"),
    name: user.get("name"),
    admin: user.get("admin"),
    needs_password_reset: user.get("needs_password_reset"),
    logged_in: true,
  }));
}

function setAuthenticatingState(state: Immutable.Map<string, any>, isAuthenticating: Immutable.Map<string, any>) {
  return state.merge(Immutable.Map({
    is_authenticating: isAuthenticating,
  }));
}

var UserGetters = {
  id: ['user', 'id'],
  username: ['user', 'username'],
  email: ['user', 'email'],
  name: ['user', 'name'],
  admin: ['user', 'admin'],
  needs_password_reset: ['user', 'needs_password_reset'],
  logged_in: ['user', 'logged_in'],
  is_authenticating: ['user', 'is_authenticating'],
};


export default UserStore;
export { UserGetters };