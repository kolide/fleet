import React, { Component } from 'react';
import Dispatcher from '#app/Dispatcher';

import Page from '#components/Page';

const Infrastructure = React.createClass({
  mixins: [
    Dispatcher.ReactMixin
  ],

  getDataBindings() {
    return {
    }
  },

  render() {
    return (
      <Page className="Infrastructure">
        <h1> Infrastructure </h1>
      </Page>
    );
  },
})

Infrastructure.getRoute = function() {
    return "/infrastructure";
}

export default Infrastructure