import { Component, PropTypes } from 'react';

const formInterface = PropTypes.shape({
  fullName: PropTypes.string,
  username: PropTypes.string,
});

class BasePageForm extends Component {
  static propTypes = {
    errors: formInterface.isRequired,
    formData: formInterface.isRequired,
    onChange: PropTypes.func.isRequired,
    onSubmit: PropTypes.func.isRequired,
  };

  constructor (props) {
    super(props);

    this.state = {
      errors: {},
    };
  }

  onChange = (fieldName) => {
    return (value) => {
      const { errors } = this.state;
      const { onChange: handleChange } = this.props;

      this.setState({
        errors: {
          ...errors,
          [fieldName]: null,
        },
      });

      return handleChange({ [fieldName]: value });
    };
  }

  onSubmit = (evt) => {
    evt.preventDefault();

    const { valid } = this;

    if (valid()) {
      const { onSubmit: handleSubmit } = this.props;

      return handleSubmit();
    }

    return false;
  }

  errors = (fieldName) => {
    const { errors } = this.state;
    const { errors: serverErrors } = this.props;

    return errors[fieldName] || serverErrors[fieldName];
  }
}

export default BasePageForm;

