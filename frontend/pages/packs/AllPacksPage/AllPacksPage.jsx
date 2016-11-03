import React, { Component, PropTypes } from 'react';

import PackInfoSidePanel from '../../../components/side_panels/PackInfoSidePanel';

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
        <PackInfoSidePanel />
      </div>
    );
  }
}

export default AllPacksPage;
