import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import PackForm from './index';

describe('PackForm - component', () => {
  it('renders the correct components', () => {
    const page = mount(<PackForm />);

    expect(page.find('InputField').length).toEqual(2);
    expect(page.find('Button').length).toEqual(1);
  });
});
