type SectionHeadingProps = {
  eyebrow: string;
  title: string;
  description?: string;
};

export function SectionHeading({ eyebrow, title, description }: SectionHeadingProps) {
  return (
    <div className="max-w-2xl space-y-3">
      <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[var(--accent)]">{eyebrow}</p>
      <h2 className="text-3xl font-semibold tracking-tight text-balance md:text-4xl">{title}</h2>
      {description ? <p className="text-base leading-7 text-[var(--text-secondary)]">{description}</p> : null}
    </div>
  );
}