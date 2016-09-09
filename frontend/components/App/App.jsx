import React, { Component, PropTypes } from 'react';
import { Style } from 'radium';
import Footer from './Footer';
import componentStyles from './styles';
import globalStyles from '../../styles/global';

export class App extends Component {
  static propTypes = {
    children: PropTypes.element,
  };

  render () {
    const { children } = this.props;

    return (
      <div>
        <Style rules={globalStyles} />
        {children}
        <Footer />
      </div>
    );
  }
}

export default App;
