/*
 * Copyright 2016-present, Kolide, Inc.
 * All rights reserved.
 *
 * @flow
 */

import React from 'react';

import Navbar from 'frontend/components/Navbar'
import Footer from 'frontend/components/Footer'

/**
 * Page is the parent component for all pages. It wraps the necessary layout
 * required to render a page properly and consistently.
 *
 * @exports Page
 */
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