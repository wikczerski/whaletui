import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    icon: '🐳',
    title: 'Container Management',
    description: 'Start, stop, restart, and manage containers with ease. Full lifecycle management through an intuitive TUI interface.',
    features: ['Container lifecycle operations', 'Resource inspection', 'Port mapping display', 'Volume mount info'],
    color: '#2496ed'
  },
  {
    icon: '📋',
    title: 'Container Logs & Inspection',
    description: 'View container logs and inspect container details with comprehensive information display.',
    features: ['Container logs viewing', 'Detailed inspection', 'Resource information', 'Configuration details'],
    color: '#17a2b8'
  },
  {
    icon: '🔒',
    title: 'SSH Integration',
    description: 'Connect to remote Docker hosts securely via SSH. Manage containers across multiple environments.',
    features: ['Secure remote access', 'Multi-host management', 'Key-based authentication', 'SSH tunneling support'],
    color: '#28a745'
  },
  {
    icon: '🖼️',
    title: 'Image Management',
    description: 'Efficiently manage Docker images with built-in tools for listing, inspection, and cleanup.',
    features: ['Image listing', 'Registry integration', 'Layer inspection', 'Cleanup tools'],
    color: '#ffc107'
  },
  {
    icon: '🌐',
    title: 'Network Control',
    description: 'List, inspect, and manage Docker networks with basic networking capabilities.',
    features: ['Network listing', 'Configuration inspection', 'Driver information', 'Network details'],
    color: '#6f42c1'
  },
  {
    icon: '💾',
    title: 'Volume & Storage',
    description: 'Manage persistent data with Docker volumes and basic storage information.',
    features: ['Volume listing', 'Storage information', 'Mount point details', 'Basic cleanup'],
    color: '#fd7e14'
  }
];

function Feature({icon, title, description, features, color}) {
  return (
    <div className={clsx('col col--4', styles.featureCol)}>
      <div className={styles.featureCard}>
        <div className={styles.featureIcon} style={{backgroundColor: color}}>
          <span className={styles.iconText}>{icon}</span>
        </div>
        <div className={styles.featureContent}>
          <h3 className={styles.featureTitle}>{title}</h3>
          <p className={styles.featureDescription}>{description}</p>
          <ul className={styles.featureList}>
            {features.map((feature, index) => (
              <li key={index} className={styles.featureItem}>
                <span className={styles.featureBullet}>•</span>
                {feature}
              </li>
            ))}
          </ul>
        </div>
        <div className={styles.featureHover} style={{backgroundColor: color}}></div>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className={styles.featuresHeader}>
          <h2 className={styles.featuresTitle}>Why Choose WhaleTUI?</h2>
          <p className={styles.featuresSubtitle}>
            Built for developers and DevOps engineers who need powerful, efficient Docker management tools
          </p>
        </div>
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
        <div className={styles.featuresFooter}>
          <div className={styles.featuresStats}>
            <div className={styles.stat}>
              <span className={styles.statNumber}>100%</span>
              <span className={styles.statLabel}>Go Native</span>
            </div>
            <div className={styles.stat}>
              <span className={styles.statNumber}>⚡</span>
              <span className={styles.statLabel}>Fast & Lightweight</span>
            </div>
            <div className={styles.stat}>
              <span className={styles.statNumber}>🔧</span>
              <span className={styles.statLabel}>Core Features Ready</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
