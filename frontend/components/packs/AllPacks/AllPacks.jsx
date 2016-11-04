import React, { Component, PropTypes } from 'react';

const baseClass = 'all-packs';

class AllPacks extends Component {
  static propTypes = {
    children: PropTypes.element,
  }

  render () {
    return (
      <div className={`${baseClass}__wrapper`}>
        <p className={`${baseClass}__title`}>
          Query Packs
        </p>
      </div>
    );
  }
}

export default AllPacks;
