import React, { Component, PropTypes } from 'react';
import Select from 'react-select';
import { noop, pick } from 'lodash';

import dropdownOptionInterface from 'interfaces/dropdownOption';
import FormField from 'components/forms/FormField';

const baseClass = 'input-dropdown';

class Dropdown extends Component {
  static propTypes = {
    className: PropTypes.string,
    clearable: PropTypes.bool,
    error: PropTypes.string,
    hint: PropTypes.oneOfType([PropTypes.array, PropTypes.string]),
    label: PropTypes.string,
    name: PropTypes.string,
    onChange: PropTypes.func,
    options: PropTypes.arrayOf(dropdownOptionInterface).isRequired,
    placeholder: PropTypes.string,
    value: PropTypes.string,
  };

  static defaultProps = {
    onChange: noop,
    clearable: false,
    placeholder: 'Select One...',
  };

  handleChange = ({ value }) => {
    const { onChange } = this.props;

    return onChange(value);
  };

  render () {
    const { handleChange } = this;
    const { className, clearable, name, options, placeholder, value } = this.props;

    const formFieldProps = pick(this.props, ['hint', 'label', 'error', 'name']);

    return (
      <FormField {...formFieldProps} type="dropdown">
        <Select
          className={`${baseClass}__select ${className}`}
          name={`${name}-select` || "targets"}
          options={options}
          onChange={handleChange}
          placeholder={placeholder}
          value={value}
          clearable={clearable}
        />
      </FormField>
    );
  }
}

export default Dropdown;
