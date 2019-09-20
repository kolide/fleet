export default {
  ADMIN_DASHBOARD: urlPrefix + '/admin',
  ADMIN_SETTINGS: urlPrefix + '/admin/settings',
  ALL_PACKS: urlPrefix + '/packs/all',
  EDIT_QUERY: (query) => {
    return `${urlPrefix}/queries/${query.id}`;
  },
  FORGOT_PASSWORD: urlPrefix + '/login/forgot',
  HOME: urlPrefix + '/',
  KOLIDE_500: urlPrefix + '/500',
  LOGIN: urlPrefix + '/login',
  LOGOUT: urlPrefix + '/logout',
  MANAGE_HOSTS: urlPrefix + '/hosts/manage',
  NEW_PACK: urlPrefix + '/packs/new',
  NEW_QUERY: urlPrefix + '/queries/new',
  RESET_PASSWORD: urlPrefix + '/login/reset',
  SETUP: urlPrefix + '/setup',
  USER_SETTINGS: urlPrefix + '/settings',
};
