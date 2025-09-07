/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  docs: [
    'intro',
    'installation',
    'quick-start',
    {
      type: 'category',
      label: 'Concepts',
      items: [
        'concepts/containers',
        'concepts/images',
        'concepts/networks',
        'concepts/volumes',
        'concepts/swarm',
        'concepts/nodes',
        'concepts/column-configuration',
        'concepts/configuration-examples',
      ],
    },
    {
      type: 'category',
      label: 'Development',
      items: [
        'development/setup',
        'development/coding-standards',
      ],
    },
  ],
};

module.exports = sidebars;
