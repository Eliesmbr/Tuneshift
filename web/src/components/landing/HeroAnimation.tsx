export function HeroAnimation() {
  return (
    <div className="flex items-center justify-center gap-6 py-8">
      {/* Spotify */}
      <div className="flex flex-col items-center gap-2 animate-[fadeInLeft_0.6s_ease-out]">
        <div className="h-16 w-16 rounded-2xl bg-spotify-green/10 flex items-center justify-center">
          <svg className="h-10 w-10 text-spotify-green" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z" />
          </svg>
        </div>
        <span className="text-xs text-surface-200 font-medium">Spotify</span>
      </div>

      {/* Animated arrow */}
      <div className="flex items-center gap-1 animate-[fadeIn_0.8s_ease-out_0.3s_both]">
        <div className="h-px w-8 bg-gradient-to-r from-spotify-green to-transparent" />
        <div className="relative">
          <div className="h-10 w-10 rounded-full bg-gradient-to-r from-spotify-green to-tidal-blue flex items-center justify-center animate-pulse">
            <svg className="h-5 w-5 text-black" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
          </div>
        </div>
        <div className="h-px w-8 bg-gradient-to-r from-transparent to-tidal-blue" />
      </div>

      {/* Tidal */}
      <div className="flex flex-col items-center gap-2 animate-[fadeInRight_0.6s_ease-out_0.2s_both]">
        <div className="h-16 w-16 rounded-2xl bg-tidal-blue/10 flex items-center justify-center">
          <svg className="h-10 w-10 text-tidal-blue" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12.012 3.992L8.008 7.996 12.012 12l4.004-4.004zM4.004 7.996L0 12l4.004 4.004L8.008 12zm15.992 0L15.992 12l4.004 4.004L24 12zM12.012 12l-4.004 4.004 4.004 4.004 4.004-4.004z" />
          </svg>
        </div>
        <span className="text-xs text-surface-200 font-medium">Tidal</span>
      </div>
    </div>
  );
}
