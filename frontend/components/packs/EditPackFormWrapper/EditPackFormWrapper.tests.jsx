import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import EditPackFormWrapper from 'components/packs/EditPackFormWrapper';
import { packStub } from 'test/stubs';

describe('EditPackFormWrapper - component', () => {
  it('does not render the EditPackForm by default', () => {
    const component = mount(<EditPackFormWrapper pack={packStub} />);

    expect(component.find('EditPackForm').length).toEqual(0);
  });

  it('renders the EditPackForm when isEditing is true', () => {
    const component = mount(<EditPackFormWrapper pack={packStub} />);

    component.setState({ isEditing: true });

    expect(component.find('EditPackForm').length).toEqual(1);
  });

  it('sets state to isEditing when the EDIT button is clicked', () => {
    const component = mount(<EditPackFormWrapper pack={packStub} />);
    const editBtn = component.find('Button').findWhere(b => b.prop('text') === 'EDIT');

    editBtn.simulate('click');

    expect(component.state('isEditing')).toEqual(true);
  });
});
