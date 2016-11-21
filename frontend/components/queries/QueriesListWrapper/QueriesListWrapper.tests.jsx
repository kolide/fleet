import React from 'react';
import expect from 'expect';
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

describe('QueriesListWrapper - component', () => {
  it('renders a QueriesList component', () => {
    const queries = [query];
    const component = mount(<QueriesListWrapper queries={queries} />);

    expect(component.find('QueriesList').length).toEqual(1);
  });

  it('updates state when a query checkbox is changed', () => {
    const component = mount(<QueriesListWrapper queries={[query]} />);
    const checkbox = component.find('Checkbox').first();

    checkbox.simulate('change');

    expect(component.state().selectedQueries).toEqual([query]);

    checkbox.simulate('change');

    expect(component.state().selectedQueries).toEqual([]);
  });
});
