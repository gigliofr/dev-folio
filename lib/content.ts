export type Project = {
  title: string;
  slug: string;
  descriptionShort: string;
  descriptionLong: string;
  technologies: string[];
  status: 'draft' | 'published' | 'archived';
  featured?: boolean;
  year: string;
  image: string;
  liveUrl?: string;
  githubUrl?: string;
};

export type Post = {
  title: string;
  slug: string;
  excerpt: string;
  category: string;
  readTime: string;
  publishedAt: string;
  tags: string[];
  sections?: PostSection[];
};

export type PostSection = {
  heading: string;
  paragraphs: string[];
  bullets?: string[];
};

export const projects: Project[] = [
  {
    title: 'DevFolio Core',
    slug: 'devfolio-core',
    descriptionShort: 'Base white-label per portfolio developer con componenti riusabili e tema brandizzabile.',
    descriptionLong:
      'Una base di partenza per costruire portfolio professionali, con layout modulare, contenuti strutturati e un design system pensato per essere rebrandizzato rapidamente.',
    technologies: ['Next.js', 'TypeScript', 'Tailwind CSS'],
    status: 'published',
    featured: true,
    year: '2026',
    image: 'https://images.unsplash.com/photo-1517694712202-14dd9538aa97?auto=format&fit=crop&w=1400&q=80',
    liveUrl: 'https://example.com',
    githubUrl: 'https://github.com/gigliofr'
  },
  {
    title: 'Content Hub',
    slug: 'content-hub',
    descriptionShort: 'Strato editoriale per articoli tecnici, contenuti SEO e gestione delle pubblicazioni.',
    descriptionLong:
      'Un progetto dedicato alla pubblicazione di articoli tecnici e aggiornamenti, con paginazione, tagging e una struttura pronta per l evoluzione verso un CMS completo.',
    technologies: ['MDX', 'SEO', 'Design Systems'],
    status: 'published',
    featured: true,
    year: '2026',
    image: 'https://images.unsplash.com/photo-1498050108023-c5249f4df085?auto=format&fit=crop&w=1400&q=80',
    githubUrl: 'https://github.com/gigliofr'
  },
  {
    title: 'Admin Console',
    slug: 'admin-console',
    descriptionShort: 'Interfaccia di amministrazione con CRUD, configurazioni e flussi di editing.',
    descriptionLong:
      'La console admin sarà il centro operativo del progetto: gestione di progetti, articoli, asset media e impostazioni del brand senza passare dal codice.',
    technologies: ['Payload CMS', 'Auth', 'CRUD'],
    status: 'draft',
    featured: false,
    year: '2026',
    image: 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?auto=format&fit=crop&w=1400&q=80'
  }
];

export const posts: Post[] = [
  {
    title: 'Come strutturare un portfolio developer che cresce con te',
    slug: 'portfolio-che-cresce',
    excerpt: 'Una base modulare evita di rifare il sito ogni volta che cambiano i contenuti o il brand.',
    category: 'Architecture',
    readTime: '4 min',
    publishedAt: '10 Apr 2026',
    tags: ['Next.js', 'Portfolio', 'CMS'],
    sections: [
      {
        heading: 'Partire da un modello di contenuto semplice',
        paragraphs: [
          'Un portfolio cresce male quando il contenuto è sparso dentro componenti hardcoded. Il primo passo è separare il dato dalla presentazione e mantenere un set minimo di campi coerenti per progetti, articoli e contatti.',
          'Nel caso di DevFolio, il modello iniziale tiene insieme titolo, slug, estratto, data e tag. Da qui si può evolvere verso sezioni più ricche senza cambiare il modo in cui le pagine leggono i contenuti.'
        ],
        bullets: [
          'Contenuto normalizzato prima del layout',
          'Slug stabili per routing e metadata',
          'Campi opzionali per l evoluzione editoriale'
        ]
      },
      {
        heading: 'Evitare la riscrittura del sito',
        paragraphs: [
          'Se il portfolio è costruito come una collezione di blocchi riusabili, cambiare brand, copy o struttura editoriale non richiede di rifare tutto il sito. Il lavoro diventa una sequenza di miglioramenti localizzati.'
        ]
      }
    ]
  },
  {
    title: 'Design system minimo, identità forte',
    slug: 'design-system-minimo',
    excerpt: 'Pochi token ben definiti bastano per far sembrare il prodotto coerente e professionale.',
    category: 'Design',
    readTime: '5 min',
    publishedAt: '6 Apr 2026',
    tags: ['Tailwind', 'Brand', 'Tokens'],
    sections: [
      {
        heading: 'I token prima delle varianti',
        paragraphs: [
          'Un design system utile parte da pochi token chiari: superfici, testo, accento, bordi e spaziatura. Questo evita il classico effetto patchwork che compare quando ogni sezione definisce i propri colori.'
        ],
        bullets: [
          'Colori derivati da variabili CSS',
          'Tipografia coerente in tutto il sito',
          'Componenti riusabili con una sola fonte di verità'
        ]
      }
    ]
  },
  {
    title: 'Phase-driven delivery per progetti personali',
    slug: 'phase-driven-delivery',
    excerpt: 'Spezzare il lavoro in fasi rende il progetto più verificabile e più facile da mantenere.',
    category: 'Process',
    readTime: '3 min',
    publishedAt: '2 Apr 2026',
    tags: ['Planning', 'Iteration', 'Quality'],
    sections: [
      {
        heading: 'Rilasciare per incrementi',
        paragraphs: [
          'Per un progetto personale la disciplina vale più della quantità di feature. Il lavoro per fasi crea checkpoint chiari: base visiva, dati, editoriale, SEO, hardening.',
          'Ogni fase deve produrre qualcosa di utilizzabile e verificabile, così il progetto resta vivo anche quando il tempo disponibile è frammentato.'
        ]
      }
    ]
  }
];

export const stats = [
  { label: 'Reusable modules', value: '12+' },
  { label: 'Target Lighthouse', value: '90+' },
  { label: 'Brand-ready layers', value: '3' },
  { label: 'Content types', value: '4' }
];

export const skills = [
  'Next.js',
  'TypeScript',
  'React',
  'Tailwind CSS',
  'Payload CMS',
  'PostgreSQL',
  'GitHub API',
  'SEO'
];