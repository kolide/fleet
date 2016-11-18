import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';

const fieldNames = ['title', 'description', 'targets'];
const validate = () => { return { valid: true, errors: {} }; };

class PackForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      title: formFieldInterface.isRequired,
      description: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func,
  };

  render () {
    const { fields, handleSubmit } = this.props;

    return (
      <form onSubmit={handleSubmit}>
        <InputField
          {...fields.title}
          placeholder="Query Pack Title"
        />
        <InputField
          {...fields.description}
          label="Description"
          placeholder="Add a description of your query"
          type="textarea"
        />
        <div>
          <Button
            text="Save Query pack"
            type="submit"
            variant="brand"
          />
        </div>
      </form>
    );
  }
}

export default Form(PackForm, {
  fields: fieldNames,
  validate,
});
