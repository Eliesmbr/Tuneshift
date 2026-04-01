import { useEffect, useState } from "react";

interface ToastItem {
  id: number;
  message: string;
  type: "error" | "success" | "info";
}

let toastId = 0;
let addToastFn: ((toast: Omit<ToastItem, "id">) => void) | null = null;

export function toast(message: string, type: "error" | "success" | "info" = "error") {
  addToastFn?.({ message, type });
}

export function ToastContainer() {
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  useEffect(() => {
    addToastFn = (t) => {
      const id = ++toastId;
      setToasts((prev) => [...prev, { ...t, id }]);
      setTimeout(() => {
        setToasts((prev) => prev.filter((x) => x.id !== id));
      }, 5000);
    };
    return () => { addToastFn = null; };
  }, []);

  if (toasts.length === 0) return null;

  return (
    <div className="fixed top-4 right-4 z-[100] flex flex-col gap-2 max-w-sm">
      {toasts.map((t) => (
        <div
          key={t.id}
          className={`rounded-xl px-4 py-3 text-sm font-medium shadow-lg backdrop-blur-sm animate-[slideIn_0.2s_ease-out] ${
            t.type === "error"
              ? "bg-red-500/90 text-white"
              : t.type === "success"
                ? "bg-spotify-green/90 text-black"
                : "bg-surface-800/90 text-white border border-surface-700"
          }`}
          onClick={() => setToasts((prev) => prev.filter((x) => x.id !== t.id))}
        >
          {t.message}
        </div>
      ))}
    </div>
  );
}
