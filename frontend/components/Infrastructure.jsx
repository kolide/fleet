/**
 * @flow
 */

import React, { Component } from 'react';
import Dispatcher from 'frontend/Dispatcher';

import Page from 'frontend/components/Page';

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