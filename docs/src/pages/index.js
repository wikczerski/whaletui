import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <div className="row">
          <div className="col col--8">
            <h1 className="hero__title">
              <span className={styles.heroTitleMain}>WhaleTUI</span>
              <span className={styles.heroTitleSub}>Docker Management Made Simple</span>
            </h1>
            <p className={styles.hero__subtitle}>
              A powerful, terminal-based Docker management tool built with Go.
              Manage containers, images, networks, and volumes with an intuitive TUI interface.
              Connect to remote hosts via SSH for distributed Docker management.
            </p>
            <div className={styles.buttons}>
              <Link
                className="button button--secondary button--lg"
                to="/docs/installation">
                Get Started
              </Link>
              <Link
                className="button button--outline button--lg"
                to="/docs/quick-start">
                Quick Start Guide
              </Link>
            </div>
          </div>
          <div className="col col--4">
            <div className={styles.heroImage}>
              <img
                src={require('@site/static/img/logo.webp').default}
                alt="WhaleTUI Logo"
                className={styles.logo}
              />
              <div className={styles.terminalMockup}>
                <div className={styles.terminalHeader}>
                  <span className={styles.terminalDot}></span>
                  <span className={styles.terminalDot}></span>
                  <span className={styles.terminalDot}></span>
                </div>
                <div className={styles.terminalContent}>
                  <span className={styles.prompt}>$</span> whaletui
                  <br />
                  <span className={styles.success}>✓</span> Connected to Docker
                  <br />
                  <span className={styles.info}>→</span> Managing 3 containers
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}

function HomepageStats() {
  return (
    <section className={styles.stats}>
      <div className="container">
        <div className="row">
          <div className="col col--3">
            <div className={styles.statItem}>
              <div className={styles.statNumber}>🚀</div>
              <div className={styles.statLabel}>Fast Performance</div>
              <div className={styles.statDesc}>Built with Go for optimal speed</div>
            </div>
          </div>
          <div className="col col--3">
            <div className={styles.statItem}>
              <div className={styles.statNumber}>🎯</div>
              <div className={styles.statLabel}>Intuitive TUI</div>
              <div className={styles.statDesc}>Clean, organized interface</div>
            </div>
          </div>
          <div className="col col--3">
            <div className={styles.statItem}>
              <div className={styles.statNumber}>🔒</div>
              <div className={styles.statLabel}>SSH Integration</div>
              <div className={styles.statDesc}>Secure remote management</div>
            </div>
          </div>
          <div className="col col--3">
            <div className={styles.statItem}>
              <div className={styles.statNumber}>📊</div>
              <div className={styles.statLabel}>Container Inspection</div>
              <div className={styles.statDesc}>Detailed container information</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

function HomepageCodeExample() {
  return (
    <section className={styles.codeExample}>
      <div className="container">
        <div className="row">
          <div className="col col--6">
            <h2>Simple Installation</h2>
            <p>Get WhaleTUI running on your system in minutes with our easy installation process.</p>
            <div className={styles.codeBlock}>
              <div className={styles.codeHeader}>Install with Go</div>
              <pre className={styles.code}>
                <code>go install github.com/wikczerski/whaletui@latest</code>
              </pre>
            </div>
            <div className={styles.codeBlock}>
              <div className={styles.codeHeader}>Or download binary</div>
              <pre className={styles.code}>
                <code>curl -L https://github.com/wikczerski/whaletui/releases/latest/download/whaletui-v1.0.0-linux-amd64 -o whaletui</code>
              </pre>
              <div className={styles.codeNote}>
                <small>Check <a href="https://github.com/wikczerski/whaletui/releases" target="_blank" rel="noopener">GitHub releases</a> for actual filenames</small>
              </div>
            </div>
          </div>
          <div className="col col--6">
            <h2>Quick Usage</h2>
            <p>Start managing your Docker containers with simple commands.</p>
            <div className={styles.codeBlock}>
              <div className={styles.codeHeader}>Launch WhaleTUI</div>
              <pre className={styles.code}>
                <code>whaletui</code>
              </pre>
            </div>
            <div className={styles.codeBlock}>
              <div className={styles.codeHeader}>Connect to remote host</div>
              <pre className={styles.code}>
                <code>whaletui connect --host ssh://user@remote-server</code>
              </pre>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

function HomepageCTA() {
  return (
    <section className={styles.cta}>
      <div className="container">
        <div className="row">
          <div className="col col--8 col--offset-2 text--center">
            <h2>Ready to Simplify Your Docker Workflow?</h2>
            <p>
              Join developers and DevOps engineers who are already using WhaleTUI
              to manage their Docker environments with a clean, efficient TUI interface.
              Perfect for both local development and remote server management.
            </p>
            <div className={styles.ctaButtons}>
              <Link
                className="button button--primary button--lg"
                to="/docs/installation">
                Install WhaleTUI
              </Link>
              <Link
                className="button button--outline button--lg"
                to="https://github.com/wikczerski/whaletui">
                View on GitHub
              </Link>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} - ${siteConfig.tagline}`}
      description="A powerful, terminal-based Docker management tool built with Go. Manage containers, images, networks, and volumes with an intuitive TUI interface.">
      <HomepageHeader />
      <main>
        <HomepageStats />
        <HomepageFeatures />
        <HomepageCodeExample />
        <HomepageCTA />
      </main>
    </Layout>
  );
}
