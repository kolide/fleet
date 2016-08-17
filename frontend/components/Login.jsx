import React from 'react';
import Dispatcher from '#app/Dispatcher';

import { UserGetters } from '#stores/User';

import AppActions from '#actions/App';

const Login = React.createClass({
  mixins: [
    Dispatcher.ReactMixin
  ],

  getDataBindings() {
    return {
      logged_in: UserGetters.logged_in,
      is_authenticating: UserGetters.is_authenticating,
    }
  },

  getInitialState() {
    return {
      username: "",
      password: ""
    }
  },

  componentWillMount() {
    if (this.state.logged_in) {
      this.props.router.push("/");
    }
  },

  onClickLoginButton(event) {
    event.preventDefault();
    AppActions.login(this.state.username, this.state.password, "foo");
  },

  onChangeUsername(event) {
    this.state.username = event.target.value;
  },

  onChangePassword(event) {
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
})

Login.getRoute = function() {
  return "/login";
}

export default Login