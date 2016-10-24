import Kolide from '../../../../kolide';
import { parseTarget } from './helpers';
import reduxConfig from '../base/reduxConfig';
import schemas from '../base/schemas';

const { TARGETS: schema } = schemas;

export default reduxConfig({
  entityName: 'targets',
  loadAllFunc: Kolide.getTargets,
  parseApiResponseFunc: (response) => { return response.targets; },
  parseEntityFunc: parseTarget,
  schema,
});
