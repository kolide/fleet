import React, { Component, PropTypes } from 'react';

import ConfigOptionForm from 'components/forms/ConfigOptionsForm/ConfigOptionForm';
import configOptionInterface from 'interfaces/config_option';
import dropdownOptionInterface from 'interfaces/dropdownOption';

class ConfigOptionsForm extends Component {
  static propTypes = {
    completedOptions: PropTypes.arrayOf(configOptionInterface),
    configNameOptions: PropTypes.arrayOf(dropdownOptionInterface),
    onRemoveOption: PropTypes.func.isRequired,
  };

  renderConfigOptionForm = (option, idx) => {
    const { configNameOptions, onRemoveOption } = this.props;

    return (
      <ConfigOptionForm
        configNameOptions={configNameOptions}
        formData={option}
        key={`config-option-form-${option.id}-${idx}`}
        onRemove={onRemoveOption}
      />
    );
  }

  render () {
    const { completedOptions } = this.props;
    const { renderConfigOptionForm } = this;

    return (
      <div>
        {completedOptions.map((option, idx) => {
          return renderConfigOptionForm(option, idx);
        })}
      </div>
    );
  }
}

export default ConfigOptionsForm;
