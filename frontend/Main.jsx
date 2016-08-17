/*
 * Copyright 2016-present, Kolide, Inc.
 * All rights reserved.
 *
 * @flow
 */

import React from 'react';
import { render } from 'react-dom';
import { Router, Route, IndexRedirect, browserHistory, withRouter } from 'react-router';
import { Promise } from 'bluebird';

import Dispatcher from 'frontend/Dispatcher';
import { requireAuthentication } from 'frontend/Authentication';

import Infrastructure from 'frontend/components/Infrastructure';
import Login from 'frontend/components/Login'

import AppActions from 'frontend/actions/App';
import UserStore from 'frontend/stores/User';


if (typeof window !== 'undefined') {
  window.Promise = window.Promise || Promise;
  window.self = window;
  require('whatwg-fetch');

  require('frontend/css');
  if (module.hot) {
    let c = 0;
    module.hot.accept('#css', () => {
      require('frontend/css');

      var a = document.createElement('a');
      var link = document.querySelector('link[rel="stylesheet"]');

      // @FlowFixMe clean up style loading
      a.href = link.href;
      a.search = '?' + c++;
      
      // @FlowFixMe clean up style loading
      link.href = a.href;
    });
  }

  Dispatcher.registerStores({
    user: UserStore,
  })

  AppActions.fetchInitialState();

  render((
    <Router history={browserHistory}>
      <Route path="/">
        <IndexRedirect to={Infrastructure.getRoute()} />
        <Route path={Login.getRoute()} component={withRouter(Login)}></Route>
        <Route path={Infrastructure.getRoute()} component={withRouter(requireAuthentication(Infrastructure))}></Route>
      </Route>
    </Router>
  ), document.getElementById('app'))
}


