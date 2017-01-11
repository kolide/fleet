import React, { Component, PropTypes } from 'react';

import Dropdown from 'components/forms/fields/Dropdown';
import dropdownOptionInterface from 'interfaces/dropdownOption';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';

const baseClass = 'config-option-form';
const fieldNames = ['name', 'value'];

class ConfigOptionForm extends Component {
  static propTypes = {
    configNameOptions: PropTypes.arrayOf(dropdownOptionInterface),
    fields: PropTypes.shape({
      name: formFieldInterface,
      value: formFieldInterface,
    }),
    formData: PropTypes.shape({
      read_only: PropTypes.bool,
    }).isRequired,
  };

  render () {
    const { configNameOptions, fields, formData } = this.props;
    const { read_only: readOnly } = formData;

    return (
      <form className={baseClass}>
        <Dropdown
          {...fields.name}
          className={`${baseClass}__field`}
          disabled={readOnly}
          options={configNameOptions}
        />
        <InputField
          {...fields.value}
          disabled={readOnly}
          inputClassName={`${baseClass}__field`}
        />
      </form>
    );
  }
}

export default Form(ConfigOptionForm, { fields: fieldNames });
