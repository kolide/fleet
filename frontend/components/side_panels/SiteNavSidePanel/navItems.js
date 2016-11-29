export default (admin) => {
  const adminNavItems = [
    {
      icon: 'admin',
      name: 'Admin',
      path: {
        regex: /^\/admin/,
        location: '/admin/users',
      },
      subItems: [
        {
          name: 'Manage Users',
          path: {
            regex: /\/users/,
            location: '/admin/users',
          },
        },
        {
          name: 'App Settings',
          path: {
            regex: /\/settings/,
            location: '/admin/settings',
          },
        },
      ],
    },
  ];

  const userNavItems = [
    {
      defaultPathname: '/hosts/manage',
      icon: 'hosts',
      name: 'Hosts',
      path: {
        regex: /^\/hosts/,
        location: '/hosts/new',
      },
      subItems: [
        {
          name: 'Manage Hosts',
          path: {
            regex: /\/manage/,
            location: '/hosts/manage',
          },
        },
        {
          name: 'Add Hosts',
          path: {
            regex: /\/new/,
            location: '/hosts/new',
          },
        },
      ],
    },
    {
      defaultPathname: '/queries/results',
      icon: 'query',
      name: 'Query',
      path: {
        regex: /^\/queries/,
        location: '/queries/new',
      },
      subItems: [
        {
          name: 'Manage Queries',
          path: {
            regex: /\/results/,
            location: '/queries/results',
          },
        },
        {
          name: 'New Query',
          path: {
            regex: /\/new/,
            location: '/queries/new',
          },
        },
      ],
    },
    {
      defaultPathname: '/packs/all',
      icon: 'packs',
      name: 'Packs',
      path: {
        regex: /^\/packs/,
        location: '/packs/all',
      },
      subItems: [
        {
          name: 'Manage Packs',
          path: {
            regex: /\/all/,
            location: '/packs/all',
          },
        },
        {
          name: 'New Pack',
          path: {
            regex: /\/new/,
            location: '/packs/new',
          },
        },
      ],
    },
    {
      icon: 'help',
      name: 'Help',
      path: {
        regex: /^\/help/,
      },
      subItems: [],
    },
  ];

  if (admin) {
    return [
      ...userNavItems,
      ...adminNavItems,
    ];
  }

  return userNavItems;
};
