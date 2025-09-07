# WhaleTUI Documentation

This directory contains the official documentation for WhaleTUI, built with [Docusaurus](https://docusaurus.io/).

## 🚀 Quick Start

### Prerequisites

- Node.js 18 or later
- npm (comes with Node.js)

### Local Development

1. **Install dependencies:**
   ```bash
   cd docs
   npm install
   ```

2. **Start the development server:**
   ```bash
   npm start
   ```
   This starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

3. **Build the documentation:**
   ```bash
   npm run build
   ```
   This command generates static content into the `build` directory and can be served using any static contents hosting service.

### 🛠️ Available Scripts

- `npm start` - Start the development server
- `npm run build` - Build the documentation for production
- `npm run serve` - Serve the built documentation locally
- `npm run deploy` - Deploy to GitHub Pages
- `npm run clear` - Clear the build cache
- `npm run write-translations` - Write translation files
- `npm run write-heading-ids` - Add heading IDs to markdown files

### 📁 Project Structure

```
docs/
├── blog/                    # Blog posts
├── docs/                    # Documentation pages
│   ├── concepts/           # Concept documentation
│   ├── development/        # Development guides
│   ├── installation.md     # Installation guide
│   ├── intro.md           # Introduction
│   └── quick-start.md     # Quick start guide
├── src/                    # Source files
│   ├── components/        # React components
│   ├── css/              # Custom CSS
│   └── pages/            # Additional pages
├── static/                # Static assets
├── docusaurus.config.js   # Docusaurus configuration
├── sidebars.js           # Sidebar configuration
└── package.json          # Dependencies and scripts
```

## 🔄 Automated Deployment

The documentation is automatically deployed to GitHub Pages when:

1. **Changes are pushed to master/main branch** that affect the `docs/` directory
2. **Manual deployment is triggered** via GitHub Actions

### GitHub Actions Workflows

- **`docs-deploy.yml`** - Automatic deployment on docs changes
- **`docs-manual-deploy.yml`** - Manual deployment trigger

### Deployment Process

1. Changes are detected in the `docs/` directory
2. GitHub Actions builds the documentation
3. Built files are deployed to the `gh-pages` branch
4. GitHub Pages serves the documentation from the `gh-pages` branch

## 📝 Writing Documentation

### Adding New Pages

1. Create a new `.md` file in the appropriate directory under `docs/`
2. Add the page to `sidebars.js` to include it in the navigation
3. Use proper frontmatter:

```markdown
---
id: my-page
title: My Page Title
sidebar_position: 1
---

# My Page Title

Content goes here...
```

### Blog Posts

1. Create a new `.md` file in the `blog/` directory
2. Use the format: `YYYY-MM-DD-title.md`
3. Include proper frontmatter:

```markdown
---
slug: my-blog-post
title: My Blog Post
authors: [wikczerski]
tags: [announcement, feature]
---

# My Blog Post

Content goes here...
```

## 🎨 Customization

### Styling

- Custom CSS can be added to `src/css/custom.css`
- The theme is configured in `docusaurus.config.js`

### Configuration

- Main configuration: `docusaurus.config.js`
- Sidebar configuration: `sidebars.js`
- Package configuration: `package.json`

## 🔧 Troubleshooting

### Common Issues

1. **Build fails**: Check Node.js version (requires 18+)
2. **Dependencies issues**: Delete `node_modules` and run `npm install`
3. **Deployment fails**: Check GitHub Pages settings and permissions

### Getting Help

- Check the [Docusaurus documentation](https://docusaurus.io/docs)
- Review the [GitHub Actions logs](https://github.com/wikczerski/whaletui/actions)
- Open an issue in the [WhaleTUI repository](https://github.com/wikczerski/whaletui/issues)

## 📄 License

This documentation is part of the WhaleTUI project and is licensed under the MIT License.
