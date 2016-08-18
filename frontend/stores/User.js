/*
 * Copyright 2016-present, Kolide, Inc.
 * All rights reserved.
 *
 * @flow
 */

import { Store } from 'nuclear-js';
import * as Immutable from 'immutable';

/**
 * UserStore is the application store for the currently logged in user.
 *
 * @exports UserStore
 */
var UserStore = new Store({
  /**
   * Return the initial state of the store on application launch.
   */
  getInitialState(): Immutable.Map<string, any> {
    return Immutable.Map({
      logged_in: false,
      is_authenticating: false,
    });
  },

  /**
   * Initialize the store by declaring what dispatcher routes this store has
   * subscribed to
   */
  initialize() {
    this.on('RECEIVE_USER_INFO', this.receiveUserInfo);
    this.on('IS_AUTHENTICATING', this.setAuthenticatingState);
  },

  /**
   * Handler for the RECEIVE_USER_INFO dispatcher event
   *
   * @param {Immutable.Map<string, any>} state The current store state
   * @param {Immutable.Map<string, any>} user User data to be received
   */
  receiveUserInfo(state: Immutable.Map<string, any>, user: Immutable.Map<string, any>): Immutable.Map<string, any> {
    return state.merge(Immutable.Map({
      id: user.get('id'),
      username: user.get('username'),
      email: user.get('email'),
      name: user.get('name'),
      admin: user.get('admin'),
      needs_password_reset: user.get('needs_password_reset'),
      logged_in: true,
    }));
  },

  /**
   * Handler for the IS_AUTHENTICATING dispatcher event
   *
   * @param {Immutable.Map<string, any>} state The current store state
   * @param {Immutable.Map<string, any>} isAuthenticating The auth state
   */
  setAuthenticatingState(state: Immutable.Map<string, any>, isAuthenticating: Immutable.Map<string, any>): Immutable.Map<string, any> {
    return state.merge(Immutable.Map({
      is_authenticating: isAuthenticating,
    }));
  },
});

/**
 * UserGetters exposes the getters which can be used to subscribe to changes
 * in the UserStore
 *
 * @exports UserGetters
 */
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