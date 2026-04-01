import type { ReactNode } from "react";

interface Props {
  children: ReactNode;
  className?: string;
}

export function Card({ children, className = "" }: Props) {
  return (
    <div className={`rounded-2xl border border-surface-800 bg-surface-900/50 p-6 backdrop-blur-sm ${className}`}>
      {children}
    </div>
  );
}
