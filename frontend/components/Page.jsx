import React from 'react';

import Navbar from 'frontend/components/Navbar'
import Footer from 'frontend/components/Footer'

const Page = React.createClass({
  render() {
    return (
      <div>
        <Navbar/>
        <div className={this.props.className}>
          {this.props.children}
        </div>
        <Footer/>
      </div>
    );
  },
})

export default Page