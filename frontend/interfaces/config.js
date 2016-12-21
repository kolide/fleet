import { PropTypes } from 'react';

export default PropTypes.shape({
  authentication_method: PropTypes.string,
  authentication_type: PropTypes.string,
  domain: PropTypes.string,
  enable_sll_tls: PropTypes.bool,
  enable_start_tls: PropTypes.bool,
  kolide_server_url: PropTypes.string,
  org_logo_url: PropTypes.string,
  org_name: PropTypes.string,
  password: PropTypes.string,
  port: PropTypes.string,
  sender_address: PropTypes.string,
  server: PropTypes.string,
  smtp_configured: PropTypes.bool,
  user_name: PropTypes.string,
  verify_sll_certs: PropTypes.bool,
});
