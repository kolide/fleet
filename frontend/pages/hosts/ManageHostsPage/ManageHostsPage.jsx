import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { filter } from 'lodash';

import componentStyles from './styles';
import entityGetter from '../../../redux/utilities/entityGetter';
import hostActions from '../../../redux/nodes/entities/hosts/actions';
import labelActions from '../../../redux/nodes/entities/labels/actions';
import labelInterface from '../../../interfaces/label';
import HostDetails from '../../../components/hosts/HostDetails';
import hostInterface from '../../../interfaces/host';
import HostSidePanel from '../../../components/side_panels/HostSidePanel';
import { setSelectedLabel } from '../../../redux/nodes/components/ManageHostsPage/actions';
import { showRightSidePanel, removeRightSidePanel } from '../../../redux/nodes/app/actions';

class ManageHostsPage extends Component {
  static propTypes = {
    allHostLabels: PropTypes.arrayOf(labelInterface),
    dispatch: PropTypes.func,
    hosts: PropTypes.arrayOf(hostInterface),
    hostPlatformLabels: PropTypes.arrayOf(labelInterface),
    hostStatusLabels: PropTypes.arrayOf(labelInterface),
    labels: PropTypes.arrayOf(labelInterface),
    selectedLabel: labelInterface,
  };

  componentWillMount () {
    const {
      allHostLabels,
      dispatch,
      hosts,
      labels,
      selectedLabel,
    } = this.props;

    dispatch(showRightSidePanel);

    if (!hosts.length) {
      dispatch(hostActions.loadAll());
    }

    if (!labels.length) {
      dispatch(labelActions.loadAll());
    }

    if (!selectedLabel) {
      dispatch(setSelectedLabel(allHostLabels[0]));
    }

    return false;
  }

  componentWillReceiveProps (nextProps) {
    const { allHostLabels, dispatch, selectedLabel } = nextProps;
    const allHostLabel = allHostLabels[0];

    if (!selectedLabel && !!allHostLabel) {
      dispatch(setSelectedLabel(allHostLabel));
    }

    return false;
  }

  componentWillUnmount () {
    const { dispatch } = this.props;

    dispatch(removeRightSidePanel);
  }

  onHostDetailActionClick = (type) => {
    return (host) => {
      return (evt) => {
        evt.preventDefault();

        console.log(type, host);
        return false;
      };
    };
  }

  onLabelClick = (selectedLabel) => {
    return (evt) => {
      evt.preventDefault();

      const { dispatch } = this.props;

      dispatch(setSelectedLabel(selectedLabel));

      return false;
    };
  }

  renderHeader = () => {
    const { selectedLabel } = this.props;
    const {
      headerHostsCountStyles,
      headerHostsTitleStyles,
      headerStyles,
    } = componentStyles;

    if (!selectedLabel) return false;

    return (
      <div style={headerStyles}>
        <i className="kolidecon-tag" />
        <span style={headerHostsTitleStyles}>{selectedLabel.label}</span>
        <span style={headerHostsCountStyles}>{selectedLabel.count} Hosts Total</span>
      </div>
    );
  }

  renderHosts = () => {
    const { hosts } = this.props;
    const { onHostDetailActionClick } = this;

    return hosts.map((host) => {
      return (
        <HostDetails
          host={host}
          key={host.hostname}
          onDisableClick={onHostDetailActionClick('disable')}
          onQueryClick={onHostDetailActionClick('query')}
        />
      );
    });
  }

  render () {
    const { containerStyles } = componentStyles;
    const {
      allHostLabels,
      hostPlatformLabels,
      hostStatusLabels,
      selectedLabel,
    } = this.props;
    const { onLabelClick, renderHeader, renderHosts } = this;

    return (
      <div style={containerStyles}>
        {renderHeader()}
        {renderHosts()}
        <HostSidePanel
          allHostGroupItems={allHostLabels}
          hostPlatformGroupItems={hostPlatformLabels}
          hostStatusGroupItems={hostStatusLabels}
          onLabelClick={onLabelClick}
          selectedLabel={selectedLabel}
        />
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: hosts } = entityGetter(state).get('hosts');
  const { entities: labels } = entityGetter(state).get('labels');
  const allHostLabels = filter(labels, { type: 'all' });
  const hostStatusLabels = filter(labels, { type: 'status' });
  const hostPlatformLabels = filter(labels, { type: 'platform' });
  const { selectedLabel } = state.components.ManageHostsPage;

  return {
    allHostLabels,
    hosts,
    labels,
    hostStatusLabels,
    hostPlatformLabels,
    selectedLabel,
  };
};

export default connect(mapStateToProps)(ManageHostsPage);
