/*
 * Copyright 2016-present, Kolide, Inc.
 * All rights reserved.
 *
 * @flow
 */

import Dispatcher from 'frontend/Dispatcher';
import { browserHistory } from 'react-router';
import * as Immutable from 'immutable';


/**
 * AppActions is a top-level set of actions that corresponds, generally, to
 * events that are happening at the scope of the entire application.
 *
 * @exports AppActions
 */
module.exports = {
  /**
   * AppActions.fetchInitialState is called when the application starts. All
   * state initialization should occur here.
   */
  fetchInitialState: function() {
  },

  /**
   * AppActions.login is the action which is called when a user attempts to
   * login to the application
   *
   * @param {string} email The email address of the user
   * @param {string} password The password of the user
   * @param {string} redirectTo The URL to redirect to after login
   */
  login(email: string, password: string, redirectTo: string) {
    Dispatcher.dispatch('IS_AUTHENTICATING', true);

    function sleep(time) {
      return new Promise((resolve) => setTimeout(resolve, time));
    }

    var p = sleep(1000);

    p.then(() => {
      var user = Immutable.Map({
        id: 'marpaia',
        username: 'marpaia',
        email: 'mike@kolide.co',
        name: 'Mike Arpaia',
        admin: true,
        needs_password_reset: false
      });

      Dispatcher.dispatch('RECEIVE_USER_INFO', user);
      browserHistory.push('/');
    });

    p.then(() => {
      Dispatcher.dispatch('IS_AUTHENTICATING', false);
    });
  },
};