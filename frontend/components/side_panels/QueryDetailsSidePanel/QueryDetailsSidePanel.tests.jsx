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

  it('renders a read-only Kolide Ace component with the query text', () => {
    const component = mount(<QueryDetailsSidePanel query={queryStub} />);
    const aceEditor = component.find('KolideAce');

    expect(aceEditor.length).toEqual(1);
    expect(aceEditor.prop('value')).toEqual(queryStub.query);
    expect(aceEditor.prop('readOnly')).toEqual(true);
  });
});
