import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import DropdownButton from './DropdownButton';

describe('DropdownButton - component', () => {
  afterEach(restoreSpies);
  const optionSpy = createSpy();
  const dropdownOptions = [{ label: 'btn1', onClick: noop }, { label: 'btn2', onClick: optionSpy }];
  const component = mount(
    <DropdownButton options={dropdownOptions}>
      New Button
    </DropdownButton>
  );

  it('Changes state on click', () => {
    expect(component.state().isOpen).toEqual(false);

    component.find('.dropdown-button').simulate('click');

    expect(component.state().isOpen).toEqual(true);
  });

  it("calls the clicked item's onClick attribute", () => {
    component.find('.dropdown-button').simulate('click');
    component.find('li').last().find('Button').simulate('click');

    expect(optionSpy).toHaveBeenCalled();
  });
});
