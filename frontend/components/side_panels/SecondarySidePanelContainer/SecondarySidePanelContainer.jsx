import React, { Component, PropTypes } from 'react';
import radium from 'radium';

import componentStyles from './styles';

class SecondarySidePanelContainer extends Component {
  static propTypes = {
    children: PropTypes.node,
    style: PropTypes.object, // eslint-disable-line react/forbid-prop-types
  };

  render () {
    const { children, style } = this.props;
    const { containerStyles } = componentStyles;

    return (
      <div style={[containerStyles, style]}>
        {children}
      </div>
    );
  }
}

export default radium(SecondarySidePanelContainer);
