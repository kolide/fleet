import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import { connectedComponent, reduxMockStore } from 'test/helpers';
import ConnectedPacksComposerPage, { PackComposerPage } from './PackComposerPage';

describe('PackComposerPage - component', () => {
  it('renders', () => {
    const page = mount(<PackComposerPage />);

    expect(page.length).toEqual(1);
  });

  it('renders a PackForm component', () => {
    const page = mount(<PackComposerPage />);

    expect(page.find('PackForm').length).toEqual(1);
  });

  it('renders a QueriesListWrapper component', () => {
    const page = mount(<PackComposerPage />);

    expect(page.find('QueriesListWrapper').length).toEqual(1);
  });

  it('loads all queries when it mounts', () => {
    const mockStore = reduxMockStore({
      components: {
        PacksPages: {},
      },
    });
    const page = mount(
      connectedComponent(ConnectedPacksComposerPage, { mockStore })
    );

    expect(page.length).toEqual(1);
    expect(mockStore.getActions()).toInclude({ type: 'queries_LOAD_REQUEST' });
  });
});
