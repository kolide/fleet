import React from 'react';
import { render } from 'react-dom';
import { Router, Route, IndexRoute, browserHistory } from 'react-router';
import { Promise } from 'bluebird';

import Dispatcher from '#app/Dispatcher';

import App from '#components/App';
import AppActions from '#actions/App';
import AppStore from '#stores/App';


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
    app: AppStore,
  })

  AppActions.fetchInitialState();

  render((
    <Router history={browserHistory}>
      <Route path="/" component={App}></Route>
    </Router>
  ), document.getElementById('app'))
}


