import React, { PropTypes } from 'react';
import radium from 'radium';
import componentStyles from './styles';
import InputField from '../InputField';

class InputFieldWithLabel extends InputField {
  static propTypes = {
    autofocus: PropTypes.bool,
    defaultValue: PropTypes.string,
    error: PropTypes.string,
    iconName: PropTypes.string,
    label: PropTypes.string,
    name: PropTypes.string,
    onChange: PropTypes.func,
    placeholder: PropTypes.string,
    style: PropTypes.object,
    type: PropTypes.string,
  };

  constructor(props) {
    super(props);

    const { defaultValue } = this.props;

    this.state = {
      value: defaultValue,
    };
  }

  render () {
    const { label, name, placeholder, style, type } = this.props;
    const { containerStyles, inputStyles, labelStyles } = componentStyles;
    const { value } = this.state;
    const { onInputChange } = this;

    return (
      <div style={[containerStyles, style.container]}>
        <label htmlFor={name} style={labelStyles}>{label}</label>
        <input
          name={name}
          onChange={onInputChange}
          placeholder={placeholder}
          ref={(r) => { this.input = r; }}
          style={[inputStyles(value, type), style.input]}
          type={type}
          value={value}
        />
      </div>
    );
  }
}

export default radium(InputFieldWithLabel);

