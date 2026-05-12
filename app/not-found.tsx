export default function NotFoundPage() {
  return (
    <section className="container-shell py-20 md:py-28">
      <div className="max-w-xl space-y-4">
        <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[var(--accent)]">404</p>
        <h1 className="text-4xl font-semibold tracking-tight">This page does not exist yet.</h1>
        <p className="text-base leading-7 text-[var(--text-secondary)]">Return to the homepage or continue with the portfolio routes already in place.</p>
      </div>
    </section>
  );
}