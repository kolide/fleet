import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import { connectedComponent, reduxMockStore } from 'test/helpers';
import ConnectedPacksComposerPage, { PackComposerPage } from './PackComposerPage';

describe('PackComposerPage - component', () => {
  it('renders', () => {
    const page = mount(connectedComponent(ConnectedPacksComposerPage));

    expect(page.length).toEqual(1);
  });

  it('renders a PackForm component', () => {
    const page = mount(connectedComponent(ConnectedPacksComposerPage));

    expect(page.find('PackForm').length).toEqual(1);
  });

  it('renders a PackInfoSidePanel component', () => {
    const page = mount(connectedComponent(ConnectedPacksComposerPage));

    expect(page.find('PackInfoSidePanel').length).toEqual(1);
  });
});
