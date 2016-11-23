import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import QueriesListItem from './index';

const query = {
  created_at: '2016-10-17T07:06:00Z',
  deleted: false,
  deleted_at: null,
  description: 'This is my query',
  differential: false,
  id: 1,
  interval: 0,
  name: 'dev_query_1',
  platform: 'darwin',
  query: 'select * from processes',
  snapshot: false,
  updated_at: '2016-10-17T07:06:00Z',
  version: '',
};

describe('QueriesListItem - component', () => {
  const component = mount(<QueriesListItem configured={false} checked={false} onSelect={noop} query={query} />);

  it('renders the query data', () => {
    expect(component.text()).toInclude(query.name);
    expect(component.text()).toInclude(query.description);
    expect(component.find('.kolidecon-apple').length).toEqual(1);
  });

  it('renders a Checkbox component when the query has not been configured', () => {
    expect(component.find('Checkbox').length).toEqual(1);
  });

  it('renders a check mark icon when the query has been configured', () => {
    const configuredComponent = mount(<QueriesListItem configured checked={false} onSelect={noop} query={query} />);

    expect(configuredComponent.find('Checkbox').length).toEqual(0);
    expect(
      configuredComponent.find('i.kolidecon-success-check').length
    ).toEqual(1);
  });
});
