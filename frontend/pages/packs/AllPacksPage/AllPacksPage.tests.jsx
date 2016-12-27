import expect from 'expect';
import { mount } from 'enzyme';

import AllPacksPage from 'pages/packs/AllPacksPage';
import { connectedComponent, fillInFormInput, reduxMockStore } from 'test/helpers';
import { packStub } from 'test/stubs';

describe('AllPacksPage - component', () => {
  it('filters the packs list', () => {
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
    const Component = connectedComponent(AllPacksPage, {
      mockStore: reduxMockStore(store),
    });
    const page = mount(Component).find('AllPacksPage');
    const packsFilterInput = page.find({ name: 'pack-filter' }).find('input');

    expect(page.node.getPacks().length).toEqual(2);

    fillInFormInput(packsFilterInput, 'My unique pack name');

    expect(page.node.getPacks().length).toEqual(1);
  });

  it('renders a PacksList component', () => {
    const page = mount(connectedComponent(AllPacksPage));

    expect(page.find('PacksList').length).toEqual(1);
  });

  it('renders the PackInfoSidePanel by default', () => {
    const page = mount(connectedComponent(AllPacksPage));

    expect(page.find('PackInfoSidePanel').length).toEqual(1);
  });
});
