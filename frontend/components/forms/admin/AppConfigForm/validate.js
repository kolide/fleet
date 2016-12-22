import { size } from 'lodash';

export default (formData) => {
  const errors = {};
  const {
    kolide_server_url: kolideServerUrl,
    smtp_sender_address: smtpSenderAddress,
    server: smtpServer,
    user_name: smtpUserName,
    password: smtpPassword,
  } = formData;

  if (!kolideServerUrl) {
    errors.kolide_server_url = 'Kolide Server URL must be present';
  }

  if (!smtpSenderAddress) {
    errors.smtp_sender_address = 'SMTP Sender Address must be present';
  }

  if (!smtpServer) {
    errors.server = 'SMTP Server must be present';
  }

  if (!smtpUserName) {
    errors.user_name = 'SMTP Username must be present';
  }

  if (!smtpPassword) {
    errors.password = 'SMTP Password must be present';
  }

  const valid = !size(errors);

  return { valid, errors };
};
