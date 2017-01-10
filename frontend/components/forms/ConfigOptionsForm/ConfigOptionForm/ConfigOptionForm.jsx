import React, { Component, PropTypes } from 'react';

import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';

const fieldNames = ['name', 'value'];

class ConfigOptionForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      name: formFieldInterface,
      value: formFieldInterface,
    }),
  };

  render () {
    const { fields } = this.props;

    return (
      <form>
        <InputField {...fields.name} />
        <InputField {...fields.value} />
      </form>
    );
  }
}

export default Form(ConfigOptionForm, { fields: fieldNames });
