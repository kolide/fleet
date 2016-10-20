import React, { Component, PropTypes } from 'react';
import radium from 'radium';

import HostSidePanelStyles from './styles';
import { iconClassForLabel } from './helpers';

class PanelGroupItem extends Component {
  static propTypes = {
    item: PropTypes.shape({
      count: PropTypes.number,
      label: PropTypes.string,
      name: PropTypes.string,
    }).isRequired,
    onLabelClick: PropTypes.func,
    selected: PropTypes.bool,
  };

  render () {
    const { item, onLabelClick, selected } = this.props;
    const {
      count,
      label,
    } = item;
    const {
      PanelGroupItemStyles: {
        containerStyles,
        itemStyles,
      },
    } = HostSidePanelStyles;

    return (
      <div onClick={onLabelClick} style={containerStyles(selected)}>
        <div style={[itemStyles, { width: '41px' }]}>
          <i className={iconClassForLabel(item)} />
        </div>
        <div style={[itemStyles, { width: '160px' }]}>{label}</div>
        <div style={[itemStyles, { width: '35px', textAlign: 'right' }]}>{count}</div>
      </div>
    );
  }
}

export default radium(PanelGroupItem);
