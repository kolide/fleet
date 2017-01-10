import React, { Component, PropTypes } from 'react';

import ConfigOptionForm from 'components/forms/ConfigOptionsForm/ConfigOptionForm';
import configOptionInterface from 'interfaces/config_option';

class ConfigOptionsForm extends Component {
  static propTypes = {
    configOptions: PropTypes.arrayOf(configOptionInterface),
  };

  renderConfigOptionForm = (option, idx) => {
    return <ConfigOptionForm key={`config-option-form-${option.id}-${idx}`} formData={option} />;
  }

  render () {
    const { configOptions } = this.props;
    const { renderConfigOptionForm } = this;

    return (
      <div>
        {configOptions.map((option, idx) => {
          return renderConfigOptionForm(option, idx);
        })}
      </div>
    );
  }
}

export default ConfigOptionsForm;
