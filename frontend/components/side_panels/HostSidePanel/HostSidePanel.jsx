import React, { Component, PropTypes } from 'react';
import radium from 'radium';

import componentStyles from './styles';
import InputField from '../../forms/fields/InputField';
import PanelGroup from './PanelGroup';
import SecondarySidePanelContainer from '../SecondarySidePanelContainer';

class HostSidePanel extends Component {
  static propTypes = {
    allHostGroupItems: PropTypes.array,
    hostPlatformGroupItems: PropTypes.array,
    hostStatusGroupItems: PropTypes.array,
    onLabelClick: PropTypes.func,
    selectedLabel: PropTypes.object,
  };

  render () {
    const {
      allHostGroupItems,
      hostPlatformGroupItems,
      hostStatusGroupItems,
      onLabelClick,
      selectedLabel,
    } = this.props;
    const {
      containerStyles,
      hrStyles,
      PanelGroupItemStyles,
    } = componentStyles;

    return (
      <SecondarySidePanelContainer style={containerStyles}>
        <PanelGroup
          groupItems={allHostGroupItems}
          onLabelClick={onLabelClick}
          selectedLabel={selectedLabel}
        />
        <hr style={hrStyles} />
        <PanelGroup
          groupItems={hostStatusGroupItems}
          onLabelClick={onLabelClick}
          selectedLabel={selectedLabel}
        />
        <hr style={hrStyles} />
        <PanelGroup
          groupItems={hostPlatformGroupItems}
          onLabelClick={onLabelClick}
          selectedLabel={selectedLabel}
        />
        <hr style={hrStyles} />
        <div style={PanelGroupItemStyles.containerStyles(false)}>
          <i className="kolidecon-tag" />
          <span style={{ marginLeft: '20px' }}>LABELS</span>
        </div>
        <div style={PanelGroupItemStyles.containerStyles(false)}>
          <InputField
            name="tags-filter"
            placeholder="Filter by Name..."
            style={{ width: '100%' }}
          />
        </div>
      </SecondarySidePanelContainer>
    );
  }
}

export default radium(HostSidePanel);
