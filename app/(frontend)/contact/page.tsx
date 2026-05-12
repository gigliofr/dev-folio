import type { Metadata } from 'next';
import { ContactForm } from './contact-form';
import { site } from '@/lib/site';

export const metadata: Metadata = {
  title: 'Contact',
  description: 'Start the conversation for the next iteration.',
  openGraph: {
    title: `Contact | ${site.name}`,
    description: 'Start the conversation for the next iteration.'
  }
};

export default function ContactPage() {
  return (
    <section className="container-shell py-16 md:py-24">
      <div className="max-w-2xl space-y-6">
        <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[var(--accent)]">Contact</p>
        <h1 className="text-4xl font-semibold tracking-tight md:text-6xl">Start the conversation for the next iteration.</h1>
        <ContactForm />
      </div>
    </section>
  );
}