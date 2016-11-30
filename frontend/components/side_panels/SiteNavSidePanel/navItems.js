export default (admin) => {
  const adminNavItems = [
    {
      icon: 'admin',
      name: 'Admin',
      path: {
        regex: /^\/admin/,
        location: '/admin/users',
      },
      subItems: [],
    },
  ];

  const userNavItems = [
    {
      defaultPathname: '/hosts/manage',
      icon: 'hosts',
      name: 'Hosts',
      path: {
        regex: /^\/hosts/,
        location: '/hosts/manage',
      },
      subItems: [
        {
          icon: 'hosts',
          name: 'Manage Hosts',
          path: {
            regex: /\/manage/,
            location: '/hosts/manage',
          },
        },
        {
          icon: 'add-plus',
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
        location: '/queries/results',
      },
      subItems: [
        {
          icon: 'query',
          name: 'Manage Queries',
          path: {
            regex: /\/results/,
            location: '/queries/results',
          },
        },
        {
          icon: 'pencil',
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
          icon: 'packs',
          name: 'Manage Packs',
          path: {
            regex: /\/all/,
            location: '/packs/all',
          },
        },
        {
          icon: 'pencil',
          name: 'New Pack',
          path: {
            regex: /\/new/,
            location: '/packs/new',
          },
        },
      ],
    },
    {
      defaultPathname: '/config/options',
      icon: 'config',
      name: 'Config',
      path: {
        regex: /^\/config/,
        location: '/config/options',
      },
      subItems: [
        {
          icon: 'config',
          name: 'Osquery Options',
          path: {
            regex: /\/options/,
            location: '/config/options',
          },
        },
        {
          icon: 'import',
          name: 'Import Config',
          path: {
            regex: /\/import/,
            location: '/config/import',
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
