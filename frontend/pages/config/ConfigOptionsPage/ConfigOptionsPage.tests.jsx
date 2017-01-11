import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import { ConfigOptionsPage } from 'pages/config/ConfigOptionsPage/ConfigOptionsPage';
import { configOptionStub } from 'test/stubs';

describe('ConfigOptionsPage - component', () => {
  const props = { configOptions: [] };
  const page = mount(<ConfigOptionsPage {...props} />);

  it('renders', () => {
    expect(page.length).toEqual(1);
  });

  it('renders reset and save buttons', () => {
    const buttons = page.find('Button');
    const resetButton = buttons.find('.config-options-page__reset-btn');
    const saveButton = buttons.find('.config-options-page__save-btn');

    expect(resetButton.length).toEqual(1);
    expect(saveButton.length).toEqual(1);
  });

  describe('removing a config option', () => {
    it('sets the option value to null in state', () => {
      const page = mount(<ConfigOptionsPage configOptions={[configOptionStub]} />);
      const removeBtn = page.find('ConfigOptionForm').find('Button');

      expect(page.state('configOptions')).toEqual([configOptionStub]);

      removeBtn.simulate('click');

      expect(page.state('configOptions')).toEqual([{
        ...configOptionStub,
        value: null,
      }]);
    });
  });
});
