import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Form from 'components/forms/Form';
import InputField from 'components/forms/fields/InputField';
import packInterface from 'interfaces/pack';

const fieldNames = ['description', 'name'];

class EditPackForm extends Component {
  static propTypes = {
    className: PropTypes.string,
    formData: packInterface,
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
