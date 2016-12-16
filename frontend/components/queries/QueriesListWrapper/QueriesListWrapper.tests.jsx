import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import { queryStub, scheduledQueryStub } from 'test/stubs';
import { fillInFormInput } from 'test/helpers';
import QueriesListWrapper from './index';

const allQueries = [queryStub];
const scheduledQueries = [
  scheduledQueryStub,
  { ...scheduledQueryStub, id: 100, name: 'mac hosts' },
];

describe('QueriesListWrapper - component', () => {
  it('renders the PackQueryConfigForm when "Add Query" is clicked', () => {
    const component = mount(
      <QueriesListWrapper
        allQueries={allQueries}
        scheduledQueries={scheduledQueries}
      />
    );

    const addQueryBtn = component.find('Button').first();

    addQueryBtn.simulate('click');
    expect(component.find('PackQueryConfigForm').length).toEqual(1);
  });

  it('filters queries', () => {
    const component = mount(
      <QueriesListWrapper
        allQueries={allQueries}
        scheduledQueries={scheduledQueries}
      />
    );

    const searchQueriesInput = component.find({ name: 'search-queries' });
    const QueriesList = component.find('QueriesList');

    expect(QueriesList.prop('scheduledQueries')).toEqual(scheduledQueries);

    fillInFormInput(searchQueriesInput, 'something that does not match');

    expect(QueriesList.prop('scheduledQueries')).toEqual([]);
  });
});
