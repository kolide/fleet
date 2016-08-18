/*
 * Copyright 2016-present, Kolide, Inc.
 * All rights reserved.
 *
 * @flow
 */

import React from 'react';
import Dispatcher from 'frontend/Dispatcher';

import { UserGetters } from 'frontend/stores/User';

import Login from 'frontend/components/Login';

/**
 * requireAuthentication wraps a component and ensure that the current user is
 * logged in before rendering the component
 *
 * @example
 * import { requireAuthentication } from 'frontend/Authentication';
 * // in react router:
 * <Route path="/" component={requireAuthentication(FooComponent)}></Route>
 *
 * @param {React.Component} higherOrderComponent The component to wrap
 * @return {React.Component}
 * @exports requireAuthentication
 */
function requireAuthentication(higherOrderComponent: $FlowFixMe) {
  return React.createClass({
    mixins: [Dispatcher.ReactMixin],

    getDataBindings() {
      return {
        logged_in: UserGetters.logged_in,
      }
    },

    getInitialState() {
      return {
        logged_in: null,
      }
    },

    componentWillMount() {
      this.checkAuth();
    },

    componentWillReceiveProps(nextProps: $FlowFixMe) {
      this.checkAuth();
    },

    checkAuth() {
      if (!this.state.logged_in) {
        // TODO: set redirect state so user can come back here after logging in
        this.props.router.push(Login.getRoute());
      }
    },

    render() {
      var Component = higherOrderComponent;
      return (
        <div>
          {this.state.logged_in === true
            ? <Component {...this.props}/> 
            : null
          }
        </div>
      )
    },
  })
}

exports.requireAuthentication = requireAuthentication