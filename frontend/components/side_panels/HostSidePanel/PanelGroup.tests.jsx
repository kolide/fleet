import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import PanelGroup from './PanelGroup';

describe('PanelGroup - component', () => {
  const validPanelGroupItems = [
    { id: 1, label: 'All Hosts', name: 'all', count: 20 },
    { id: 2, label: 'MAC OS', name: 'macs', count: 10 },
    { id: 3, label: 'ONLINE', name: 'online', count: 10 },
  ];

  const component = mount(
    <PanelGroup groupItems={validPanelGroupItems} />
  );

  it('renders a PanelGroupItem for each group item', () => {
    const panelGroupItems = component.find('PanelGroupItem');

    expect(panelGroupItems.length).toEqual(3);
  });
});

