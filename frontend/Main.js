import React from 'react';
import { render } from 'react-dom';
import { Router, Route, IndexRedirect, browserHistory, withRouter } from 'react-router';
import { Promise } from 'bluebird';

import Dispatcher from '#app/Dispatcher';
import { requireAuthentication } from '#app/Authentication';

import Infrastructure from '#components/Infrastructure';
import Login from '#components/Login'

import AppActions from '#actions/App';
import UserStore from '#stores/User';


if (typeof window !== 'undefined') {
  window.Promise = window.Promise || Promise;
  window.self = window;
  require('whatwg-fetch');

  require('#css');
  if (module.hot) {
    let c = 0;
    module.hot.accept('#css', () => {
      require('#css');
      const a = document.createElement('a');
      const link = document.querySelector('link[rel="stylesheet"]');
      a.href = link.href;
      a.search = '?' + c++;
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


