/*
 * Copyright 2016-present, Kolide, Inc.
 * All rights reserved.
 *
 * @flow
 */

import React from 'react';
import Dispatcher from 'frontend/Dispatcher';

import { UserGetters } from 'frontend/stores/User';

import AppActions from 'frontend/actions/App';

/**
 * Login is the main page component for the login page.
 *
 * @exports Login
 */
const Login = React.createClass({
  mixins: [
    Dispatcher.ReactMixin
  ],

  getDataBindings() {
    return {
      logged_in: UserGetters.logged_in,
      is_authenticating: UserGetters.is_authenticating,
    };
  },

  getInitialState() {
    return {
      logged_in: null,
      is_authenticating: null,
      username: '',
      password: ''
    };
  },

  componentWillMount() {
    if (this.state.logged_in) {
      this.props.router.push('/');
    }
  },

  onClickLoginButton(event: any) {
    event.preventDefault();
    AppActions.login(this.state.username, this.state.password, 'foo');
  },

  onChangeUsername(event: any) {
    this.state.username = event.target.value;
  },

  onChangePassword(event: any) {
    this.state.password = event.target.value;
  },

  render() {
    return (
      <div className="Login">
        <h1> Login </h1>
        <form role='form'>
          <div>
            <input type='text'
                  className='input-lg'
                  placeholder='Username' 
                  onChange={this.onChangeUsername} />
          </div>
          <div>
            <input type='password'
                  className='input-lg'
                  placeholder='Password'
                  onChange={this.onChangePassword} />
          </div>
          <button type='submit'
                  className='btn btn-lg'
                  disabled={this.state.is_authenticating}
                  onClick={this.onClickLoginButton}>Submit</button>
        </form>
      </div>
    );
  },
});

Login.getRoute = function() {
  return '/login';
};

export default Login;