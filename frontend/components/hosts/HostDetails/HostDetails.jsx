import React, { PropTypes } from 'react';
import radium from 'radium';

import componentStyles from './styles';
import { humanMemory, humanUptime } from './helpers';

const {
  containerStyles,
  contentSeparatorStyles,
  hostContentItemStyles,
  hostnameStyles,
  statusStyles,
} = componentStyles;
export const STATUSES = {
  ONLINE: 'ONLINE',
  OFFLINE: 'OFFLINE',
  UPGRADE: 'NEEDS_UPGRADE',
};

const HostDetails = ({ host }) => {
  const status = STATUSES.ONLINE;
  const { hostname, ip, mac, memory, platform, uptime } = host;

  return (
    <div style={containerStyles(status)}>
      <div style={statusStyles(status)}>
        {status}
      </div>
      <p style={hostnameStyles}>{hostname}</p>
      <div style={contentSeparatorStyles}>
        <div>
          <span style={[hostContentItemStyles, { textTransform: 'capitalize' }]}>{platform}</span>
        </div>
        <div>
          <span style={hostContentItemStyles}>{platform}</span>
        </div>
        <div>
          <span style={hostContentItemStyles}>{humanMemory(memory)}</span>
          <span style={hostContentItemStyles}>{humanUptime(uptime)}</span>
        </div>
        <div>
          <span style={hostContentItemStyles}>{mac}</span>
        </div>
        <div>
          <span style={hostContentItemStyles}>{ip}</span>
        </div>
      </div>
      <div style={contentSeparatorStyles}>
        <div>
          <span style={[hostContentItemStyles, { textTransform: 'capitalize' }]}>Tags go here</span>
        </div>
      </div>
    </div>
  );
};

HostDetails.propTypes = {
  host: PropTypes.shape({
    hostname: PropTypes.string,
    ip: PropTypes.string,
    mac: PropTypes.string,
    memory: PropTypes.number,
    platform: PropTypes.string,
    uptime: PropTypes.number,
  }).isRequired,
};

export default radium(HostDetails);
