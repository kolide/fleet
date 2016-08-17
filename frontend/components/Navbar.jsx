import React from 'react';
import { Link } from 'react-router';
import Dispatcher from '#app/Dispatcher';

import { UserGetters } from '#stores/User';

import Infrastructure from '#components/Infrastructure';

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