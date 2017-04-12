import React from 'react';
import expect from 'expect';
import { noop } from 'lodash';
import { mount } from 'enzyme';

import HostContainer from './HostContainer';
import { hostStub } from 'test/stubs';

const allHostsLabel = { id: 1, display_text: 'All Hosts', slug: 'all-hosts', type: 'all', count: 22 };
const customLabel = { id: 6, display_text: 'Custom Label', slug: 'custom-label', type: 'custom', count: 3 };

describe.only('HostsContainer - component', () => {
  const props = {
    hosts: [hostStub],
    selectedLabel: allHostsLabel,
    loadingHosts: false,
    displayType: 'Grid',
    toggleAddHostModal: noop,
    toggleDeleteHostModal: noop,
    onQueryHost: noop,
  };

  it('renders Spinner while hosts are loading', () => {
    const loadingProps = { ...props, loadingHosts: true };
    const page = mount(<HostContainer {...loadingProps} hosts={[]} selectedLabel={allHostsLabel} />);

    expect(page.find('Spinner').length).toEqual(1);
  });

  it('render LonelyHost if no hosts available', () => {
    const page = mount(<HostContainer {...props} hosts={[]} selectedLabel={allHostsLabel} />);

    expect(page.find('LonelyHost').length).toEqual(1);
  });

  it('renders message if no hosts available and not on All Hosts', () => {
    const page = mount(<HostContainer {...props} hosts={[]} selectedLabel={customLabel} />);

    expect(page.find('.host-container__no-hosts').length).toEqual(1);
  });

  it('renders hosts as HostDetails by default', () => {
    const page = mount(<HostContainer {...props} />);

    expect(page.find('HostDetails').length).toEqual(1);
  });
});






  //   it('does not render sidebar if labels are loading', () => {
  //     const loadingProps = { ...props, loadingLabels: true };
  //     const page = mount(<ManageHostsPage {...loadingProps} hosts={[]} selectedLabel={allHostsLabel} />);

  //     expect(page.find('HostSidePanel').length).toEqual(0);
  //   });







  //   it('renders hosts as HostsTable when the display is "List"', () => {
  //     const page = mount(<ManageHostsPage {...props} display="List" hosts={[hostStub]} />);

  //     expect(page.find('HostsTable').length).toEqual(1);
  //   });

  //   it('toggles between displays', () => {
  //     const ownProps = { location: {}, params: {} };
  //     const component = connectedComponent(ConnectedManageHostsPage, { props: ownProps, mockStore });
  //     const page = mount(component);
  //     const button = page.find('Rocker').find('button');
  //     const toggleDisplayAction = {
  //       type: 'SET_DISPLAY',
  //       payload: {
  //         display: 'List',
  //       },
  //     };

  //     button.simulate('click');

  //     expect(mockStore.getActions()).toInclude(toggleDisplayAction);
  //   });

  //   it('filters hosts', () => {
  //     const allHostsLabelPageNode = mount(
  //       <ManageHostsPage
  //         {...props}
  //         hosts={[hostStub, offlineHost]}
  //         selectedLabel={allHostsLabel}
  //       />
  //     ).node;
  //     const offlineHostsLabelPageNode = mount(
  //       <ManageHostsPage
  //         {...props}
  //         hosts={[hostStub, offlineHost]}
  //         selectedLabel={offlineHostsLabel}
  //       />
  //     ).node;

  //     expect(allHostsLabelPageNode.filterHosts()).toEqual([hostStub, offlineHost]);
  //     expect(offlineHostsLabelPageNode.filterHosts()).toEqual([offlineHost]);
  //   });
  // });
