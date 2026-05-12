export type NavItem = {
  label: string;
  href: string;
};

export const site = {
  name: 'DevFolio',
  tagline: 'Developer portfolio platform',
  description: 'Portfolio digitale personalizzabile per sviluppatori, con base modulare e pronta per il white-label.',
  url: process.env.NEXT_PUBLIC_SITE_URL ?? 'http://localhost:3000',
  accent: '#1a56db'
};

export const navItems: NavItem[] = [
  { label: 'Home', href: '/' },
  { label: 'Projects', href: '/projects' },
  { label: 'Blog', href: '/blog' },
  { label: 'About', href: '/about' },
  { label: 'Contact', href: '/contact' }
];

export const socialLinks = [
  { label: 'GitHub', href: 'https://github.com/gigliofr' },
  { label: 'LinkedIn', href: 'https://www.linkedin.com/in/gigliofrancesco/' },
  { label: 'Email', href: 'mailto:hello@devfolio.local' }
];