import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import AppConfigForm from 'components/forms/admin/AppConfigForm';
import { itBehavesLikeAFormInputElement } from 'test/helpers';

describe.only('AppConfigForm - form', () => {
  const form = mount(<AppConfigForm handleSubmit={noop} />);

  describe('Organization Name input', () => {
    it('renders an input field', () => {
      itBehavesLikeAFormInputElement(form, 'org_name');
    });
  });

  describe('Organization Logo URL input', () => {
    it('renders an input field', () => {
      itBehavesLikeAFormInputElement(form, 'org_logo_url');
    });
  });
});
