import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import { noop } from 'lodash';
import componentStyles from './styles';

class Dropdown extends Component {
  static propTypes = {
    containerStyles: PropTypes.object,
    fieldName: PropTypes.string,
    initialOption: PropTypes.shape({
      text: PropTypes.string,
      value: PropTypes.string,
    }),
    options: PropTypes.arrayOf(PropTypes.shape({
      text: PropTypes.string,
      value: PropTypes.string,
    })),
    onSelect: PropTypes.func,
  };

  static defaultProps = {
    onSelect: noop,
  };

  constructor (props) {
    super(props);

    const { initialOption, options } = props;

    this.state = {
      expanded: false,
      selectedOption: initialOption || options[0],
    };
  }

  onOptionClick = (selectedOption) => {
    return () => {
      const { fieldName, onSelect } = this.props;
      const { value } = selectedOption;

      this.setState({ selectedOption });
      this.toggleShowOptions();

      onSelect({
        [fieldName]: value,
      });

      return false;
    };
  }

  toggleShowOptions = () => {
    const { expanded } = this.state;

    this.setState({
      expanded: !expanded,
    });
  }

  renderOption = (option) => {
    const { value, text } = option;
    const { optionWrapperStyles } = componentStyles;

    return (
      <div key={value} onClick={this.onOptionClick(option)} style={optionWrapperStyles}>
        {text}
      </div>
    );
  }

  render () {
    const { containerStyles, options } = this.props;
    const { expanded, selectedOption } = this.state;
    const { text } = selectedOption;
    const {
      chevronStyles,
      chevronWrapperStyles,
      optionsWrapperStyles,
      selectedOptionStyles,
      selectedTextStyles,
    } = componentStyles;

    return (
      <div style={[{ position: 'relative' }, containerStyles]}>
        <div onClick={this.toggleShowOptions} style={selectedOptionStyles}>
          <span style={selectedTextStyles}>{text}</span>
          <div style={chevronWrapperStyles}><i className="kolidecon-chevron-bold-down" style={chevronStyles} /></div>
          <div style={{ clear: 'both' }} />
        </div>
        <div style={optionsWrapperStyles(expanded)}>
          {options.map(option => {
            return this.renderOption(option);
          })}
        </div>
      </div>
    );
  }
}

export default radium(Dropdown);
