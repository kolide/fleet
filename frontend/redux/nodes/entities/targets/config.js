import Kolide from '../../../../kolide';
import { parseTarget } from './helpers';
import reduxConfig from '../base/reduxConfig';
import schemas from '../base/schemas';

const { TARGETS: schema } = schemas;

export default reduxConfig({
  entityName: 'targets',
  loadAllFunc: Kolide.getTargets,
  parseFunc: parseTarget,
  schema,
});
