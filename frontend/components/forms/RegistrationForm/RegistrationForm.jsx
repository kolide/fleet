import React, { Component, PropTypes } from 'react';

import AdminDetails from 'components/forms/RegistrationForm/AdminDetails';
import OrgDetails from 'components/forms/RegistrationForm/OrgDetails';

class RegistrationForm extends Component {
  static propTypes = {
    onNextPage: PropTypes.func,
    onSubmit: PropTypes.func,
    page: PropTypes.number,
  };

  constructor (props) {
    super(props);

    this.state = {
      errors: {},
      formData: {
        full_name: '',
        username: '',
        password: '',
        password_confirmation: '',
        email: '',
        org_name: '',
        org_web_url: '',
        org_logo_url: '',
      },
    };
  }

  onInputFieldChange = (inputData) => {
    const { errors, formData } = this.state;

    this.setState({
      errors: {
        ...errors,
        [Object.keys(formData)]: null,
      },
      formData: {
        ...formData,
        ...inputData,
      },
    });

    return false;
  }

  onPageFormSubmit = () => {
    const { onNextPage } = this.props;

    return onNextPage();
  }

  onSubmit = (evt) => {
    evt.preventDefault();

    const { formData } = this.state;
    const { onSubmit: handleSubmit } = this.props;

    return handleSubmit(formData);
  }

  renderPageForm = () => {
    const { errors, formData } = this.state;
    const { onInputFieldChange, onPageFormSubmit } = this;
    const { page } = this.props;

    if (page === 1) {
      return (
        <AdminDetails
          errors={errors}
          formData={formData}
          onChange={onInputFieldChange}
          onSubmit={onPageFormSubmit}
        />
      );
    }

    if (page === 2) {
      return (
        <OrgDetails
          errors={errors}
          formData={formData}
          onChange={onInputFieldChange}
          onSubmit={onPageFormSubmit}
        />
      );
    }

    return false;
  }

  render () {
    const { onSubmit } = this.props;
    const { renderPageForm } = this;

    return (
      <form onSubmit={onSubmit}>
        {renderPageForm()}
      </form>
    );
  }
}

export default RegistrationForm;
