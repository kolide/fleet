import BaseConfig from 'redux/nodes/entities/base/base_config';

class ReduxConfig extends BaseConfig {
  constructor (inputs) {
    super(inputs);
  }

  get actions () {
    return this.allActions();
  }
};

export default ReduxConfig;
