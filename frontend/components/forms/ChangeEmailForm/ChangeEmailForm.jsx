import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';

class ChangeEmailForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      password: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func.isRequired,
  };

  render () {
    const { fields, handleSubmit } = this.props;

    return (
      <form onSubmit={handleSubmit}>
        <InputField
          {...fields.password}
          autofocus
          label="Password"
          type="password"
        />
        <Button block type="submit" variant="brand">
          Submit
        </Button>
      </form>
    );
  }
}

export default Form(ChangeEmailForm, {
  fields: ['password'],
  validate: (formData) => {
    if (!formData.password) {
      return {
        valid: false,
        errors: { password: 'Password must be present' },
      };
    }

    return { valid: true, errors: {} };
  },
});
