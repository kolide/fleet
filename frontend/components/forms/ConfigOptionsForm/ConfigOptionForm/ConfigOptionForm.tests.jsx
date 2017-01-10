import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import ConfigOptionForm from 'components/forms/ConfigOptionsForm/ConfigOptionForm';
import { itBehavesLikeAFormInputElement } from 'test/helpers';

describe('ConfigOptionForm - form', () => {
  afterEach(restoreSpies);

  it('renders form fields for the config option name and value', () => {
    const form = mount(<ConfigOptionForm handleSubmit={noop} />);

    itBehavesLikeAFormInputElement(form, 'name');
    itBehavesLikeAFormInputElement(form, 'value');
  });

  it('calls the onFormUpdate prop when the form updates', () => {
    const spy = createSpy();
    const form = mount(<ConfigOptionForm handleSubmit={noop} onFormUpdate={spy} />);

    itBehavesLikeAFormInputElement(form, 'name', 'InputField', 'new config option name');

    expect(spy).toHaveBeenCalledWith({
      errors: { base: null, name: null },
      formData: { name: 'new config option name' },
    });
  });
});

