import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import ConfigOptionsPage from 'pages/config/ConfigOptionsPage';
import { configOptionStub } from 'test/stubs';

describe.only('ConfigOptionsPage - component', () => {
  const props = {};
  const page = mount(<ConfigOptionsPage props={props} />);

  it('renders', () => {
    expect(page.length).toEqual(1);
    expect(page.state('configOptions')).toEqual([]);
  });

  it('renders reset and save buttons', () => {
    const buttons = page.find('Button');
    const resetButton = buttons.find('.config-options-page__reset-btn');
    const saveButton = buttons.find('.config-options-page__save-btn');

    expect(resetButton.length).toEqual(1);
    expect(saveButton.length).toEqual(1);
  });

  it('updates the states configOptions when props change', () => {
    expect(page.state('configOptions')).toEqual([]);

    page.setProps({ config_options: [configOptionStub] });
    expect(page.state('configOptions')).toEqual([configOptionStub]);
  });
});
