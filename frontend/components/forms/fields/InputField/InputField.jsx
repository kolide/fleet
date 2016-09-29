import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import componentStyles from './styles';

class InputField extends Component {
  static propTypes = {
    autofocus: PropTypes.bool,
    error: PropTypes.string,
    inputWrapperStyles: PropTypes.object,
    inputOptions: PropTypes.object,
    label: PropTypes.string,
    labelStyles: PropTypes.object,
    name: PropTypes.string,
    onChange: PropTypes.func,
    placeholder: PropTypes.string,
    style: PropTypes.object,
    type: PropTypes.string,
  };

  static defaultProps = {
    autofocus: false,
    inputWrapperStyles: {},
    inputOptions: {},
    label: null,
    labelStyles: {},
    style: {},
    type: 'text',
  };

  constructor (props) {
    super(props);
    this.state = { value: null };
  }

  componentDidMount () {
    const { autofocus } = this.props;
    const { input } = this;

    if (autofocus) input.focus();

    return false;
  }

  onInputChange = (evt) => {
    evt.preventDefault();

    const { value } = evt.target;
    const { onChange } = this.props;

    this.setState({ value });
    return onChange(evt);
  }

  renderLabel = () => {
    const { componentLabelStyles } = componentStyles;
    const { error, label, labelStyles, name } = this.props;

    if (!label) return false;

    return (
      <label htmlFor={name} style={[componentLabelStyles(error), labelStyles]}>
        {error || label}
      </label>
    );
  }

  render () {
    const { error, inputOptions, inputWrapperStyles, name, placeholder, style, type } = this.props;
    const { inputErrorStyles, inputStyles } = componentStyles;
    const { value } = this.state;
    const { onInputChange, renderLabel } = this;

    if (type === 'textarea') {
      return (
        <div style={inputWrapperStyles}>
          {renderLabel()}
          <textarea
            name={name}
            onChange={onInputChange}
            className="input-with-icon"
            placeholder={placeholder}
            ref={(r) => { this.input = r; }}
            style={[inputStyles(type, value), inputErrorStyles(error), style]}
            type={type}
            {...inputOptions}
          />
        </div>
      );
    }

    return (
      <div style={inputWrapperStyles}>
        {renderLabel()}
        <input
          name={name}
          onChange={onInputChange}
          className="input-with-icon"
          placeholder={placeholder}
          ref={(r) => { this.input = r; }}
          style={[inputStyles(type, value), inputErrorStyles(error), style]}
          type={type}
          {...inputOptions}
        />
      </div>
    );
  }
}

export default radium(InputField);

