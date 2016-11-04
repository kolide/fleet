import React, { Component, PropTypes } from 'react';

import AllPacks from '../../../components/packs/AllPacks';
import PackInfoSidePanel from '../../../components/side_panels/PackInfoSidePanel';

class AllPacksPage extends Component {
  static propTypes = {
    children: PropTypes.element,
  }

  render () {
    return (
      <div>
        <AllPacks />
        <PackInfoSidePanel />
      </div>
    );
  }
}

export default AllPacksPage;
