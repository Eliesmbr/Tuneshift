import { useCallback, useEffect, useState } from "react";
import { Header } from "./components/layout/Header";
import { StepIndicator } from "./components/layout/StepIndicator";
import { HeroAnimation } from "./components/landing/HeroAnimation";
import { SourceSelector } from "./components/landing/SourceSelector";
import { CSVUpload } from "./components/upload/CSVUpload";
import { PlaylistPreview } from "./components/upload/PlaylistPreview";
import { SelectionPanel } from "./components/library/SelectionPanel";
import { ConnectButton } from "./components/auth/ConnectButton";
import { MigrationProgressView } from "./components/migration/MigrationProgress";
import { MigrationSummary } from "./components/migration/MigrationSummary";
import { TakeoutUpload } from "./components/upload/TakeoutUpload";
import { ToastContainer, toast } from "./components/ui/Toast";
import { Card } from "./components/ui/Card";
import { Button } from "./components/ui/Button";
import { useAuth } from "./hooks/useAuth";
import { useMigration } from "./hooks/useMigration";
import { api } from "./api/client";
import type { Source, Step, UploadedPlaylist } from "./types";

const STORAGE_KEY = "tuneshift_upload";

interface PersistedState {
  uploadSessionId: string;
  playlists: UploadedPlaylist[];
  totalTracks: number;
  selectedPlaylists: string[];
  sourceSelected: boolean;
  selectedSource: Source | null;
}

function saveState(state: PersistedState) {
  sessionStorage.setItem(STORAGE_KEY, JSON.stringify(state));
}

