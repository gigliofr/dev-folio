import type { Metadata } from 'next';
import { Inter, JetBrains_Mono } from 'next/font/google';
import './globals.css';
import { ThemeProvider } from '@/components/theme-provider';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { site } from '@/lib/site';

const inter = Inter({ subsets: ['latin'], variable: '--font-sans' });
const jetBrainsMono = JetBrains_Mono({ subsets: ['latin'], variable: '--font-mono' });

export const metadata: Metadata = {
  metadataBase: new URL(site.url),
  title: {
    default: site.name,
    template: `%s | ${site.name}`
  },
  description: site.description,
  openGraph: {
    title: site.name,
    description: site.description,
    url: site.url,
    siteName: site.name,
    type: 'website'
  },
  twitter: {
    card: 'summary_large_image',
    title: site.name,
    description: site.description
  }
};

export default function RootLayout({ children }: Readonly<{ children: import('react').ReactNode }>) {
  return (
    <html lang="it" suppressHydrationWarning>
      <body className={`${inter.variable} ${jetBrainsMono.variable} font-sans antialiased`}>
        <script
          dangerouslySetInnerHTML={{
            __html: `(() => { try { const theme = localStorage.getItem('devfolio-theme'); const preferred = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'; document.documentElement.dataset.theme = theme === 'dark' || theme === 'light' ? theme : preferred; } catch (error) {} })();`
          }}
        />
        <ThemeProvider>
          <div className="relative isolate min-h-screen">
            <div className="absolute inset-x-0 top-0 -z-10 h-[32rem] bg-[radial-gradient(circle_at_top,_rgba(26,86,219,0.16),_transparent_46%)]" />
            <SiteHeader />
            <main>{children}</main>
            <SiteFooter />
          </div>
        </ThemeProvider>
      </body>
    </html>
  );
}