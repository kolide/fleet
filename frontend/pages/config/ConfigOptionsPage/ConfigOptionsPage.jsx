import React, { Component, PropTypes } from 'react';
import { isEqual } from 'lodash';

import Button from 'components/buttons/Button';
import configOptionInterface from 'interfaces/config_option';

const baseClass= "config-options-page";

class ConfigOptionsPage extends Component {
  static propTypes = {
    config_options: PropTypes.arrayOf(configOptionInterface),
  };

  static defaultProps = {
    configOptions: [],
  };

  constructor (props) {
    const { config_options: configOptions } = props;

    super(props);

    this.state = {
      configOptions: configOptions || [],
    };
  }

  componentWillReceiveProps ({ config_options: configOptions }) {
    if (!isEqual(configOptions, this.state.configOptions)) {
      this.setState({ configOptions });
    }

    return false;
  }

  renderOptions = () => {
    return (
      <div className={`${baseClass}__options-wrapper`}>
      </div>
    );
  }

  render () {
    const { renderOptions } = this;

    return (
      <div className={`body-wrap ${baseClass} ${baseClass}__wrapper`}>
        <div>
          <h1>Manage Additional Osquery Options</h1>
          <p>
            Osquery allows you to set a number of configuration options (Osquery Documentation).
            Since Kolide manages your Osquery configuration, you can set these additional desired
            options on this screen. Some options that Kolide needs to function correctly will be ignored.
          </p>
        </div>
        <div className={`${baseClass}__btn-wrapper`}>
          <Button block className={`${baseClass}__reset-btn`} variant="inverse">
            RESET TO DEFAULT
          </Button>
          <Button block className={`${baseClass}__save-btn`} variant="brand">
            SAVE OPTIONS
          </Button>
        </div>
        {renderOptions()}
      </div>
    );
  }
}

export default ConfigOptionsPage;
