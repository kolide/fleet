import stateEntityGetter from 'react-entity-getter';

const pathToEntities = (entityName) => {
  return `entities[${entityName}].entities`;
};

export default stateEntityGetter(pathToEntities);
