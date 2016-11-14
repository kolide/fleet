import React, { Component, PropTypes } from 'react';

import InputField from 'components/forms/fields/InputField';
import labelInterface from 'interfaces/label';
import PanelGroup from 'components/side_panels/HostSidePanel/PanelGroup';
import SecondarySidePanelContainer from 'components/side_panels/SecondarySidePanelContainer';

const baseClass = 'host-side-panel';

class HostSidePanel extends Component {
  static propTypes = {
    allHostGroupItems: PropTypes.arrayOf(labelInterface),
    hostPlatformGroupItems: PropTypes.arrayOf(labelInterface),
    hostStatusGroupItems: PropTypes.arrayOf(labelInterface),
    onAddLabelClick: PropTypes.func,
    onLabelClick: PropTypes.func,
    selectedLabel: labelInterface,
  };

  render () {
    const {
      allHostGroupItems,
      hostPlatformGroupItems,
      hostStatusGroupItems,
      onAddLabelClick,
      onLabelClick,
      selectedLabel,
    } = this.props;

    return (
      <SecondarySidePanelContainer className={`${baseClass}__wrapper`}>
        <PanelGroup
          groupItems={allHostGroupItems}
          onLabelClick={onLabelClick}
          selectedLabel={selectedLabel}
        />
        <hr className={`${baseClass}__hr`} />
        <PanelGroup
          groupItems={hostStatusGroupItems}
          onLabelClick={onLabelClick}
          selectedLabel={selectedLabel}
        />
        <hr className={`${baseClass}__hr`} />
        <PanelGroup
          groupItems={hostPlatformGroupItems}
          onLabelClick={onLabelClick}
          selectedLabel={selectedLabel}
        />
        <hr className={`${baseClass}__hr`} />
        <div className={`${baseClass}__panel-group-item`}>
          <i className="kolidecon-tag" />
          <span className="title">LABELS</span>
        </div>
        <div className={`${baseClass}__panel-group-item`}>
          <InputField
            name="tags-filter"
            placeholder="Filter by Name..."
          />
        </div>
        <hr className={`${baseClass}__hr`} />
        <button className={`${baseClass}__add-label-btn button button--unstyled`} onClick={onAddLabelClick}>
          <i className="kolidecon-add-button" />
          ADD NEW LABEL
          <i className="kolidecon-label" />
        </button>
      </SecondarySidePanelContainer>
    );
  }
}

export default HostSidePanel;
