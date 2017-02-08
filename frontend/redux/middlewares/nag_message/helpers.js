import moment from 'moment';

const shouldNagUser = ({ license }) => {
  const { allowed_hosts: allowedHosts, expiry, hosts } = license;

  const hostsOverenrolled = hosts > allowedHosts;
  const licenseExpired = moment().isAfter(moment(expiry));

  return hostsOverenrolled || licenseExpired;
};

export default { shouldNagUser };
