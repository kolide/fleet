import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';

import Checkbox from './Checkbox';

describe('Checkbox - component', () => {
  afterEach(restoreSpies);

  it('renders', () => {
    expect(mount(<Checkbox />)).toExist();
  });

  it('calls the "onChange" handler when changed', () => {
    const onCheckedComponentChangeSpy = createSpy();
    const onUncheckedComponentChangeSpy = createSpy();

    const checkedComponent = mount(
      <Checkbox
        checked
        name="checkbox"
        onChange={onCheckedComponentChangeSpy}
      />
    ).find('input');

    const uncheckedComponent = mount(
      <Checkbox
        checked={false}
        name="checkbox"
        onChange={onUncheckedComponentChangeSpy}
      />
    ).find('input');

    checkedComponent.simulate('change');
    uncheckedComponent.simulate('change');

    expect(onCheckedComponentChangeSpy).toHaveBeenCalledWith(false);
    expect(onUncheckedComponentChangeSpy).toHaveBeenCalledWith(true);
  });
});
