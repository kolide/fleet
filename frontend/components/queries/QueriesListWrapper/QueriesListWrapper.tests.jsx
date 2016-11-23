import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';

import QueriesListWrapper from './index';

const query = {
  created_at: '2016-10-17T07:06:00Z',
  deleted: false,
  deleted_at: null,
  description: '',
  differential: false,
  id: 1,
  interval: 0,
  name: 'dev_query_1',
  platform: '',
  query: 'select * from processes',
  snapshot: false,
  updated_at: '2016-10-17T07:06:00Z',
  version: '',
};
const queries = [query];

describe('QueriesListWrapper - component', () => {
  afterEach(restoreSpies);

  it('renders the PackQueryConfigForm when there are staged queries', () => {
    const componentWithoutStagedQueries = mount(
      <QueriesListWrapper
        configuredQueries={[]}
        queries={queries}
        stagedQueries={[]}
      />
    );
    const componentWithStagedQueries = mount(
      <QueriesListWrapper
        configuredQueries={[]}
        queries={queries}
        stagedQueries={queries}
      />
    );

    expect(
      componentWithoutStagedQueries.find('PackQueryConfigForm').length
    ).toEqual(0);
    expect(
      componentWithStagedQueries.find('PackQueryConfigForm').length
    ).toEqual(1);
  });

  it('calls the onSelectQuery prop when a query checkbox is selected', () => {
    const onSelectQuerySpy = createSpy();
    const component = mount(
      <QueriesListWrapper
        configuredQueries={[]}
        onSelectQuery={onSelectQuerySpy}
        queries={queries}
        stagedQueries={[]}
      />
    );
    const checkbox = component.find('Checkbox').first();

    checkbox.simulate('change');

    expect(onSelectQuerySpy).toHaveBeenCalledWith(query);
  });

  it('calls the onSelectQuery prop when a query checkbox is changed', () => {
    const onDeselectQuerySpy = createSpy();
    const component = mount(
      <QueriesListWrapper
        configuredQueries={[]}
        onDeselectQuery={onDeselectQuerySpy}
        queries={queries}
        stagedQueries={queries}
      />
    );
    const checkbox = component.find('Checkbox').first();

    checkbox.simulate('change');

    expect(onDeselectQuerySpy).toHaveBeenCalledWith(query);
  });
});
