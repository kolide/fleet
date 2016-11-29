import Kolide from 'kolide';
import reduxConfig from 'redux/nodes/entities/base/reduxConfig';
import schemas from 'redux/nodes/entities/base/schemas';

const { CAMPAIGNS: schema } = schemas;

export default reduxConfig({
  createFunc: Kolide.runQuery,
  entityName: 'campaigns',
  schema,
});

