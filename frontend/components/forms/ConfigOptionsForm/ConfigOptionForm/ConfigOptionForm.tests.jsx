import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import ConfigOptionForm from 'components/forms/ConfigOptionsForm/ConfigOptionForm';
import {
  itBehavesLikeAFormInputElement,
  itBehavesLikeAFormDropdownElement,
} from 'test/helpers';

describe('ConfigOptionForm - form', () => {
  afterEach(restoreSpies);

  it('renders form fields for the config option name and value', () => {
    const configNameOptions = [{ label: 'My option', value: 'my_option' }];
    const form = mount(
      <ConfigOptionForm configNameOptions={configNameOptions} handleSubmit={noop} />
    );

    itBehavesLikeAFormDropdownElement(form, 'name');
    itBehavesLikeAFormInputElement(form, 'value');
  });

  it('calls the onFormUpdate prop when the form updates', () => {
    const spy = createSpy();
    const configNameOptions = [{ label: 'My option', value: 'my_option' }];
    const form = mount(
      <ConfigOptionForm configNameOptions={configNameOptions} handleSubmit={noop} onFormUpdate={spy} />
    );

    itBehavesLikeAFormInputElement(form, 'value', 'InputField', 'new config option value');
    itBehavesLikeAFormDropdownElement(form, 'name');

    expect(spy).toHaveBeenCalledWith({
      errors: { base: null, name: null, value: null },
      formData: {
        name: 'my_option',
        value: 'new config option value',
      },
    });
  });
});

