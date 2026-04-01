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
      className={`relative transition-colors cursor-pointer ${
        dragOver
          ? "border-spotify-green bg-spotify-green/5"
          : "hover:border-surface-700"
      }`}
    >
      <div
        className="flex flex-col items-center justify-center py-12 text-center"
        onDragOver={(e) => {
          e.preventDefault();
          setDragOver(true);
        }}
        onDragLeave={() => setDragOver(false)}
        onDrop={handleDrop}
        onClick={() => inputRef.current?.click()}
      >
        {loading ? (
          <>
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-surface-700 border-t-spotify-green mb-4" />
            <p className="text-sm text-surface-200">Parsing CSV files...</p>
          </>
        ) : (
          <>
            <svg
              className="h-12 w-12 text-surface-200 mb-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={1}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              />
            </svg>
            <p className="text-lg font-semibold mb-1">
              Drop your Exportify CSV files here
            </p>
            <p className="text-sm text-surface-200 mb-4">
              or click to browse
            </p>
            <p className="text-xs text-surface-700 max-w-sm">
              Export your Spotify playlists at{" "}
              <a
                href="https://exportify.app"
                target="_blank"
                rel="noopener noreferrer"
                className="text-spotify-green hover:underline"
                onClick={(e) => e.stopPropagation()}
              >
                exportify.app
              </a>
              {" "}and upload the CSV files here
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
