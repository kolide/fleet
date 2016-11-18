import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import { PackComposerPage } from './PackComposerPage';

describe('PackComposerPage - component', () => {
  it('renders', () => {
    const page = mount(<PackComposerPage />);

    expect(page.length).toEqual(1);
  });

  it('renders a PackForm component', () => {
    const page = mount(<PackComposerPage />);

    expect(page.find('PackForm').length).toEqual(1);
  });
});
