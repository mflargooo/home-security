import { useEffect, useRef, useState, useCallback } from 'react';

export type ClipState = 'loading' | 'ready' | 'error' | 'unsupported';

interface UseHLSClipOptions {
  /** Seek to this position (seconds) when the clip loads. Defaults to 0. */
  startPosition?: number;
  onReady?: (duration: number) => void;
  onError?: (message: string) => void;
}

/**
 * VOD/clip HLS loader.
 * Optimised for scrubbing: large buffer, full seek range, user-initiated playback.
 * Does NOT autoplay — the component drives play/pause.
 */
export function useHLSClip(
  videoRef: React.RefObject<HTMLVideoElement>,
  url: string | null,
  options: UseHLSClipOptions = {}
) {
  const [clipState, setClipState] = useState<ClipState>('loading');
  const [duration, setDuration] = useState<number>(0);
  const hlsRef = useRef<import('hls.js').default | null>(null);
  const { startPosition = 0, onReady, onError } = options;

  // Expose a seek function so the parent modal can jump to a timestamp
  const seek = useCallback((seconds: number) => {
    const el = videoRef.current;
    if (!el) return;
    el.currentTime = Math.max(0, Math.min(seconds, el.duration || 0));
  }, []);

  useEffect(() => {
    if (!url || !videoRef.current) {
      setClipState('unsupported');
      return;
    }

    setClipState('loading');
    setDuration(0);

    const videoEl = videoRef.current;
    let destroyed = false;

    const cleanup = () => {
      destroyed = true;
      if (hlsRef.current) {
        hlsRef.current.detachMedia();
        hlsRef.current.destroy();
        hlsRef.current = null;
      }
    };

    const load = async () => {
      try {
        const Hls = (await import('hls.js')).default;

        if (destroyed) return;

        if (Hls.isSupported()) {
          const hls = new Hls({
            enableWorker: true,

            // VOD/scrub tuning — opposite priorities to live
            lowLatencyMode: false,
            startPosition,                  // Honour caller's start time
            maxBufferLength: 60,            // Large buffer — user will scrub
            maxMaxBufferLength: 180,        // Allow up to 3 min pre-buffered
            progressive: true,              // Start playing before full load
            testBandwidth: true,            // Pick best quality for connection

            // After a seek, abandon any in-flight fragment and start fresh
            startFragPrefetch: false,

            // Flush the buffer on seek instead of trying to stitch
            backBufferLength: 0,

            // Don't try to be clever about rebuffering after seek
            maxBufferHole: 0.5,

            // Give it longer to recover after a seek stall
            highBufferWatchdogPeriod: 10,
            nudgeMaxRetry: 10,
          });

          hlsRef.current = hls;
          hls.loadSource(url);
          hls.attachMedia(videoEl);

          hls.on(Hls.Events.MANIFEST_PARSED, () => {
            if (destroyed) return;
            // Don't autoplay — let the user initiate
            setClipState('ready');
          });

          hls.on(Hls.Events.LEVEL_LOADED, (_, data) => {
            if (destroyed) return;
            const d = data.details.totalduration;
            if (d > 0) {
              setDuration(d);
              onReady?.(d);
            }
          });

          hls.on(Hls.Events.ERROR, (_, data) => {
            if (destroyed || !data.fatal) return;
            const msg = `${data.type}: ${data.details}`;
            setClipState('error');
            onError?.(msg);
          });

// What HLS.js thinks it's doing
hls.on(Hls.Events.FRAG_BUFFERED, (_, data) => {
  console.log('buffered frag', data.frag.sn, 
    'start:', data.frag.start, 
    'duration:', data.frag.duration);
});

hls.on(Hls.Events.BUFFER_APPENDED, (_, data) => {
  console.log('buffer ranges:', data.timeRanges);
});

// What the video element thinks it's doing
videoRef.current.addEventListener('waiting', () => 
  console.log('video waiting at', videoRef.current.currentTime));
videoRef.current.addEventListener('stalled', () => 
  console.log('video stalled at', videoRef.current.currentTime));
videoRef.current.addEventListener('suspend', () => 
  console.log('video suspended at', videoRef.current.currentTime));

// Poll the actual buffer state every second
const interval = setInterval(() => {
  const el = videoRef.current;
  if (!el) return;
  const buffered = [];
  for (let i = 0; i < el.buffered.length; i++) {
    buffered.push(`${el.buffered.start(i).toFixed(2)}-${el.buffered.end(i).toFixed(2)}`);
  }
  console.log('t:', el.currentTime.toFixed(2), 
    'paused:', el.paused,
    'readyState:', el.readyState,
    'buffered:', buffered.join(', '));
}, 1000);

// Clean up in the useEffect return
return () => clearInterval(interval);

        } else if (videoEl.canPlayType('application/vnd.apple.mpegurl')) {
          // Native HLS — Safari / iOS
          videoEl.src = url;
          videoEl.currentTime = startPosition;

          const onLoaded = () => {
            if (destroyed) return;
            setDuration(videoEl.duration || 0);
            setClipState('ready');
            onReady?.(videoEl.duration || 0);
          };

          const onErr = () => {
            if (destroyed) return;
            setClipState('error');
            onError?.('Native HLS error');
          };

          videoEl.addEventListener('loadedmetadata', onLoaded, { once: true });
          videoEl.addEventListener('error', onErr, { once: true });

          (videoEl as any).__hlsCleanup = () => {
            videoEl.removeEventListener('loadedmetadata', onLoaded);
            videoEl.removeEventListener('error', onErr);
            videoEl.src = '';
            videoEl.load();
          };

        } else {
          setClipState('unsupported');
        }

      } catch {
        if (!destroyed) setClipState('error');
      }
    };

    load();

    return () => {
      cleanup();
      const el = videoRef.current as any;
      if (el?.__hlsCleanup) {
        el.__hlsCleanup();
        delete el.__hlsCleanup;
      }
    };
  }, [url, startPosition]);

  return { clipState, setClipState, duration, seek };
}