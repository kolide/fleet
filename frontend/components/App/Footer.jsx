import React from 'react';
import radium from 'radium';
import componentStyles from './styles';
import footerLogo from './footer-logo.svg'

const { footerStyles } = componentStyles;

const Footer = () => {

  return (
    <footer style={footerStyles}>
      <img src={footerLogo}/>
    </footer>
  );
};

export default radium(Footer);
