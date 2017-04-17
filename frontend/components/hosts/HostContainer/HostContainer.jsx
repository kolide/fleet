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
let CURRENT_PAGE = 0;

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
      hostsPerPage: 20,
      pagedHosts: [],
      showSpinner: false,
    };
  }

  componentWillMount () {
    this.buildSortedHosts();
  }

  componentWillReceiveProps (nextProps) {
    if (isEqual(nextProps, this.props)) {
      return false;
    }

    this.buildSortedHosts();
    return true;
  }

  shouldComponentUpdate (nextProps, nextState) {
    if (isEqual(nextProps, this.props) && isEqual(nextState, this.state)) {
      return false;
    }

    this.buildSortedHosts();
    return true;
  }

  buildSortedHosts = () => {
    const { filterHosts, sortHosts } = this;
    const { hostsPerPage } = this.state;

    const sortedHosts = sortHosts(filterHosts());

    const currentPage = CURRENT_PAGE - 1 < 0 ? 0 : CURRENT_PAGE - 1;
    const fromIndex = currentPage * hostsPerPage;
    const toIndex = fromIndex + hostsPerPage;

    const pagedHosts = slice(sortedHosts, fromIndex, toIndex);

    this.setState({
      allHostCount: sortedHosts.length,
      pagedHosts,
      showSpinner: false,
    });
  }

  filterHosts = () => {
    const { hosts, selectedLabel } = this.props;
    const { filterHosts } = helpers;

    return filterHosts(hosts, selectedLabel);
  }

  sortHosts = (hosts) => {
    const alphaHosts = sortBy(hosts, (h) => { return h.hostname; });
    const orderedHosts = orderBy(alphaHosts, 'status', 'desc');

    return orderedHosts;
  }

  handlePaginationChange = (page) => {
    CURRENT_PAGE = page;
    this.buildSortedHosts();

    return true;
  }

  handlePerPageChange = (option) => {
    this.setState({
      hostsPerPage: Number(option.value),
      showSpinner: true,
    });

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
    const { allHostCount, hostsPerPage } = this.state;

    const paginationSelectOpts = [
      { value: 20, label: '20 Hosts' },
      { value: 100, label: '100 Hosts' },
      { value: 500, label: '500 Hosts' },
      { value: 1000, label: '1,000 Hosts' },
    ];
    const currentPage = CURRENT_PAGE === 0 ? 1 : CURRENT_PAGE;
    const startRange = currentPage === 1 ? 1 : ((currentPage - 1) * hostsPerPage) + 1;
    const endRange = (currentPage * hostsPerPage) > allHostCount ? allHostCount : (currentPage * hostsPerPage);

    return (
      <div className={`${baseClass}__pager-wrap`}>
        <Pagination
          onChange={handlePaginationChange}
          current={currentPage}
          total={allHostCount}
          pageSize={hostsPerPage}
          className={`${baseClass}__pagination`}
          locale={enUs}
          showLessItems
        />
        <p className={`${baseClass}__pager-range`}>{`${startRange} - ${endRange} of ${allHostCount} items`}</p>
        <div className={`${baseClass}__pager-count`}>
          <Select
            name="pager-host-count"
            value={hostsPerPage}
            options={paginationSelectOpts}
            onChange={handlePerPageChange}
            className={`${baseClass}__count-select`}
            clearable={false}
          />
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
