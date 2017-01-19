import { find } from 'lodash';

import Kolide from '../../../../kolide';
import reduxConfig from '../base/reduxConfig';
import schemas from '../base/schemas';

const { HOSTS: schema } = schemas;

export default reduxConfig({
  destroyFunc: Kolide.hosts.destroy,
  entityName: 'hosts',
  loadAllFunc: Kolide.hosts.loadAll,
  parseEntityFunc: (host) => {
    const { network_interfaces: networkInterfaces } = host;
    const networkInterface = networkInterfaces && find(networkInterfaces, { id: host.primary_ip_id });
    const clockSpeed = host.cpu_brand.split('@ ')[1] || host.cpu_brand.split('@')[1];

    const additionalAttrs = {
      host_cpu: `${host.cpu_physical_cores} x ${clockSpeed}`,
      target_type: 'hosts',
    };

    if (networkInterface) {
      additionalAttrs.host_ip_address = networkInterface.address;
      additionalAttrs.host_mac = networkInterface.mac;
    }

    return {
      ...host,
      ...additionalAttrs,
    };
  },
  schema,
});
