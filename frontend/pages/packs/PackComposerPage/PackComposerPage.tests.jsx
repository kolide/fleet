import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import { PackComposerPage } from './PackComposerPage';

describe('PackComposerPage - component', () => {
  it('renders', () => {
    const page = mount(<PackComposerPage />);

    expect(page.length).toEqual(1);
  });
});
