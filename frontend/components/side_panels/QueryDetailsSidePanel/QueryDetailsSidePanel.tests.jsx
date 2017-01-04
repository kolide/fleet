import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import QueryDetailsSidePanel from 'components/side_panels/QueryDetailsSidePanel';
import { queryStub } from 'test/stubs';

describe('QueryDetailsSidePanel - component', () => {
  it('renders', () => {
    const component = mount(<QueryDetailsSidePanel query={queryStub} />);

    expect(component.length).toEqual(1);
  });
});
