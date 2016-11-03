import React, { Component, PropTypes } from 'react';

class NewPackPage extends Component {
  static propTypes = {
    children: PropTypes.node,
  };

  render () {
    const { children } = this.props;

    return (
      <div>
        <h1>New Pack</h1>
        {children}
      </div>
    );
  }
}

export default NewPackPage;
