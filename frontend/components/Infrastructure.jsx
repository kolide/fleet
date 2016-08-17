/*
 * Copyright 2016-present, Kolide, Inc.
 * All rights reserved.
 *
 * @flow
 */

import React, { Component } from 'react';
import Dispatcher from 'frontend/Dispatcher';

import Page from 'frontend/components/Page';

/**
 * Infrastructure is the main page component for viewing infrastructure
 *
 * @exports Infrastructure
 */
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