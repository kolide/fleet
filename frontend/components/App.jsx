import React, { Component } from 'react';
import Dispatcher from '#app/Dispatcher'

import { AppGetters } from '#stores/App';

export default React.createClass({
  mixins: [Dispatcher.ReactMixin],

  getDataBindings() {
    return {
      username: AppGetters.username,
    }
  },

  render() {
    return (
      <div className="App">
        <h1>Kolide</h1>
        <p>If you can read this, React is rendering correctly!</p>
        <p>{this.state.username}</p>
      </div>
    );
  },
})