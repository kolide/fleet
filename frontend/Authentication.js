import React from 'react';
import Dispatcher from '#app/Dispatcher';

import { UserGetters } from '#stores/User';

import Login from '#components/Login';

function requireAuthentication(higherOrderComponent) {
  return React.createClass({
    mixins: [Dispatcher.ReactMixin],

    getDataBindings() {
      return {
        logged_in: UserGetters.logged_in,
      }
    },

    componentWillMount() {
      this.checkAuth();
    },

    componentWillReceiveProps(nextProps) {
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