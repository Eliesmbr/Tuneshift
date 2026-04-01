import { useCallback, useRef, useState } from "react";
import { Card } from "../ui/Card";

interface Props {
  onUpload: (files: File[]) => void;
  loading: boolean;
}

export function CSVUpload({ onUpload, loading }: Props) {
  const [dragOver, setDragOver] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setDragOver(false);
      const files = Array.from(e.dataTransfer.files).filter((f) =>
        f.name.endsWith(".csv"),
      );
      if (files.length > 0) onUpload(files);
    },
    [onUpload],
  );

  const handleFileSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const files = Array.from(e.target.files ?? []);
      if (files.length > 0) onUpload(files);
    },
    [onUpload],
  );

  return (
    <Card
      className={`relative transition-all duration-200 cursor-pointer ${
        dragOver
          ? "border-spotify-green bg-spotify-green/5 scale-[1.02]"
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
            <div className="h-10 w-10 animate-spin rounded-full border-2 border-surface-700 border-t-spotify-green" />
            <p className="text-sm text-surface-200">Parsing your playlists...</p>
          </div>
        ) : (
          <>
            <div className="mb-4 h-14 w-14 rounded-2xl bg-gradient-to-br from-spotify-green/20 to-tidal-blue/20 flex items-center justify-center">
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
              Drop your CSV files here
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
        accept=".csv"
        multiple
        className="hidden"
        onChange={handleFileSelect}
      />
    </Card>
  );
}
