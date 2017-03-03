import { schema } from 'normalizr';

const campaignsSchema = new schema.Entity('campaigns');
const configOptionsSchema = new schema.Entity('config_options');
const hostsSchema = new schema.Entity('hosts');
const invitesSchema = new schema.Entity('invites');
const labelsSchema = new schema.Entity('labels');
const packsSchema = new schema.Entity('packs');
const queriesSchema = new schema.Entity('queries');
const scheduledQueriesSchema = new schema.Entity('scheduled_queries');
const targetsSchema = new schema.Entity('targets');
const usersSchema = new schema.Entity('users');

export default {
  CAMPAIGNS: campaignsSchema,
  CONFIG_OPTIONS: configOptionsSchema,
  HOSTS: hostsSchema,
  INVITES: invitesSchema,
  LABELS: labelsSchema,
  PACKS: packsSchema,
  QUERIES: queriesSchema,
  SCHEDULED_QUERIES: scheduledQueriesSchema,
  TARGETS: targetsSchema,
  USERS: usersSchema,
};
