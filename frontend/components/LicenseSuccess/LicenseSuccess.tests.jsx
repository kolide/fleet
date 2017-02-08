import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import { licenseStub } from 'test/stubs';
import LicenseSuccess from 'components/LicenseSuccess';

const defaultProps = {
  license: licenseStub,
  onConfirmLicense: noop,
};

describe('LicenseSuccess - component', () => {
  afterEach(restoreSpies);

  describe('rendering', () => {
    it('renders', () => {
      expect(mount(<LicenseSuccess {...defaultProps} />).length).toEqual(1, 'Expected LicenseSuccess component to render');
    });
  });

  it('calls the onConfirmLicense prop when the button is clicked', () => {
    const spy = createSpy();
    const props = { ...defaultProps, onConfirmLicense: spy };
    const Component = mount(<LicenseSuccess {...props} />);

    Component
      .find('Button')
      .simulate('click');

    expect(spy).toHaveBeenCalled();
  });
});
