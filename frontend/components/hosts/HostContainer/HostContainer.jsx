import React, { Component, PropTypes } from 'react';
import { isEqual, orderBy, slice, sortBy } from 'lodash';
import Pagination from 'rc-pagination';
import Select from 'react-select';
import 'rc-pagination/assets/index.css';

import enUs from 'rc-pagination/lib/locale/en_US';
import hostInterface from 'interfaces/host';
import labelInterface from 'interfaces/label';
import HostsTable from 'components/hosts/HostsTable';
import HostDetails from 'components/hosts/HostDetails';
import LonelyHost from 'components/hosts/LonelyHost';
import Spinner from 'components/loaders/Spinner';
import helpers from './helpers';


const baseClass = 'host-container';

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
      currentPage: 0,
      hostsPerPage: 20,
      pagedHosts: [],
      showSpinner: false,
    };
  }

  componentWillMount () {
    this.buildSortedHosts();
  }

  componentWillUpdate (nextProps, nextState) {
    if (isEqual(nextProps, this.props) && isEqual(nextState, this.state)) {
      return false;
    }

    this.buildSortedHosts(nextProps, nextState);
    return true;
  }

  buildSortedHosts = (nextProps, nextState) => {
    const { filterHosts, sortHosts } = this;
    const { currentPage, hostsPerPage } = nextState || this.state;
    const { hosts, selectedLabel } = nextProps || this.props;

    const sortedHosts = sortHosts(filterHosts(hosts, selectedLabel));

    const fromIndex = currentPage * hostsPerPage;
    const toIndex = fromIndex + hostsPerPage;

    const pagedHosts = slice(sortedHosts, fromIndex, toIndex);

    this.setState({
      allHostCount: sortedHosts.length,
      pagedHosts,
      showSpinner: false,
    });
  }

  filterHosts = (hosts, selectedLabel) => {
    const { filterHosts } = helpers;

    return filterHosts(hosts, selectedLabel);
  }

  sortHosts = (hosts) => {
    const alphaHosts = sortBy(hosts, (h) => { return h.hostname; });
    const orderedHosts = orderBy(alphaHosts, 'status', 'desc');

    return orderedHosts;
  }

  handlePaginationChange = (page) => {
    const { scrollToTop } = helpers;

    this.setState({
      currentPage: page - 1,
    });

    scrollToTop();

    return true;
  }

  handlePerPageChange = (option) => {
    const { scrollToTop } = helpers;

    this.setState({
      currentPage: 0,
      hostsPerPage: Number(option.value),
      showSpinner: true,
    });

    scrollToTop();

    return true;
  }

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
    const { pagedHosts } = this.state;

    if (displayType === 'Grid') {
      return pagedHosts.map((host) => {
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
    }

    return (
      <HostsTable
        hosts={pagedHosts}
        onDestroyHost={toggleDeleteHostModal}
        onQueryHost={onQueryHost}
      />
    );
  }

  renderPagination = () => {
    const { handlePaginationChange, handlePerPageChange } = this;
    const { allHostCount, currentPage, hostsPerPage } = this.state;

    const paginationSelectOpts = [
      { value: 20, label: '20' },
      { value: 100, label: '100' },
      { value: 500, label: '500' },
      { value: 1000, label: '1,000' },
    ];

    const humanPage = currentPage + 1;
    const startRange = (currentPage * hostsPerPage) + 1;
    const endRange = Math.min(humanPage * hostsPerPage, allHostCount);

    return (
      <div className={`${baseClass}__pager-wrap`}>
        <Pagination
          onChange={handlePaginationChange}
          current={humanPage}
          total={allHostCount}
          pageSize={hostsPerPage}
          className={`${baseClass}__pagination`}
          locale={enUs}
          showLessItems
        />
        <p className={`${baseClass}__pager-range`}>{`${startRange} - ${endRange} of ${allHostCount} hosts`}</p>
        <div className={`${baseClass}__pager-count`}>
          <Select
            name="pager-host-count"
            value={hostsPerPage}
            options={paginationSelectOpts}
            onChange={handlePerPageChange}
            className={`${baseClass}__count-select`}
            clearable={false}
          /> <span>Hosts per page</span>
        </div>
      </div>
    );
  }

  render () {
    const { renderHosts, renderNoHosts, renderPagination } = this;
    const { allHostCount, showSpinner } = this.state;
    const { displayType, loadingHosts, selectedLabel, toggleAddHostModal } = this.props;

    if (loadingHosts || showSpinner) {
      return <Spinner />;
    }

    if (allHostCount === 0) {
      if (selectedLabel && selectedLabel.type === 'all') {
        return <LonelyHost onClick={toggleAddHostModal} />;
      }

      return renderNoHosts();
    }

    return (
      <div className={`${baseClass} ${baseClass}--${displayType.toLowerCase()}`}>
        {renderHosts()}
        {renderPagination()}
      </div>
    );
  }
}

export default HostContainer;
