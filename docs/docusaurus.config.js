// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion.

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'WhaleTUI Documentation',
  tagline: 'Docker Management Made Simple',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://wikczerski.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment from gh-pages branch, use root path
  baseUrl: '/whaletui/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'wikczerski', // Usually your GitHub org/user name.
  projectName: 'whaletui', // Usually your repo name.
  deploymentBranch: 'gh-pages',
  trailingSlash: false,

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/wikczerski/whaletui/edit/main/docs/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/wikczerski/whaletui/edit/main/docs/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/theme-common').UserThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/social_card.png',

      // Force dark theme only
      colorMode: {
        defaultMode: 'dark',
        disableSwitch: true,
        respectPrefersColorScheme: false,
      },

      navbar: {
        title: 'WhaleTUI',
        logo: {
          alt: 'WhaleTUI Logo',
          src: 'img/logo.webp',
          srcDark: 'img/logo.webp',
          width: 32,
          height: 32,
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'docs',
            position: 'left',
            label: 'Documentation',
          },
          {to: '/blog', label: 'Blog', position: 'left'},
          {
            href: 'https://github.com/wikczerski/whaletui',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Getting Started',
                to: '/docs/intro',
              },
              {
                label: 'Installation',
                to: '/docs/installation',
              },
            ],
          },
                      {
              title: 'Community',
              items: [
                {
                  label: 'GitHub',
                  href: 'https://github.com/wikczerski/whaletui',
                },
                {
                  label: 'Issues',
                  href: 'https://github.com/wikczerski/whaletui/issues',
                },
              ],
            },
          {
            title: 'More',
            items: [
              {
                label: 'Blog',
                to: '/blog',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} WhaleTUI. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
