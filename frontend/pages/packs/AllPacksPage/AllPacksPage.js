import React, { Component, PropTypes } from 'react';

class AllPacksPage extends Component {
  static propTypes = {
    children: PropTypes.node,
  };

  render () {
    const { children } = this.props;

    return (
      <div>
        <h1>All Packs</h1>
        {children}
      </div>
    );
  }
}

export default AllPacksPage;
