import React from 'react';
import expect, { createSpy, restoreSpies } from 'expect';
import { mount } from 'enzyme';
import { noop } from 'lodash';

import AdminDetails from 'components/forms/RegistrationForm/AdminDetails';
import { fillInFormInput } from 'test/helpers';

const noErrors = {};

describe('AdminDetails - form', () => {
  afterEach(restoreSpies);

  describe('full name input', () => {
    it('renders an input field', () => {
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={noop} />
      );
      const fullNameField = form.find({ name: 'full name' });

      expect(fullNameField.length).toEqual(1);
    });

    it('calls the onChange prop when the field changes', () => {
      const onChangeSpy = createSpy();
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={onChangeSpy} />
      );
      const fullNameField = form.find({ name: 'full name' }).find('input');

      fillInFormInput(fullNameField, 'The Gnar Co');

      expect(onChangeSpy).toHaveBeenCalledWith('full_name', 'The Gnar Co');
    });
  });

  describe('username input', () => {
    it('renders an input field', () => {
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={noop} />
      );
      const usernameField = form.find({ name: 'username' });

      expect(usernameField.length).toEqual(1);
    });

    it('calls the onChange prop when the field changes', () => {
      const onChangeSpy = createSpy();
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={onChangeSpy} />
      );
      const usernameField = form.find({ name: 'username' }).find('input');

      fillInFormInput(usernameField, 'Gnar');

      expect(onChangeSpy).toHaveBeenCalledWith('username', 'Gnar');
    });
  });

  describe('password input', () => {
    it('renders an input field', () => {
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={noop} />
      );
      const passwordField = form.find({ name: 'password' });

      expect(passwordField.length).toEqual(1);
    });

    it('calls the onChange prop when the field changes', () => {
      const onChangeSpy = createSpy();
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={onChangeSpy} />
      );
      const passwordField = form.find({ name: 'password' }).find('input');

      fillInFormInput(passwordField, 'p@ssw0rd');

      expect(onChangeSpy).toHaveBeenCalledWith('password', 'p@ssw0rd');
    });
  });

  describe('password confirmation input', () => {
    it('renders an input field', () => {
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={noop} />
      );
      const passwordField = form.find({ name: 'password confirmation' });

      expect(passwordField.length).toEqual(1);
    });

    it('calls the onChange prop when the field changes', () => {
      const onChangeSpy = createSpy();
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={onChangeSpy} />
      );
      const passwordField = form.find({ name: 'password confirmation' }).find('input');

      fillInFormInput(passwordField, 'p@ssw0rd');

      expect(onChangeSpy).toHaveBeenCalledWith('password_confirmation', 'p@ssw0rd');
    });
  });

  describe('email input', () => {
    it('renders an input field', () => {
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={noop} />
      );
      const emailField = form.find({ name: 'email' });

      expect(emailField.length).toEqual(1);
    });

    it('calls the onChange prop when the field changes', () => {
      const onChangeSpy = createSpy();
      const form = mount(
        <AdminDetails errors={noErrors} formData={{}} onChange={onChangeSpy} />
      );
      const emailField = form.find({ name: 'email' }).find('input');

      fillInFormInput(emailField, 'hi@gnar.dog');

      expect(onChangeSpy).toHaveBeenCalledWith('email', 'hi@gnar.dog');
    });
  });

  describe('submitting the form', () => {
    it('validates the email field', () => {
      const onSubmitSpy = createSpy();
      const form = mount(
        <AdminDetails
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
        email: 'Email must be present',
        full_name: 'Full name must be present',
        password: 'Password must be present',
        password_confirmation: 'Password confirmation must be present',
        username: 'Username must be present',
      });
    });

    it('validates the email field', () => {
      const onSubmitSpy = createSpy();
      const form = mount(
        <AdminDetails
          errors={noErrors}
          formData={{ email: 'invalid-email' }}
          onChange={noop}
          onSubmit={onSubmitSpy}
        />
      );
      const submitBtn = form.find('Button');


      submitBtn.simulate('click');

      expect(onSubmitSpy).toNotHaveBeenCalled();
      expect(form.state().errors).toInclude({
        email: 'Email must be a valid email',
      });
    });

    it('validates the password fields match', () => {
      const onSubmitSpy = createSpy();
      const formData = {
        password: 'p@ssw0rd',
        password_confirmation: 'password123',
      };
      const form = mount(
        <AdminDetails
          errors={noErrors}
          formData={formData}
          onChange={noop}
          onSubmit={onSubmitSpy}
        />
      );
      const submitBtn = form.find('Button');

      submitBtn.simulate('click');

      expect(onSubmitSpy).toNotHaveBeenCalled();
      expect(form.state().errors).toInclude({
        password_confirmation: 'Password confirmation does not match password',
      });
    });

    it('submits the form when valid', () => {
      const onSubmitSpy = createSpy();
      const formData = {
        email: 'hi@gnar.dog',
        full_name: 'Gnar Dog',
        password: 'p@ssw0rd',
        password_confirmation: 'p@ssw0rd',
        username: 'gnardog',
      };
      const form = mount(
        <AdminDetails
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
