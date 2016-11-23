import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import QueriesList from './index';

const query1 = {
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

const query2 = {
  created_at: '2016-10-17T07:06:00Z',
  deleted: false,
  deleted_at: null,
  description: '',
  differential: false,
  id: 2,
  interval: 0,
  name: 'dev_query_2',
  platform: '',
  query: 'select * from time',
  snapshot: false,
  updated_at: '2016-10-17T07:06:00Z',
  version: '',
};

describe('QueriesList - component', () => {
  it('renders a QueriesListItem for each query', () => {
    const queries = [query1, query2];
    const onSelectQuery = () => {
      return noop;
    };

    const component = mount(
      <QueriesList
        onSelectQuery={onSelectQuery}
        queries={queries}
        selectedQueries={[]}
      />
    );

    expect(component.find('QueriesListItem').length).toEqual(2);
  });
});
