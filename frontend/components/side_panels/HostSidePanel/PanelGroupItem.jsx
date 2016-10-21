import React, { Component, PropTypes } from 'react';
import radium from 'radium';

import HostSidePanelStyles from './styles';
import { iconClassForLabel } from './helpers';
import labelInterface from '../../../interfaces/label';

class PanelGroupItem extends Component {
  static propTypes = {
    item: labelInterface.isRequired,
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
      <button className="btn--unstyled" onClick={onLabelClick} style={containerStyles(selected)}>
        <div style={[itemStyles, { width: '41px' }]}>
          <i className={iconClassForLabel(item)} />
        </div>
        <div style={[itemStyles, { width: '160px' }]}>{label}</div>
        <div style={[itemStyles, { width: '35px', textAlign: 'right' }]}>{count}</div>
      </button>
    );
  }
}

export default radium(PanelGroupItem);
