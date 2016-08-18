/*
 * Copyright 2016-present, Kolide, Inc.
 * All rights reserved.
 *
 * @flow
 */

import React from 'react';
import { Link } from 'react-router';
import Dispatcher from 'frontend/Dispatcher';

import { UserGetters } from 'frontend/stores/User';

import Infrastructure from 'frontend/components/Infrastructure';

/**
 * Navbar is the react component for the application navbar.
 *
 * @exports Navbar
 */
const Navbar = React.createClass({
  mixins: [
    Dispatcher.ReactMixin
  ],

  getDataBindings() {
    return {
      name: UserGetters.name,
    }
  },

  render() {
    return (
      <div className="Navbar">
        <h3> Hello, {this.state.name}! </h3>
        <ul>
          <li><Link to={Infrastructure.getRoute()}>Infrastructure</Link></li>
        </ul>
      </div>
    );
  },
})

export default Navbar