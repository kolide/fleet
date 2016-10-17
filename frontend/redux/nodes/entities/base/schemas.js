import { Schema } from 'normalizr';

const hostsSchema = new Schema('hosts');
const invitesSchema = new Schema('invites');
const labelsSchema = new Schema('labels');
const usersSchema = new Schema('users');

export default {
  HOSTS: hostsSchema,
  INVITES: invitesSchema,
  LABELS: labelsSchema,
  USERS: usersSchema,
};
