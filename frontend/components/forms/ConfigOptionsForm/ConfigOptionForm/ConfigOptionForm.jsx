import React, { Component, PropTypes } from 'react';

import Dropdown from 'components/forms/fields/Dropdown';
import dropdownOptionInterface from 'interfaces/dropdownOption';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';

const fieldNames = ['name', 'value'];

class ConfigOptionForm extends Component {
  static propTypes = {
    configNameOptions: PropTypes.arrayOf(dropdownOptionInterface),
    fields: PropTypes.shape({
      name: formFieldInterface,
      value: formFieldInterface,
    }),
  };

  render () {
    const { configNameOptions, fields } = this.props;

    return (
      <form>
        <Dropdown
          {...fields.name}
          options={configNameOptions}
        />
        <InputField {...fields.value} />
      </form>
    );
  }
}

export default Form(ConfigOptionForm, { fields: fieldNames });
