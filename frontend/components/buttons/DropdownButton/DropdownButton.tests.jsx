import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import DropdownButton from './DropdownButton';

const fauxClick = () => {};

describe('DropdownButton - component', () => {
  it('Changes state on click', () => {
    const component = mount(
      <DropdownButton
        options={[{
          label: 'Button 1',
          onClick: fauxClick,
        }, {
          label: 'Button 2',
          onClick: fauxClick,
        }]}
      >
        New Button
      </DropdownButton>
    );

    expect(component.state().isOpen).toEqual(false);

    component.find('.dropdown-button').simulate('click');

    expect(component.state().isOpen).toEqual(true);
  });
});
