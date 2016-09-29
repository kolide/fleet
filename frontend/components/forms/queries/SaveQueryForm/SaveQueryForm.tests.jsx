import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';
import helpers from '../../../../test/helpers';
import SaveQueryForm from './index';

const { fillInFormInput } = helpers;
const queryName = 'My New Query';

describe('SaveQueryForm - component', () => {
  afterEach(restoreSpies);

  it('handles query name input changes', () => {
    const form = mount(
      <SaveQueryForm onSubmit={noop} />
    );
    const queryNameInput = form.find({ name: 'queryName' });

    fillInFormInput(queryNameInput, queryName);

    const { formData } = form.state();

    expect(formData).toEqual({
      queryDescription: null,
      queryDuration: 'short',
      queryHosts: 'all',
      queryHostsPercentage: null,
      queryName,
      queryPlatform: 'all',
      scanInterval: 0,
    });
  });

  it('does not submit the form if it is invalid', () => {
    const onSubmit = createSpy();
    const form = mount(
      <SaveQueryForm onSubmit={onSubmit} />
    );

    form.simulate('submit');

    expect(onSubmit).toNotHaveBeenCalled();
  });

  it('calls onSubmit with the formData', () => {
    const onSubmit = createSpy();
    const form = mount(
      <SaveQueryForm onSubmit={onSubmit} />
    );
    const queryNameInput = form.find({ name: 'queryName' });

    fillInFormInput(queryNameInput, queryName);
    form.simulate('submit');

    expect(onSubmit).toHaveBeenCalledWith({
      runType: 'run',
      formData: {
        queryDescription: null,
        queryDuration: 'short',
        queryHosts: 'all',
        queryHostsPercentage: null,
        queryName,
        queryPlatform: 'all',
        scanInterval: 0,
      },
    });
  });
});
