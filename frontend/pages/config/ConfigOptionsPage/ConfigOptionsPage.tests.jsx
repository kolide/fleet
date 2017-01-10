import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import { ConfigOptionsPage } from 'pages/config/ConfigOptionsPage/ConfigOptionsPage';

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
});
