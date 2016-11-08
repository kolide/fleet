import React from 'react';
import { size } from 'lodash';

import BasePageForm from 'components/forms/RegistrationForm/BasePageForm';
import Button from 'components/buttons/Button';
import InputFieldWithIcon from 'components/forms/fields/InputFieldWithIcon';
import validateEquality from 'components/forms/validators/validate_equality';
import validEmail from 'components/forms/validators/valid_email';

class AdminDetails extends BasePageForm {
  valid = () => {
    const clientErrors = {};
    const { errors } = this.state;
    const {
      formData: {
        email,
        full_name: fullName,
        password,
        password_confirmation: passwordConfirmation,
        username,
      },
    } = this.props;

    if (!validEmail(email)) {
      clientErrors.email = 'Email must be a valid email';
    }

    if (!email) {
      clientErrors.email = 'Email must be present';
    }

    if (!fullName) {
      clientErrors.full_name = 'Full name must be present';
    }

    if (!username) {
      clientErrors.username = 'Username must be present';
    }

    if (password && passwordConfirmation && !validateEquality(password, passwordConfirmation)) {
      clientErrors.password_confirmation = 'Password confirmation does not match password';
    }

    if (!password) {
      clientErrors.password = 'Password must be present';
    }

    if (!passwordConfirmation) {
      clientErrors.password_confirmation = 'Password confirmation must be present';
    }

    if (size(clientErrors)) {
      this.setState({
        errors: {
          ...errors,
          ...clientErrors,
        },
      });

      return false;
    }

    return true;
  }

  render () {
    const { formData } = this.props;
    const { errors, onChange, onSubmit } = this;

    return (
      <div>
        <InputFieldWithIcon
          error={errors('full_name')}
          name="full name"
          onChange={onChange('full_name')}
          placeholder="Full Name"
          value={formData.full_name}
        />
        <InputFieldWithIcon
          error={errors('username')}
          iconName="kolidecon-username"
          name="username"
          onChange={onChange('username')}
          placeholder="Username"
          value={formData.username}
        />
        <InputFieldWithIcon
          error={errors('password')}
          iconName="kolidecon-password"
          name="password"
          onChange={onChange('password')}
          placeholder="Password"
          type="password"
          value={formData.password}
        />
        <InputFieldWithIcon
          error={errors('password_confirmation')}
          iconName="kolidecon-password"
          name="password confirmation"
          onChange={onChange('password_confirmation')}
          placeholder="Confirm Password"
          type="password"
          value={formData.password_confirmation}
        />
        <InputFieldWithIcon
          error={errors('email')}
          iconName="kolidecon-email"
          name="email"
          onChange={onChange('email')}
          placeholder="Email"
          value={formData.email}
        />
        <Button
          onClick={onSubmit}
          text="Submit"
          variant="gradient"
        />
      </div>
    );
  }
}

export default AdminDetails;
