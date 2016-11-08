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

      expect(onChangeSpy).toHaveBeenCalledWith({
        org_name: 'The Gnar Co',
      });
    });
  });

  describe('organization web URL input', () => {
    it('renders an input field', () => {
      const form = mount(
        <OrgDetails
          errors={noErrors}
          formData={{}}
          onChange={noop}
          onSubmit={noop}
        />
      );
      const orgURLField = form.find({ name: 'org web url' });

      expect(orgURLField.length).toEqual(1);
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
      const orgURLField = form.find({ name: 'org web url' }).find('input');

      fillInFormInput(orgURLField, 'http://www.thegnar.co');

      expect(onChangeSpy).toHaveBeenCalledWith({
        org_web_url: 'http://www.thegnar.co',
      });
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

      expect(onChangeSpy).toHaveBeenCalledWith({
        org_logo_url: 'http://www.thegnar.co/logo.png',
      });
    });
  });

  describe('submitting the form', () => {
    it('submits the form when valid', () => {
      const formData = {
        org_name: 'The Gnar Co.',
        org_web_url: 'http://www.thegnar.co',
        org_logo_url: 'http://www.thegnar.co/logo.png',
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

