import { useCallback, useRef, useState } from "react";
import { Card } from "../ui/Card";

interface Props {
  onUpload: (file: File) => void;
  loading: boolean;
}

export function TakeoutUpload({ onUpload, loading }: Props) {
  const [dragOver, setDragOver] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setDragOver(false);
      const file = Array.from(e.dataTransfer.files).find((f) =>
        f.name.endsWith(".zip"),
      );
      if (file) onUpload(file);
    },
    [onUpload],
  );

  const handleFileSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (file) onUpload(file);
    },
    [onUpload],
  );

  return (
    <Card
      className={`relative transition-all duration-200 cursor-pointer ${
        dragOver
          ? "border-red-400 bg-red-400/5 scale-[1.02]"
          : "hover:border-surface-700 hover:bg-surface-900/80"
      }`}
    >
      <div
        className="flex flex-col items-center justify-center py-10 text-center"
        onDragOver={(e) => {
          e.preventDefault();
          setDragOver(true);
        }}
        onDragLeave={() => setDragOver(false)}
        onDrop={handleDrop}
        onClick={() => inputRef.current?.click()}
      >
        {loading ? (
          <div className="flex flex-col items-center gap-3">
            <div className="h-10 w-10 animate-spin rounded-full border-2 border-surface-700 border-t-red-400" />
            <p className="text-sm text-surface-200">Parsing your playlists...</p>
          </div>
        ) : (
          <>
            <div className="mb-4 h-14 w-14 rounded-2xl bg-gradient-to-br from-red-500/20 to-tidal-blue/20 flex items-center justify-center">
              <svg
                className="h-7 w-7 text-surface-200"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                strokeWidth={1.5}
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5"
                />
              </svg>
            </div>
            <p className="text-lg font-semibold mb-1">
              Drop your Takeout ZIP here
            </p>
            <p className="text-sm text-surface-200">
              or click to browse
            </p>
          </>
        )}
      </div>
      <input
        ref={inputRef}
        type="file"
        accept=".zip"
        className="hidden"
        onChange={handleFileSelect}
      />
    </Card>
  );
}
