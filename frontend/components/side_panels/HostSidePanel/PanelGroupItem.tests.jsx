import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import PanelGroupItem from './PanelGroupItem';

describe('PanelGroupItem - component', () => {
  const validPanelGroupItem = {
    id: 1,
    label: 'All Hosts',
    name: 'all',
    count: 20,
  };

  const component = mount(
    <PanelGroupItem item={validPanelGroupItem} />
  );

  it('renders the icon', () => {
    const icon = component.find('i.kolidecon-hosts');

    expect(icon.length).toEqual(1);
  });

  it('renders the item text', () => {
    expect(component.text()).toContain(validPanelGroupItem.label);
  });

  it('renders the item count', () => {
    expect(component.text()).toContain(validPanelGroupItem.count);
  });
});