function loadState(): PersistedState | null {
  const raw = sessionStorage.getItem(STORAGE_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

function clearState() {
  sessionStorage.removeItem(STORAGE_KEY);
}

export default function App() {
  const auth = useAuth();
  const migration = useMigration();
  const restored = loadState();

  const [step, setStep] = useState<Step>("upload");
  const [selectedSource, setSelectedSource] = useState<Source | null>(restored?.selectedSource ?? "spotify");
  const [sourceSelected, setSourceSelected] = useState(restored?.sourceSelected ?? false);
  const [uploading, setUploading] = useState(false);
  const [uploadSessionId, setUploadSessionId] = useState<string | null>(restored?.uploadSessionId ?? null);
  const [playlists, setPlaylists] = useState<UploadedPlaylist[]>(restored?.playlists ?? []);
  const [totalTracks, setTotalTracks] = useState(restored?.totalTracks ?? 0);
  const [selectedPlaylists, setSelectedPlaylists] = useState<string[]>(restored?.selectedPlaylists ?? []);

  useEffect(() => {
    if (auth.loading) return;

    if (migration.done) {
      setStep("done");
    } else if (migration.running) {
      setStep("migrate");
    } else if (selectedPlaylists.length > 0) {
      setStep("connect-tidal");
    } else if (playlists.length > 0) {
      setStep("select");
    } else {
      setStep("upload");
    }
  }, [
    auth.loading,
    auth.tidal.connected,
    playlists.length,
    selectedPlaylists.length,
    migration.running,
    migration.done,
  ]);

  const handleSelectSource = useCallback((source: Source) => {
    setSelectedSource(source);
  }, []);

  const handleConfirmSource = useCallback(() => {
    if (!selectedSource) return;
    setSourceSelected(true);
  }, [selectedSource]);

  const handleUpload = useCallback(async (files: File[]) => {
    setUploading(true);
    try {
      const result = await api.uploadCSV(files);
      setUploadSessionId(result.session_id);
      setPlaylists(result.playlists);
      setTotalTracks(result.total_tracks);
      toast(`${result.playlists.length} playlist${result.playlists.length !== 1 ? "s" : ""} loaded`, "success");
      saveState({
        uploadSessionId: result.session_id,
        playlists: result.playlists,
        totalTracks: result.total_tracks,
        selectedPlaylists: [],
        sourceSelected: true,
        selectedSource: "spotify",
      });
    } catch (err) {
      toast(err instanceof Error ? err.message : "Upload failed", "error");
    } finally {
      setUploading(false);
    }
  }, []);

  const handleTakeoutUpload = useCallback(async (file: File) => {
    setUploading(true);
    try {
      const result = await api.uploadTakeout(file);
      setUploadSessionId(result.session_id);
      setPlaylists(result.playlists);
      setTotalTracks(result.total_tracks);
      toast(`${result.playlists.length} playlist${result.playlists.length !== 1 ? "s" : ""} loaded`, "success");
      saveState({
        uploadSessionId: result.session_id,
        playlists: result.playlists,
        totalTracks: result.total_tracks,
        selectedPlaylists: [],
        sourceSelected: true,
        selectedSource: "youtube-music",
      });
    } catch (err) {
      toast(err instanceof Error ? err.message : "Upload failed", "error");
    } finally {
      setUploading(false);
    }
  }, []);

  const handleSelection = useCallback(
    (names: string[]) => {
      setSelectedPlaylists(names);
      if (uploadSessionId) {
        saveState({
          uploadSessionId,
          playlists,
          totalTracks,
          selectedPlaylists: names,
          sourceSelected: true,
          selectedSource: selectedSource,
        });
      }
    },
    [uploadSessionId, playlists, totalTracks, selectedSource],
  );

  const handleStartMigration = useCallback(() => {
    if (uploadSessionId) {
      migration.start({
        upload_session_id: uploadSessionId,
        playlists: selectedPlaylists,
      });
    }
  }, [uploadSessionId, selectedPlaylists, migration]);

  const handleStartOver = useCallback(() => {
    clearState();
    window.location.reload();
  }, []);

  if (auth.loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-surface-700 border-t-white" />
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <Header onLogoClick={handleStartOver} />
      <ToastContainer />
      <main className="mx-auto max-w-2xl px-6 py-8">
        <StepIndicator current={step} source={selectedSource} />

        {/* Source selection / Upload */}
        {step === "upload" && (
          <div className="space-y-8">
            {!sourceSelected ? (
              <>
                <div className="text-center animate-[fadeIn_0.5s_ease-out]">
                  <h2 className="text-4xl font-bold mb-3 tracking-tight">
                    Move your music to{" "}
                    <span className="bg-gradient-to-r from-spotify-green to-tidal-blue bg-clip-text text-transparent">
                      Tidal
                    </span>
                  </h2>
                  <p className="text-surface-200 max-w-md mx-auto">
                    Transfer your playlists, keep your music. Free and open source.
                  </p>
                </div>

                <HeroAnimation source={selectedSource} />
                <SourceSelector selected={selectedSource} onSelectSource={handleSelectSource} />

                {selectedSource === "youtube-music" && (
                  <div className="rounded-xl border border-surface-700/50 bg-surface-900/50 px-4 py-3 text-sm text-surface-200 text-center animate-[fadeIn_0.3s_ease-out]">
                    YouTube's API requires a costly security audit for public apps, so I use{" "}
                    <a href="https://takeout.google.com" target="_blank" rel="noopener noreferrer" className="text-red-400 hover:underline">Google Takeout</a>{" "}
                    instead — your data, exported directly by Google.
                  </div>
                )}

                {selectedSource && (
                  <div className="flex justify-center animate-[fadeIn_0.3s_ease-out]">
                    <button
                      onClick={handleConfirmSource}
                      className={`px-8 py-3 text-base font-semibold rounded-xl transition-all hover:opacity-90 cursor-pointer ${
                        selectedSource === "spotify"
                          ? "bg-gradient-to-r from-spotify-green to-tidal-blue text-black"
                          : "bg-gradient-to-r from-red-500 to-tidal-blue text-white"
                      }`}
                    >
                      Start migration from {selectedSource === "spotify" ? "Spotify" : "YouTube Music"}
                    </button>
                  </div>
                )}
              </>
            ) : selectedSource === "youtube-music" ? (
              <>
                <div className="text-center animate-[fadeIn_0.3s_ease-out]">
                  <h2 className="text-2xl font-bold mb-2">Upload your YouTube Music data</h2>
                  <p className="text-surface-200 text-sm">
                    Export your data at{" "}
                    <a
                      href="https://takeout.google.com"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-red-400 hover:underline"
                    >
                      takeout.google.com
                    </a>{" "}
                    and upload the ZIP file below.
                  </p>
                </div>

                <TakeoutUpload onUpload={handleTakeoutUpload} loading={uploading} />

                <Card className="p-4 animate-[slideUp_0.4s_ease-out_0.2s_both]">
                  <h3 className="text-sm font-semibold mb-3">How it works</h3>
                  <ol className="space-y-2 text-sm text-surface-200">
                    <li className="flex gap-3">
                      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-800 text-xs font-bold text-white">1</span>
                      <span>Go to <a href="https://takeout.google.com" target="_blank" rel="noopener noreferrer" className="text-red-400 hover:underline">takeout.google.com</a></span>
                    </li>
                    <li className="flex gap-3">
                      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-800 text-xs font-bold text-white">2</span>
                      <span>Click "Deselect all", then select only "YouTube and YouTube Music"</span>
                    </li>
                    <li className="flex gap-3">
                      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-800 text-xs font-bold text-white">3</span>
                      <span>Click "All YouTube data included" and select <strong>Music library</strong> and <strong>Playlists</strong></span>
                    </li>
                    <li className="flex gap-3">
                      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-800 text-xs font-bold text-white">4</span>
                      <span>Create the export, wait for the email, then upload the ZIP here</span>
                    </li>
                  </ol>
                </Card>
              </>
            ) : (
              <>
                <div className="text-center animate-[fadeIn_0.3s_ease-out]">
                  <h2 className="text-2xl font-bold mb-2">Upload your Spotify playlists</h2>
                  <p className="text-surface-200 text-sm">
                    Export your playlists at{" "}
                    <a
                      href="https://exportify.app"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-spotify-green hover:underline"
                    >
                      exportify.app
                    </a>{" "}
                    and drop the CSV files below.
                  </p>
                </div>

                <CSVUpload onUpload={handleUpload} loading={uploading} />

                <Card className="p-4 animate-[slideUp_0.4s_ease-out_0.2s_both]">
                  <h3 className="text-sm font-semibold mb-3">How it works</h3>
                  <ol className="space-y-2 text-sm text-surface-200">
                    <li className="flex gap-3">
                      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-800 text-xs font-bold text-white">1</span>
                      <span>Go to <a href="https://exportify.app" target="_blank" rel="noopener noreferrer" className="text-spotify-green hover:underline">exportify.app</a> and log in with Spotify</span>
                    </li>
                    <li className="flex gap-3">
                      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-800 text-xs font-bold text-white">2</span>
                      <span>Click "Export" on each playlist you want to transfer</span>
                    </li>
                    <li className="flex gap-3">
                      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-800 text-xs font-bold text-white">3</span>
                      <span>Upload the downloaded CSV files here</span>
                    </li>
                    <li className="flex gap-3">
                      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-800 text-xs font-bold text-white">4</span>
                      <span>Connect your Tidal account and start the migration</span>
                    </li>
                  </ol>
                </Card>
              </>
            )}
          </div>
        )}

        {/* Select playlists */}
        {step === "select" && (
          <div className="space-y-6">
            <PlaylistPreview playlists={playlists} totalTracks={totalTracks} />
            <h3 className="text-lg font-semibold">Select playlists to migrate</h3>
            <SelectionPanel
              playlists={playlists}
              onSelect={handleSelection}
            />
          </div>
        )}

        {/* Connect Tidal */}
        {step === "connect-tidal" && (
          <div className="space-y-6">
            <Card className="p-4">
              <p className="text-sm text-surface-200">
                {selectedPlaylists.length} playlist{selectedPlaylists.length !== 1 ? "s" : ""} selected
              </p>
            </Card>

            <ConnectButton
              service="tidal"
              connected={auth.tidal.connected}
              userName={auth.tidal.user?.name}
              onConnect={auth.connectTidal}
              onDisconnect={auth.disconnectTidal}
            />

            {auth.tidal.connected && (
              <Card>
                <h3 className="font-semibold mb-3">Ready to migrate</h3>
                <ul className="space-y-1 text-sm text-surface-200 mb-4">
                  {selectedPlaylists.map((name) => (
                    <li key={name}>+ {name}</li>
                  ))}
                </ul>
                <Button
                  variant="primary"
                  onClick={handleStartMigration}
                  className="w-full bg-gradient-to-r from-spotify-green to-tidal-blue text-black"
                >
                  Start Migration
                </Button>
              </Card>
            )}
          </div>
        )}

        {/* Migration in progress */}
        {step === "migrate" && (
          <MigrationProgressView
            events={migration.events}
            running={migration.running}
            error={migration.error}
          />
        )}

        {/* Done */}
        {step === "done" && (
          <MigrationSummary
            events={migration.events}
            onStartOver={handleStartOver}
          />
        )}

        <footer className="mt-16 pb-8 text-center text-xs text-surface-700 space-y-2">
          <p>
            Tuneshift is free and open source. Your data is parsed
            server-side and never stored permanently.
          </p>
          <p>
            <a href="/privacy" className="hover:text-surface-200 underline">Privacy Policy</a>
          </p>
        </footer>
      </main>
    </div>
  );
}
