'use client';

import { useState } from 'react';
import { submitContact } from '@/lib/backend';

export function ContactForm() {
  const [formData, setFormData] = useState({ name: '', email: '', message: '' });
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setMessage(null);

    const success = await submitContact(formData.name, formData.email, formData.message);

    if (success) {
      setMessage({ type: 'success', text: 'Thank you! We received your message.' });
      setFormData({ name: '', email: '', message: '' });
    } else {
      setMessage({ type: 'error', text: 'Failed to send message. Please try again later.' });
    }

    setIsLoading(false);
  };

  return (
    <div className="rounded-3xl border border-[var(--border)] bg-[var(--surface)] p-8">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="mb-2 block text-sm font-medium text-[var(--text-primary)]">Name</label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            placeholder="Your name"
            className="w-full rounded-lg border border-[var(--border)] bg-[var(--bg)] px-4 py-2 text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)]"
            required
          />
        </div>

        <div>
          <label className="mb-2 block text-sm font-medium text-[var(--text-primary)]">Email</label>
          <input
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            placeholder="your.email@example.com"
            className="w-full rounded-lg border border-[var(--border)] bg-[var(--bg)] px-4 py-2 text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)]"
            required
          />
        </div>

        <div>
          <label className="mb-2 block text-sm font-medium text-[var(--text-primary)]">Message</label>
          <textarea
            name="message"
            value={formData.message}
            onChange={handleChange}
            placeholder="Your message..."
            rows={5}
            className="w-full rounded-lg border border-[var(--border)] bg-[var(--bg)] px-4 py-2 text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)]"
            required
          />
        </div>

        <button
          type="submit"
          disabled={isLoading}
          className="w-full rounded-lg bg-[var(--accent)] px-6 py-3 font-medium text-[var(--bg)] transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          {isLoading ? 'Sending...' : 'Send Message'}
        </button>
      </form>

      {message && (
        <div
          className={`mt-4 rounded-lg border p-4 text-sm ${
            message.type === 'success'
              ? 'border-green-900/50 bg-green-900/20 text-green-400'
              : 'border-red-900/50 bg-red-900/20 text-red-400'
          }`}
        >
          {message.text}
        </div>
      )}
    </div>
  );
}
