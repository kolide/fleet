import React, { Component } from 'react';
import NewQuery from '../../../components/Queries/NewQuery';

class NewQueryPage extends Component {
  static propTypes = {
    children: PropTypes.node,
  };

  render () {
    return (
      <div>
        <h1>New Query Page</h1>
      </div>
    )
  }
}

export default NewQueryPage;

