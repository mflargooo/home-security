import { useEffect, useRef, useState } from 'react';

export type StreamState = 'loading' | 'playing' | 'error' | 'unsupported';

interface UseHLSStreamOptions {
  /** Called when the stream successfully starts playing */
  onPlaying?: () => void;
  /** Called on fatal error */
  onError?: (message: string) => void;
}

/**
 * Live RTSP→HLS stream loader.
 * Optimised for low-latency live feeds: small buffer, live-edge start,
 * automatic fatal-error recovery via reload.
 */
export function useHLSStream(
  videoRef: React.RefObject<HTMLVideoElement>,
  url: string | null,
  options: UseHLSStreamOptions = {}
) {
  const [streamState, setStreamState] = useState<StreamState>('loading');
  const hlsRef = useRef<import('hls.js').default | null>(null);
  const { onPlaying, onError } = options;

  useEffect(() => {
    if (!url) {
      return;
    }

    setStreamState('loading');

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

            // Live-stream tuning
            lowLatencyMode: true,
            liveSyncDurationCount: 2,       // Stay 2 segments behind live edge
            liveMaxLatencyDurationCount: 5, // Seek back if >5 segments behind
            maxBufferLength: 10,            // Keep buffer short — it's live
            maxMaxBufferLength: 20,
            startPosition: -1,              // -1 = live edge
          });

          hlsRef.current = hls;
          hls.loadSource(url);
          hls.attachMedia(videoEl);

          hls.on(Hls.Events.MANIFEST_PARSED, () => {
            if (destroyed) return;
            videoEl.play().catch(() => {});
            setStreamState('playing');
            onPlaying?.();
          });

          hls.on(Hls.Events.ERROR, (_, data) => {
            if (destroyed || !data.fatal) return;

            const msg = `${data.type}: ${data.details}`;
            if (data.type === Hls.ErrorTypes.NETWORK_ERROR) {
              // Attempt to recover network errors by restarting the stream
              hls.startLoad();
            } else {
              setStreamState('error');
              onError?.(msg);
            }
          });

        } else if (videoEl.canPlayType('application/vnd.apple.mpegurl')) {
          // Native HLS — Safari / iOS
          videoEl.src = url;

          const onLoaded = () => {
            if (destroyed) return;
            videoEl.play().catch(() => {});
            setStreamState('playing');
            onPlaying?.();
          };

          const onErr = () => {
            if (destroyed) return;
            setStreamState('error');
            onError?.('Native HLS error');
          };

          videoEl.addEventListener('loadedmetadata', onLoaded, { once: true });
          videoEl.addEventListener('error', onErr, { once: true });

          // Store cleanup on the ref so the return() can reach it
          (videoEl as any).__hlsCleanup = () => {
            videoEl.removeEventListener('loadedmetadata', onLoaded);
            videoEl.removeEventListener('error', onErr);
            videoEl.src = '';
            videoEl.load();
          };

        } else {
          setStreamState('unsupported');
        }

      } catch {
        if (!destroyed) setStreamState('error');
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
  }, [url]);

  return { streamState, setStreamState };
}