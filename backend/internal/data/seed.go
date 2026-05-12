package data

import "devfolio/backend/internal/domain"

func Seed() domain.Data {
	return domain.Data{
		Site: domain.Site{
			Name:        "DevFolio",
			Tagline:     "Developer portfolio platform",
			Description: "Portfolio digitale personalizzabile per sviluppatori, con base modulare e pronta per il white-label.",
			Accent:      "#1a56db",
		},
		Projects: []domain.Project{
			{
				Title:            "DevFolio Core",
				Slug:             "devfolio-core",
				DescriptionShort: "Base white-label per portfolio developer con componenti riusabili e tema brandizzabile.",
				DescriptionLong:  "Una base di partenza per costruire portfolio professionali, con layout modulare, contenuti strutturati e un design system pensato per essere rebrandizzato rapidamente.",
				Technologies:     []string{"Next.js", "TypeScript", "Tailwind CSS"},
				Status:           "published",
				Featured:         true,
				Year:             "2026",
				Image:            "https://images.unsplash.com/photo-1517694712202-14dd9538aa97?auto=format&fit=crop&w=1400&q=80",
				LiveURL:          "https://example.com",
				GitHubURL:        "https://github.com/gigliofr",
			},
			{
				Title:            "Content Hub",
				Slug:             "content-hub",
				DescriptionShort: "Strato editoriale per articoli tecnici, contenuti SEO e gestione delle pubblicazioni.",
				DescriptionLong:  "Un progetto dedicato alla pubblicazione di articoli tecnici e aggiornamenti, con paginazione, tagging e una struttura pronta per l'evoluzione verso un CMS completo.",
				Technologies:     []string{"MDX", "SEO", "Design Systems"},
				Status:           "published",
				Featured:         true,
				Year:             "2026",
				Image:            "https://images.unsplash.com/photo-1498050108023-c5249f4df085?auto=format&fit=crop&w=1400&q=80",
				GitHubURL:        "https://github.com/gigliofr",
			},
			{
				Title:            "Admin Console",
				Slug:             "admin-console",
				DescriptionShort: "Interfaccia di amministrazione con CRUD, configurazioni e flussi di editing.",
				DescriptionLong:  "La console admin sarà il centro operativo del progetto: gestione di progetti, articoli, asset media e impostazioni del brand senza passare dal codice.",
				Technologies:     []string{"Payload CMS", "Auth", "CRUD"},
				Status:           "draft",
				Featured:         false,
				Year:             "2026",
				Image:            "https://images.unsplash.com/photo-1555066931-4365d14bab8c?auto=format&fit=crop&w=1400&q=80",
			},
		},
		Posts: []domain.Post{
			{
				Title:       "Come strutturare un portfolio developer che cresce con te",
				Slug:        "portfolio-che-cresce",
				Excerpt:     "Una base modulare evita di rifare il sito ogni volta che cambiano i contenuti o il brand.",
				Category:    "Architecture",
				ReadTime:    "4 min",
				PublishedAt: "10 Apr 2026",
				Tags:        []string{"Next.js", "Portfolio", "CMS"},
				Sections: []domain.PostSection{
					{
						Heading: "Partire da un modello di contenuto semplice",
						Paragraphs: []string{
							"Un portfolio cresce male quando il contenuto è sparso dentro componenti hardcoded. Il primo passo è separare il dato dalla presentazione e mantenere un set minimo di campi coerenti per progetti, articoli e contatti.",
							"Nel caso di DevFolio, il modello iniziale tiene insieme titolo, slug, estratto, data e tag. Da qui si può evolvere verso sezioni più ricche senza cambiare il modo in cui le pagine leggono i contenuti.",
						},
						Bullets: []string{
							"Contenuto normalizzato prima del layout",
							"Slug stabili per routing e metadata",
							"Campi opzionali per l'evoluzione editoriale",
						},
					},
					{
						Heading: "Evitare la riscrittura del sito",
						Paragraphs: []string{
							"Se il portfolio è costruito come una collezione di blocchi riusabili, cambiare brand, copy o struttura editoriale non richiede di rifare tutto il sito. Il lavoro diventa una sequenza di miglioramenti localizzati.",
						},
					},
				},
			},
			{
				Title:       "Design system minimo, identità forte",
				Slug:        "design-system-minimo",
				Excerpt:     "Pochi token ben definiti bastano per far sembrare il prodotto coerente e professionale.",
				Category:    "Design",
				ReadTime:    "5 min",
				PublishedAt: "6 Apr 2026",
				Tags:        []string{"Tailwind", "Brand", "Tokens"},
				Sections: []domain.PostSection{
					{
						Heading: "I token prima delle varianti",
						Paragraphs: []string{
							"Un design system utile parte da pochi token chiari: superfici, testo, accento, bordi e spaziatura. Questo evita il classico effetto patchwork che compare quando ogni sezione definisce i propri colori.",
						},
						Bullets: []string{
							"Colori derivati da variabili CSS",
							"Tipografia coerente in tutto il sito",
							"Componenti riusabili con una sola fonte di verità",
						},
					},
				},
			},
			{
				Title:       "Phase-driven delivery per progetti personali",
				Slug:        "phase-driven-delivery",
				Excerpt:     "Spezzare il lavoro in fasi rende il progetto più verificabile e più facile da mantenere.",
				Category:    "Process",
				ReadTime:    "3 min",
				PublishedAt: "2 Apr 2026",
				Tags:        []string{"Planning", "Iteration", "Quality"},
				Sections: []domain.PostSection{
					{
						Heading: "Rilasciare per incrementi",
						Paragraphs: []string{
							"Per un progetto personale la disciplina vale più della quantità di feature. Il lavoro per fasi crea checkpoint chiari: base visiva, dati, editoriale, SEO, hardening.",
							"Ogni fase deve produrre qualcosa di utilizzabile e verificabile, così il progetto resta vivo anche quando il tempo disponibile è frammentato.",
						},
					},
				},
			},
		},
		Stats: []domain.Stat{
			{Label: "Reusable modules", Value: "12+"},
			{Label: "Target Lighthouse", Value: "90+"},
			{Label: "Brand-ready layers", Value: "3"},
			{Label: "Content types", Value: "4"},
		},
		Skills: []string{"Next.js", "TypeScript", "React", "Tailwind CSS", "Payload CMS", "PostgreSQL", "GitHub API", "SEO"},
	}
}