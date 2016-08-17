import React from 'react';

import Navbar from '#components/Navbar'
import Footer from '#components/Footer'

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