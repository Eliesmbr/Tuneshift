import type { Source, Step } from "../../types";

interface StepDef {
  key: Step;
  label: string;
}

const spotifySteps: StepDef[] = [
  { key: "upload", label: "Upload CSV" },
  { key: "select", label: "Select" },
  { key: "connect-tidal", label: "Connect Tidal" },
  { key: "migrate", label: "Migrate" },
  { key: "done", label: "Done" },
];

const youtubeSteps: StepDef[] = [
  { key: "connect-source", label: "Connect YouTube" },
  { key: "fetch-playlists", label: "Select" },
  { key: "connect-tidal", label: "Connect Tidal" },
  { key: "migrate", label: "Migrate" },
  { key: "done", label: "Done" },
];

interface Props {
  current: Step;
  source?: Source | null;
}

export function StepIndicator({ current, source }: Props) {
  const steps = source === "youtube-music" ? youtubeSteps : spotifySteps;
  const currentIdx = steps.findIndex((s) => s.key === current);

  return (
    <div className="flex items-center justify-center gap-2 py-6">
      {steps.map((step, i) => {
        const isActive = i === currentIdx;
        const isDone = i < currentIdx;

        return (
          <div key={step.key} className="flex items-center gap-2">
            <div className="flex items-center gap-2">
              <div
                className={`flex h-8 w-8 items-center justify-center rounded-full text-xs font-bold transition-colors ${
                  isDone
                    ? "bg-spotify-green text-black"
                    : isActive
                      ? "bg-white text-surface-950"
                      : "bg-surface-800 text-surface-200"
                }`}
              >
                {isDone ? (
                  <svg className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                  </svg>
                ) : (
                  i + 1
                )}
              </div>
              <span
                className={`hidden text-xs font-medium sm:block ${
                  isActive ? "text-white" : isDone ? "text-spotify-green" : "text-surface-200"
                }`}
              >
                {step.label}
              </span>
            </div>
            {i < steps.length - 1 && (
              <div
                className={`h-px w-8 ${
                  isDone ? "bg-spotify-green" : "bg-surface-800"
                }`}
              />
            )}
          </div>
        );
      })}
    </div>
  );
}
