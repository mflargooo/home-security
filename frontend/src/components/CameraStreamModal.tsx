import { useRef, useEffect, useState } from "react"
import type { Camera } from "../types/Camera";
import { getStreamUrl } from "../utils/api"
import Hls from "hls.js"
import { Modal } from "./Modal";
import { CameraSettingsModal } from "./CameraSettingsModal";

type CameraProps = {
    camera: Camera
    onClose: () => void
}

export function CameraStreamModal({ camera, onClose } : CameraProps) {
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const hlsRef = useRef<Hls | null>(null);
  const [streamState, setStreamState] = useState('loading'); // loading | playing | error | unsupported
  const [settingsIsOpen, setSettingsIsOpen] = useState<boolean>(false); 

  const isStreamable = camera.status === 'active';

  /*
  useEffect(() => {
    const handleKey = (e : KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    document.addEventListener('keydown', handleKey);
    return () => document.removeEventListener('keydown', handleKey);
  }, [onClose]);
  */
  streamState;
  useEffect(() => {
    if (!isStreamable || !videoRef.current) {
      setStreamState(camera.status === 'active' ? 'error' : 'unsupported');
      return;
    }

    const streamUrl = getStreamUrl(camera.id);

    const tryHLS = async (timeoutID : number) => {
      try {
        const Hls = (await import('hls.js')).default;
        const videoEl = videoRef.current;

        if (!videoEl) return;

        if (Hls.isSupported()) {
          const hls = new Hls({
            enableWorker: true,
            lowLatencyMode: true,
          });
          hlsRef.current = hls;
          hls.loadSource(streamUrl);
          hls.attachMedia(videoEl);
          hls.on(Hls.Events.MANIFEST_PARSED, () => {
            videoEl.play().catch(() => {});
            setStreamState('playing');
            clearTimeout(timeoutID)
          });
          hls.on(Hls.Events.ERROR, (_, data) => {
            if (data.fatal) setStreamState('error');
          });
        } else if (videoEl.canPlayType('application/vnd.apple.mpegurl')) {
          // Native HLS (Safari)
          videoEl.src = streamUrl;
          videoEl.addEventListener('loadedmetadata', () => {
            videoEl.play().catch(() => {});
            setStreamState('playing');
          }, { once: true });
          clearTimeout(timeoutID);
          videoEl.addEventListener('error', () => setStreamState('error'));
        } else {
          setStreamState('unsupported');
        }
      } catch {
        setStreamState('error');
      }
    };

    // Simulate brief loading then show error (no real stream in demo)
    const t = setTimeout(() => {
      setStreamState('error'); // Will show "connect your backend" message
    }, 2000);

    // Uncomment to attempt real HLS stream:
    tryHLS(t);

    return () => {
      clearTimeout(t);
      if (hlsRef.current) {
        hlsRef.current.destroy();
        hlsRef.current = null;
      }
    };
  }, [camera.id, isStreamable, camera.status]);

  return (
    <>
      <Modal className="w-full max-w-5xl" title={camera.name} onClose={() => { setSettingsIsOpen(false); onClose() }} ellipsisOnClick={() => setSettingsIsOpen(true)}>
        <div style={{ aspectRatio: '16/9' }}>
          <video 
              ref={videoRef} 
              className={`w-full h-full`} 
              autoPlay 
              playsInline
              controls
          />
        </div>
      </Modal>

      {settingsIsOpen && 
        <CameraSettingsModal camera={camera} onClose={() => setSettingsIsOpen(false)} />
      }
    </>
  )
}