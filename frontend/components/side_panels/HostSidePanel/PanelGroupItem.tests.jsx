import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import PanelGroupItem from './PanelGroupItem';

describe('PanelGroupItem - component', () => {
  const validPanelGroupItem = {
    hosts_count: 20,
    type: 'all',
    title: 'All Hosts',
  };

  const component = mount(
    <PanelGroupItem item={validPanelGroupItem} />
  );

  it('renders the icon', () => {
    const icon = component.find('i.kolidecon-hosts');

    expect(icon.length).toEqual(1);
  });

  it('renders the item text', () => {
    expect(component.text()).toContain(validPanelGroupItem.title);
  });

  it('renders the item count', () => {
    expect(component.text()).toContain(validPanelGroupItem.hosts_count);
  });
});
