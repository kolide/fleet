import React, { PropTypes } from 'react';

import Icon from 'components/icons/Icon';
import { Link } from 'react-router';
import packInterface from 'interfaces/pack';
import SecondarySidePanelContainer from 'components/side_panels/SecondarySidePanelContainer';
import Slider from 'components/forms/fields/Slider';

const baseClass = 'pack-details-side-panel';

const PackDetailsSidePanel = ({ onUpdateSelectedPack, pack }) => {
  const { disabled } = pack;
  const updatePackStatus = (value) => {
    return onUpdateSelectedPack(pack, { disabled: !value });
  };

  return (
    <SecondarySidePanelContainer className={baseClass}>
      <div>
        <Icon name="packs" /><span>{pack.name}</span>
        <Slider
          activeText="ENABLED"
          inactiveText="DISABLED"
          label="Status"
          onChange={updatePackStatus}
          value={!disabled}
        />
        <Link className="button button--inverse" to={`/packs/${pack.id}`}>Edit Pack</Link>
      </div>
    </SecondarySidePanelContainer>
  );
};

PackDetailsSidePanel.propTypes = {
  onUpdateSelectedPack: PropTypes.func,
  pack: packInterface.isRequired,
};

export default PackDetailsSidePanel;

