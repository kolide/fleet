import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import KolideDetails from 'components/forms/RegistrationForm/KolideDetails';
import { fillInFormInput } from 'test/helpers';

const noErrors = {};

describe('KolideDetails - form', () => {
  afterEach(restoreSpies);

  describe('kolide web address input', () => {
    it('renders an input field', () => {
      const form = mount(
        <KolideDetails
          errors={noErrors}
          formData={{}}
          onChange={noop}
          onSubmit={noop}
        />
      );
      const kolideWebAddressField = form.find({ name: 'kolide web address' });

      expect(kolideWebAddressField.length).toEqual(1);
    });

    it('calls the onChange prop when the field changes', () => {
      const onChangeSpy = createSpy();
      const form = mount(
        <KolideDetails
          errors={noErrors}
          formData={{}}
          onChange={onChangeSpy}
          onSubmit={noop}
        />
      );
      const kolideWebAddressField = form.find({ name: 'kolide web address' }).find('input');

      fillInFormInput(kolideWebAddressField, 'https://gnar.kolide.co');

      expect(onChangeSpy).toHaveBeenCalledWith('kolide_web_address', 'https://gnar.kolide.co');
    });
  });

  describe('submitting the form', () => {
    it('validates the presence of the kolide web address field', () => {
      const onSubmitSpy = createSpy();
      const form = mount(
        <KolideDetails
          errors={noErrors}
          formData={{}}
          onChange={noop}
          onSubmit={onSubmitSpy}
        />
      );
      const submitBtn = form.find('Button');

      submitBtn.simulate('click');

      expect(onSubmitSpy).toNotHaveBeenCalled();
      expect(form.state().errors).toInclude({
        kolide_web_address: 'Kolide web address must be completed',
      });
    });

    it('validates the kolide web address field starts with https://', () => {
      const onSubmitSpy = createSpy();
      const form = mount(
        <KolideDetails
          errors={noErrors}
          formData={{ kolide_web_address: 'http://google.com' }}
          onChange={noop}
          onSubmit={onSubmitSpy}
        />
      );
      const submitBtn = form.find('Button');

      submitBtn.simulate('click');

      expect(onSubmitSpy).toNotHaveBeenCalled();
      expect(form.state().errors).toInclude({
        kolide_web_address: 'Kolide web address must start with https://',
      });
    });

    it('submits the form when valid', () => {
      const formData = {
        kolide_web_address: 'https://gnar.kolide.co',
      };
      const onSubmitSpy = createSpy();
      const form = mount(
        <KolideDetails
          errors={noErrors}
          formData={formData}
          onChange={noop}
          onSubmit={onSubmitSpy}
        />
      );
      const submitBtn = form.find('Button');

      submitBtn.simulate('click');

      expect(onSubmitSpy).toHaveBeenCalled();
    });
  });
});

