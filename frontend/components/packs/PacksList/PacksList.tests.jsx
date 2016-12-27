import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import PacksList from 'components/packs/PacksList';
import { packStub } from 'test/stubs';

describe('PacksList - component', () => {
  it('renders', () => {
    expect(mount(<PacksList packs={[packStub]} />).length).toEqual(1);
  });

  it('selects all packs when the Checkbox is checked', () => {
    const component = mount(<PacksList packs={[packStub]} />);

    expect(component.state('allPacksChecked')).toEqual(false);

    component.find({ name: 'select-all-packs' }).simulate('change');

    expect(component.state('allPacksChecked')).toEqual(true);
  });

  it('updates the checked pack IDs when an individual pack is checked', () => {
    const component = mount(<PacksList packs={[packStub]} />);
    const packCheckbox = component.find({ name: `select-pack-${packStub.id}` });

    expect(component.state('checkedPackIDs')).toEqual([]);

    packCheckbox.simulate('change');

    expect(component.state('checkedPackIDs')).toEqual([packStub.id]);

    packCheckbox.simulate('change');

    expect(component.state('checkedPackIDs')).toEqual([]);
  });
});
