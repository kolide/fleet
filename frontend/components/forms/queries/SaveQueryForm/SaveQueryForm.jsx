import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import componentStyles from './styles';
import InputFieldWithIcon from '../../fields/InputFieldWithIcon';
import GradientButton from '../../../buttons/GradientButton';

class SaveQueryForm extends Component {
  static propTypes = {
    onSubmit: PropTypes.func,
  };

  constructor (props) {
    super(props);

    this.state = {
      formData: {
        queryName: null,
      },
    };
  }

  onFieldChange = (fieldName) => {
    return ({ target }) => {
      const { formData } = this.state;

      this.setState({
        formData: {
          ...formData,
          [fieldName]: target.value,
        },
      });
    };
  }

  onFormSubmit = (evt) => {
    evt.preventDefault();

    const { formData } = this.state;
    const { onSubmit } = this.props;

    return onSubmit(formData);
  }

  render () {
    const { buttonStyles } = componentStyles;
    const { onFieldChange, onFormSubmit } = this;

    return (
      <form onSubmit={onFormSubmit}>
        <label htmlFor="queryName">Query Name</label>
        <InputFieldWithIcon onChange={onFieldChange('queryName')} name="queryName" />
        <GradientButton
          style={buttonStyles}
          text="Run Query"
          type="submit"
        />
      </form>
    );
  }
}

export default radium(SaveQueryForm);
