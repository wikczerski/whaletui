# WhaleTUI Documentation

This directory contains the official documentation for WhaleTUI, built with [Docusaurus](https://docusaurus.io/).

## What is Docusaurus?

Docusaurus is a modern static website generator that makes it easy to build, deploy, and maintain open source project websites.

## Project Structure

```
docs/
├── blog/                    # Blog posts
├── docs/                    # Documentation pages
│   ├── concepts/           # Core Docker concepts
│   ├── ui/                 # User interface guides
│   ├── advanced/           # Advanced topics
│   └── development/        # Development guides
├── src/                    # React components and styles
│   ├── components/         # Custom React components
│   ├── css/               # Custom CSS styles
│   └── pages/             # Additional pages
├── static/                 # Static assets (images, etc.)
├── docusaurus.config.js    # Main configuration
├── sidebars.js            # Sidebar navigation
├── package.json           # Node.js dependencies
└── README.md              # This file
```

## Getting Started

### Prerequisites

- Node.js 16+
- npm or yarn

### Installation

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run start
   ```

3. Open [http://localhost:3000](http://localhost:3000) in your browser.

### Building for Production

```bash
npm run build
```

The built files will be in the `build/` directory.

### Deployment

The site is configured for GitHub Pages deployment. The build output can be deployed to the `gh-pages` branch.

## Development

### Adding New Documentation

1. Create new `.md` files in the appropriate `docs/` subdirectory
2. Add frontmatter with proper metadata
3. Update `sidebars.js` to include new pages
4. Test locally with `npm run start`

### Adding New Blog Posts

1. Create new `.md` files in the `blog/` directory
2. Use the date format: `YYYY-MM-DD-title.md`
3. Include proper frontmatter with title, authors, and tags

### Customizing the Theme

- Modify `src/css/custom.css` for global styles
- Update `docusaurus.config.js` for theme configuration
- Create custom React components in `src/components/`

## Configuration

### Docusaurus Configuration

The main configuration is in `docusaurus.config.js`:

- Site metadata and branding
- Navigation and footer links
- Plugin configurations
- Theme settings

### Sidebar Configuration

Navigation structure is defined in `sidebars.js`:

- Document organization
- Category grouping
- Page ordering

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally
5. Submit a pull request

## Resources

- [Docusaurus Documentation](https://docusaurus.io/docs)
- [Docusaurus Tutorial](https://docusaurus.io/docs/tutorial)
- [Markdown Guide](https://www.markdownguide.org/)
- [React Documentation](https://react.dev/)

## Support

If you have questions about the documentation:

- Check the [Docusaurus documentation](https://docusaurus.io/docs)
- Open an issue on GitHub
- Join our community discussions

---

*Built with ❤️ using [Docusaurus](https://docusaurus.io/)*
