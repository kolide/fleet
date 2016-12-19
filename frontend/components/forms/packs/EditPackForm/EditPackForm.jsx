import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';

const fieldNames = ['description', 'name'];

class EditPackForm extends Component {
  static propTypes = {
    className: PropTypes.string,
    fields: PropTypes.arrayOf(formFieldInterface).isRequired,
    handleSubmit: PropTypes.func.isRequired,
  };

  render () {
    const { className, fields, handleSubmit } = this.props;

    return (
      <form className={className} onSubmit={handleSubmit}>
        <InputField
          {...fields.name}
        />
        <InputField
          {...fields.description}
        />
        <Button
          text="SAVE"
          type="submit"
          variant="brand"
        />
      </form>
    );
  }
}

export default Form(EditPackForm, {
  fields: fieldNames,
});
