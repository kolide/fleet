import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';
import { push } from 'react-router-redux';

import Breadcrumbs from 'pages/RegistrationPage/Breadcrumbs';
import paths from 'router/paths';
import RegistrationForm from 'components/forms/RegistrationForm';
import { setup } from 'redux/nodes/auth/actions';
import { showBackgroundImage } from 'redux/nodes/app/actions';

export class RegistrationPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func.isRequired,
  };

  static defaultProps = {
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = { page: 1 };

    return false;
  }

  componentWillMount () {
    const { dispatch } = this.props;

    dispatch(showBackgroundImage);

    return false;
  }

  onNextPage = () => {
    const { page } = this.state;
    this.setState({ page: page + 1 });

    return false;
  }

  onRegistrationFormSubmit = (formData) => {
    const { dispatch } = this.props;
    const { LOGIN } = paths;

    return dispatch(setup(formData))
      .then(() => { return dispatch(push(LOGIN)); })
      .catch(() => { return false; });
  }

  onSetPage = (page) => {
    this.setState({ page });

    return false;
  }

  render () {
    const { page } = this.state;
    const { onRegistrationFormSubmit, onNextPage, onSetPage } = this;

    return (
      <div>
        <Breadcrumbs onClick={onSetPage} page={page} />
        <RegistrationForm page={page} onNextPage={onNextPage} onSubmit={onRegistrationFormSubmit} />
      </div>
    );
  }
}

export default connect()(RegistrationPage);
