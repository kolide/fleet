import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import ConnectedAllPacksPage, { AllPacksPage } from 'pages/packs/AllPacksPage/AllPacksPage';
import { connectedComponent, fillInFormInput, reduxMockStore } from 'test/helpers';
import { packStub } from 'test/stubs';

const store = {
  entities: {
    packs: {
      data: {
        [packStub.id]: packStub,
        101: {
          ...packStub,
          id: 101,
          name: 'My unique pack name',
        },
      },
    },
  },
};

describe('AllPacksPage - component', () => {
  it('filters the packs list', () => {
    const Component = connectedComponent(ConnectedAllPacksPage, {
      mockStore: reduxMockStore(store),
    });
    const page = mount(Component).find('AllPacksPage');
    const packsFilterInput = page.find({ name: 'pack-filter' }).find('input');

    expect(page.node.getPacks().length).toEqual(2);

    fillInFormInput(packsFilterInput, 'My unique pack name');

    expect(page.node.getPacks().length).toEqual(1);
  });

  it('renders a PacksList component', () => {
    const page = mount(connectedComponent(ConnectedAllPacksPage));

    expect(page.find('PacksList').length).toEqual(1);
  });

  it('renders the PackInfoSidePanel by default', () => {
    const page = mount(connectedComponent(ConnectedAllPacksPage));

    expect(page.find('PackInfoSidePanel').length).toEqual(1);
  });

  it('updates checkedPackIDs in state when the select all packs Checkbox is toggled', () => {
    const page = mount(<AllPacksPage packs={[packStub]} />);
    const selectAllPacks = page.find({ name: 'select-all-packs' });

    expect(page.state('checkedPackIDs')).toEqual([]);

    selectAllPacks.simulate('change');

    expect(page.state('checkedPackIDs')).toEqual([packStub.id]);

    selectAllPacks.simulate('change');

    expect(page.state('checkedPackIDs')).toEqual([]);
  });

  it('updates checkedPackIDs in state when a pack row Checkbox is toggled', () => {
    const page = mount(<AllPacksPage packs={[packStub]} />);
    const selectPack = page.find({ name: `select-pack-${packStub.id}` });

    expect(page.state('checkedPackIDs')).toEqual([]);

    selectPack.simulate('change');

    expect(page.state('checkedPackIDs')).toEqual([packStub.id]);

    selectPack.simulate('change');

    expect(page.state('checkedPackIDs')).toEqual([]);
  });
});
