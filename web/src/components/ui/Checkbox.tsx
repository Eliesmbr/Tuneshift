interface Props {
  checked: boolean;
  onChange: (checked: boolean) => void;
  label: string;
  description?: string;
}

export function Checkbox({ checked, onChange, label, description }: Props) {
  return (
    <label className="flex items-start gap-3 cursor-pointer group">
      <div className="relative mt-0.5">
        <input
          type="checkbox"
          checked={checked}
          onChange={(e) => onChange(e.target.checked)}
          className="peer sr-only"
        />
        <div className="h-5 w-5 rounded-md border-2 border-surface-700 bg-surface-800 transition-colors peer-checked:border-spotify-green peer-checked:bg-spotify-green group-hover:border-surface-200">
          {checked && (
            <svg className="h-5 w-5 text-black" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
            </svg>
          )}
        </div>
      </div>
      <div>
        <span className="text-sm font-medium text-white">{label}</span>
        {description && <p className="text-xs text-surface-200 mt-0.5">{description}</p>}
      </div>
    </label>
  );
}
