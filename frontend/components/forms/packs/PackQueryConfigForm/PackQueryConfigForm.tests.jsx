import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import { fillInFormInput } from 'test/helpers';
import PackQueryConfigForm from './PackQueryConfigForm';

const formData = {
  queries: [],
};

describe('PackQueryConfigForm - component', () => {
  describe('interval field', () => {
    const form = mount(
      <PackQueryConfigForm formData={formData} handleSubmit={noop} />
    );

    it('renders an input field', () => {
      const intervalField = form.find('InputField').find({ name: 'interval' });

      expect(intervalField.length).toEqual(1);
    });

    it('updates state on field change', () => {
      const intervalField = form.find('InputField').find({ name: 'interval' });

      fillInFormInput(intervalField, '3600');

      expect(form.state().formData).toInclude({ interval: '3600' });
    });
  });
});
