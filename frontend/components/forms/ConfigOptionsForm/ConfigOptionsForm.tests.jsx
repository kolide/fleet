import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import ConfigOptionsForm from 'components/forms/ConfigOptionsForm';
import { configOptionStub } from 'test/stubs';

describe('ConfigOptionsForm - form', () => {
  it('renders a ConfigOptionForm for each completed config option', () => {
    const formWithOneOption = mount(<ConfigOptionsForm completedOptions={[configOptionStub]} />);
    const formWithTwoOptions = mount(<ConfigOptionsForm completedOptions={[configOptionStub, configOptionStub]} />);

    expect(formWithOneOption.find('ConfigOptionForm').length).toEqual(1);
    expect(formWithTwoOptions.find('ConfigOptionForm').length).toEqual(2);
  });
});
