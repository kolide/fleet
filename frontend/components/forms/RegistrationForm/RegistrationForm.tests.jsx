import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import RegistrationForm from 'components/forms/RegistrationForm';

describe('RegistrationForm - component', () => {
  it('renders AdminDetails on the first page', () => {
    const form = mount(<RegistrationForm page={1} />);

    expect(form.find('AdminDetails').length).toEqual(1);
  });

  it('renders OrgDetails on the second page', () => {
    const form = mount(<RegistrationForm page={2} />);

    expect(form.find('OrgDetails').length).toEqual(1);
  });

  it('renders KolideDetails on the third page', () => {
    const form = mount(<RegistrationForm page={3} />);

    expect(form.find('KolideDetails').length).toEqual(1);
  });
});

