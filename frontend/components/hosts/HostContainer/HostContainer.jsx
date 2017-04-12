import React, { Component, PropTypes } from 'react';
import { isEqual, orderBy, sortBy } from 'lodash';
// import Pagination from 'rc-pagination';

import hostInterface from 'interfaces/host';
import labelInterface from 'interfaces/label';
import HostsTable from 'components/hosts/HostsTable';
import HostDetails from 'components/hosts/HostDetails';
import LonelyHost from 'components/hosts/LonelyHost';
import Spinner from 'components/loaders/Spinner';
import helpers from './helpers';

const baseClass = 'host-container';
// const PAGE_SIZE = 2;

// const en_US = {
//   // Options.jsx
//   items_per_page: '/ page',
//   jump_to: 'Goto',
//   page: '',

//   // Pagination.jsx
//   prev_page: 'Previous Page',
//   next_page: 'Next Page',
//   prev_5: 'Previous 5 Pages',
//   next_5: 'Next 5 Pages',
//   prev_3: 'Previous 3 Pages',
//   next_3: 'Next 3 Pages',
// };

class HostContainer extends Component {
  static propTypes = {
    hosts: PropTypes.arrayOf(hostInterface),
    selectedLabel: labelInterface,
    loadingHosts: PropTypes.bool.isRequired,
    displayType: PropTypes.oneOf(['Grid', 'List']),
    toggleAddHostModal: PropTypes.func,
    toggleDeleteHostModal: PropTypes.func,
    onQueryHost: PropTypes.func,
  };

  constructor (props) {
    super(props);

    this.state = {
      allHostCount: 0,
      currentPagination: 0,
      sortedHosts: [],
    };
  }

  componentWillReceiveProps (nextProps) {
    if (isEqual(nextProps, this.props)) {
      return false;
    }

    const { filterHosts, sortHosts } = this;

    const filteredHosts = filterHosts();
    const sortedHosts = sortHosts(filteredHosts);

    this.setState({
      allHostCount: filteredHosts.length,
      sortedHosts,
    });

    return true;
  }

  shouldComponentUpdate (nextProps, nextState) {
    if (isEqual(nextProps, this.props) &&
        isEqual(nextState, this.state)) {
      return false;
    }

    return true;
  }

  filterHosts = () => {
    const { hosts, selectedLabel } = this.props;

    return helpers.filterHosts(hosts, selectedLabel);
  }

  sortHosts = (hosts) => {
    const alphaHosts = sortBy(hosts, (h) => { return h.hostname; });
    const orderedHosts = orderBy(alphaHosts, 'status', 'desc');

    return orderedHosts;
  }

  // handlePaginationChange = () => {
  //   console.log('clicky clicky');
  // }

  renderNoHosts = () => {
    const { selectedLabel } = this.props;
    const { type } = selectedLabel || '';
    const isCustom = type === 'custom';

    return (
      <div className={`${baseClass}__no-hosts`}>
        <h1>No matching hosts found.</h1>
        <h2>Where are the missing hosts?</h2>
        <ul>
          {isCustom && <li>Check your SQL query above to confirm there are no mistakes.</li>}
          <li>Check to confirm that your hosts are online.</li>
          <li>Confirm that your expected hosts have osqueryd installed and configured.</li>
        </ul>

        <div className={`${baseClass}__no-hosts-contact`}>
          <p>Still having trouble? Want to talk to a human?</p>
          <p>Contact Kolide Support:</p>
          <p><a href="mailto:support@kolide.co">support@kolide.co</a></p>
        </div>
      </div>
    );
  }

  renderHosts = () => {
    const { displayType, toggleDeleteHostModal, onQueryHost } = this.props;
    const { sortedHosts } = this.state;

    if (displayType === 'Grid') {
      return sortedHosts.map((host) => {
        const isLoading = !host.hostname;

        return (
          <HostDetails
            host={host}
            key={`host-${host.id}-details`}
            onDestroyHost={toggleDeleteHostModal}
            onQueryHost={onQueryHost}
            isLoading={isLoading}
          />
        );
      });
    } else {
      return (
        <HostsTable
          hosts={sortedHosts}
          onDestroyHost={toggleDeleteHostModal}
          onQueryHost={onQueryHost}
        />
      );
    }
  }

  render () {
    const { renderHosts, renderNoHosts } = this;
    const { allHostCount } = this.state;
    const { displayType, loadingHosts, selectedLabel, toggleAddHostModal } = this.props;

    if (loadingHosts) {
      return <Spinner />;
    }

    if (allHostCount === 0) {
      if (selectedLabel && selectedLabel.type === 'all') {
        return <LonelyHost onClick={toggleAddHostModal} />;
      }

      return renderNoHosts();
    }

    return(
      <div className={`${baseClass} ${baseClass}--${displayType.toLowerCase()}`}>
        {renderHosts()}
      </div>
    );
  }
}

export default HostContainer;
