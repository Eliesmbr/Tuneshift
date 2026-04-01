export function Header() {
  return (
    <header className="border-b border-surface-800 bg-surface-950/80 backdrop-blur-md sticky top-0 z-50">
      <div className="mx-auto max-w-5xl px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-spotify-green to-tidal-blue" />
          <h1 className="text-xl font-bold tracking-tight">Tuneshift</h1>
        </div>
        <p className="text-xs text-surface-200">Free & Open Source</p>
      </div>
    </header>
  );
}
