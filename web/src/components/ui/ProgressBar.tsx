interface Props {
  current: number;
  total: number;
  className?: string;
}

export function ProgressBar({ current, total, className = "" }: Props) {
  const pct = total > 0 ? Math.round((current / total) * 100) : 0;

  return (
    <div className={`w-full ${className}`}>
      <div className="flex justify-between text-xs text-surface-200 mb-1">
        <span>{current} / {total}</span>
        <span>{pct}%</span>
      </div>
      <div className="h-2 w-full rounded-full bg-surface-800 overflow-hidden">
        <div
          className="h-full rounded-full transition-all duration-300 bg-gradient-to-r from-spotify-green to-tidal-blue"
          style={{ width: `${pct}%` }}
        />
      </div>
    </div>
  );
}
