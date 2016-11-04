import React, { Component, PropTypes } from 'react';

const baseClass = 'new-pack';

class AllPacks extends Component {
  static propTypes = {
    children: PropTypes.element,
  }

  render () {
    return (
      <div className={`${baseClass}__wrapper`}>
        <p className={`${baseClass}__title`}>
          Query Pack Title
        </p>
      </div>
    );
  }
}

export default AllPacks;
