import React, { PropTypes } from 'react';

import Icon from 'components/icons/Icon';
import { Link } from 'react-router';
import packInterface from 'interfaces/pack';
import ScheduledQueriesSection from 'components/side_panels/PackDetailsSidePanel/ScheduledQueriesSection';
import scheduledQueryInterface from 'interfaces/scheduled_query';
import SecondarySidePanelContainer from 'components/side_panels/SecondarySidePanelContainer';
import Slider from 'components/forms/fields/Slider';

const baseClass = 'pack-details-side-panel';

const Description = ({ pack }) => {
  if (!pack.description) {
    return false;
  }

  return (
    <div>
      <p>Description</p>
      <p>{pack.description}</p>
    </div>
  );
};

const PackDetailsSidePanel = ({ onUpdateSelectedPack, pack, scheduledQueries = [] }) => {
  const { disabled } = pack;
  const updatePackStatus = (value) => {
    return onUpdateSelectedPack(pack, { disabled: !value });
  };

  return (
    <SecondarySidePanelContainer className={baseClass}>
      <Icon name="packs" /><span>{pack.name}</span>
      <Slider
        activeText="ENABLED"
        inactiveText="DISABLED"
        label="Status"
        onChange={updatePackStatus}
        value={!disabled}
      />
      <Link className="button button--inverse" to={`/packs/${pack.id}`}>Edit Pack</Link>
      <Description pack={pack} />
      <ScheduledQueriesSection scheduledQueries={scheduledQueries} />
    </SecondarySidePanelContainer>
  );
};

Description.propTypes = {
  pack: packInterface.isRequired,
};

PackDetailsSidePanel.propTypes = {
  onUpdateSelectedPack: PropTypes.func,
  pack: packInterface.isRequired,
  scheduledQueries: PropTypes.arrayOf(scheduledQueryInterface),
};

export default PackDetailsSidePanel;

