import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import OrgDetails from 'components/forms/RegistrationForm/OrgDetails';
import { fillInFormInput } from 'test/helpers';

const noErrors = {};

describe('OrgDetails - form', () => {
  afterEach(restoreSpies);

  describe('organization name input', () => {
    it('renders an input field', () => {
      const form = mount(
        <OrgDetails
          errors={noErrors}
          formData={{}}
          onChange={noop}
          onSubmit={noop}
        />
      );
      const orgNameField = form.find({ name: 'organization name' });

      expect(orgNameField.length).toEqual(1);
    });

    it('calls the onChange prop when the field changes', () => {
      const onChangeSpy = createSpy();
      const form = mount(
        <OrgDetails
          errors={noErrors}
          formData={{}}
          onChange={onChangeSpy}
          onSubmit={noop}
        />
      );
      const orgNameField = form.find({ name: 'organization name' }).find('input');

      fillInFormInput(orgNameField, 'The Gnar Co');

      expect(onChangeSpy).toHaveBeenCalledWith('org_name', 'The Gnar Co');
    });
  });

  describe('organization logo URL input', () => {
    it('renders an input field', () => {
      const form = mount(
        <OrgDetails
          errors={noErrors}
          formData={{}}
          onChange={noop}
          onSubmit={noop}
        />
      );
      const orgLogoField = form.find({ name: 'org logo url' });

      expect(orgLogoField.length).toEqual(1);
    });

    it('calls the onChange prop when the field changes', () => {
      const onChangeSpy = createSpy();
      const form = mount(
        <OrgDetails
          errors={noErrors}
          formData={{}}
          onChange={onChangeSpy}
          onSubmit={noop}
        />
      );
      const orgLogoField = form.find({ name: 'org logo url' }).find('input');

      fillInFormInput(orgLogoField, 'http://www.thegnar.co/logo.png');

      expect(onChangeSpy).toHaveBeenCalledWith('org_logo_url', 'http://www.thegnar.co/logo.png');
    });
  });

  describe('submitting the form', () => {
    it('validates presence of all fields', () => {
      const onSubmitSpy = createSpy();
      const form = mount(
        <OrgDetails
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
        org_name: 'Organization name must be present',
        org_logo_url: 'Organization logo URL must be present',
      });
    });

    it('validates the logo url field starts with https://', () => {
      const onSubmitSpy = createSpy();
      const form = mount(
        <OrgDetails
          errors={noErrors}
          formData={{ org_logo_url: 'http://google.com' }}
          onChange={noop}
          onSubmit={onSubmitSpy}
        />
      );
      const submitBtn = form.find('Button');

      submitBtn.simulate('click');

      expect(onSubmitSpy).toNotHaveBeenCalled();
      expect(form.state().errors).toInclude({
        org_logo_url: 'Organization logo URL must start with https://',
      });
    });

    it('submits the form when valid', () => {
      const formData = {
        org_name: 'The Gnar Co.',
        org_logo_url: 'https://thegnar.co/assets/logo.png',
      };
      const onSubmitSpy = createSpy();
      const form = mount(
        <OrgDetails
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

