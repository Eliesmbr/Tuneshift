import type { ButtonHTMLAttributes, ReactNode } from "react";

type Variant = "spotify" | "tidal" | "youtube" | "primary" | "secondary";

const variants: Record<Variant, string> = {
  spotify: "bg-spotify-green hover:bg-spotify-green/90 text-black font-semibold",
  tidal: "bg-tidal-blue hover:bg-tidal-blue/90 text-black font-semibold",
  youtube: "bg-red-500 hover:bg-red-500/90 text-white font-semibold",
  primary: "bg-white hover:bg-white/90 text-surface-950 font-semibold",
  secondary: "bg-surface-800 hover:bg-surface-700 text-white border border-surface-700",
};

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  children: ReactNode;
  icon?: ReactNode;
}

export function Button({ variant = "primary", children, icon, className = "", ...props }: Props) {
  return (
    <button
      className={`inline-flex items-center justify-center gap-2 rounded-xl px-6 py-3 text-sm transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed ${variants[variant]} ${className}`}
      {...props}
    >
      {icon}
      {children}
    </button>
  );
}
