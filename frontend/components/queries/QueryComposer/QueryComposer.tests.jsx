import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import { createAceSpy, fillInFormInput } from 'test/helpers';
import QueryComposer from './index';

describe('QueryComposer - component', () => {
  beforeEach(() => {
    createAceSpy();
  });
  afterEach(restoreSpies);

  it('does not render the SaveQueryForm by default', () => {
    const component = mount(
      <QueryComposer
        onOsqueryTableSelect={noop}
        onTextEditorInputChange={noop}
        selectedTargets={[]}
        textEditorText="Hello world"
      />
    );

    expect(component.find('SaveQueryForm').length).toEqual(0);
  });

  it('renders the SaveQueryFormModal when "Save Query" is clicked', () => {
    const component = mount(
      <QueryComposer
        onOsqueryTableSelect={noop}
        onTextEditorInputChange={noop}
        selectedTargets={[]}
        textEditorText="Hello world"
      />
    );

    component.find('.query-composer__save-query-btn').simulate('click');

    expect(component.find('SaveQueryForm').length).toEqual(1);
  });

  it('renders the UpdateQueryForm when the query prop is present', () => {
    const query = {
      id: 1,
      query: 'SELECT * FROM users',
      name: 'Get all users',
      description: 'This gets all of the users',
    };
    const component = mount(
      <QueryComposer
        onOsqueryTableSelect={noop}
        onTextEditorInputChange={noop}
        query={query}
        selectedTargets={[]}
        textEditorText="Hello world"
      />
    );

    const form = component.find('UpdateQueryForm');

    expect(form.length).toEqual(1);
    expect(form.find('InputField').length).toEqual(2);
  });

  it('renders the Run Query button as disabled without selected targets', () => {
    const component = mount(
      <QueryComposer
        onOsqueryTableSelect={noop}
        onTextEditorInputChange={noop}
        selectedTargets={[]}
        textEditorText="Hello world"
      />
    );

    const runQueryBtn = component.find('.query-composer__run-query-btn');

    expect(runQueryBtn.props()).toInclude({
      disabled: true,
    });
  });

  it('hides the SaveQueryFormModal after the form is submitted', () => {
    const component = mount(
      <QueryComposer
        onSaveQueryFormSubmit={noop}
        selectedTargets={[]}
        textEditorText="SELECT * FROM users"
      />
    );

    component.find('.query-composer__save-query-btn').simulate('click');

    const form = component.find('SaveQueryForm');

    fillInFormInput(form.find({ name: 'name' }), 'My query name');
    form.simulate('submit');

    expect(component.find('SaveQueryForm').length).toEqual(0);
  });

  it('calls onSaveQueryFormSubmit with appropriate data from SaveQueryFormModal', () => {
    const onSaveQueryFormSubmitSpy = createSpy();
    const query = 'SELECT * FROM users';
    const selectedTargets = [{ name: 'my target' }];
    const component = mount(
      <QueryComposer
        onSaveQueryFormSubmit={onSaveQueryFormSubmitSpy}
        selectedTargets={selectedTargets}
        textEditorText={query}
      />
    );

    component.find('.query-composer__save-query-btn').simulate('click');

    const form = component.find('SaveQueryForm');

    fillInFormInput(form.find({ name: 'name' }), 'My query name');
    fillInFormInput(form.find({ name: 'description' }), 'My query description');
    form.simulate('submit');

    expect(onSaveQueryFormSubmitSpy).toHaveBeenCalledWith({
      description: 'My query description',
      name: 'My query name',
    });
  });

  it('calls onRunQuery when "Run Query" is clicked', () => {
    const onRunQuerySpy = createSpy();
    const query = 'SELECT * FROM users';
    const selectedTargets = [{ name: 'my target' }];
    const component = mount(
      <QueryComposer
        onRunQuery={onRunQuerySpy}
        selectedTargets={selectedTargets}
        textEditorText={query}
      />
    );
    component.find('.query-composer__run-query-btn').simulate('click');

    expect(onRunQuerySpy).toHaveBeenCalled();
  });

  it('calls onTargetSelectInputChange when changing the select target input text', () => {
    const onTargetSelectInputChangeSpy = createSpy();
    const component = mount(
      <QueryComposer
        onTargetSelectInputChange={onTargetSelectInputChangeSpy}
        selectedTargets={[]}
      />
    );
    const selectTargetsInput = component.find('.Select-input input');

    fillInFormInput(selectTargetsInput, 'my target');

    expect(onTargetSelectInputChangeSpy).toHaveBeenCalledWith('my target');
  });
});
