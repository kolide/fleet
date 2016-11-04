import React, { Component, PropTypes } from 'react';

import NewPack from 'components/packs/NewPack';
import PackInfoSidePanel from 'components/side_panels/PackInfoSidePanel';

class NewPackPage extends Component {
  static propTypes = {
    children: PropTypes.element,
  }

  render () {
    return (
      <div>
        <NewPack />
        <PackInfoSidePanel />
      </div>
    );
  }
}

export default NewPackPage;
