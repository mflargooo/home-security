import { useRef, useEffect, useState, type Key } from "react"
import type { Camera } from "../types/Camera";
import { getStreamUrl } from "../utils/api"
import type { Hls } from "hls.js"

import { EllipsisVertical, X } from "lucide-react"

type CameraProps = {
    camera: Camera
    onClose: () => void
}

export function CameraModal({ camera, onClose } : CameraProps) {
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const hlsRef = useRef<Hls | null>(null);
  const [streamState, setStreamState] = useState('loading'); // loading | playing | error | unsupported
  const [fullscreen, setFullscreen] = useState(false);
  const [activeTab, setActiveTab] = useState('stream');

  const isStreamable = camera.status === 'active';

  useEffect(() => {
    const handleKey = (e : KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    document.addEventListener('keydown', handleKey);
    return () => document.removeEventListener('keydown', handleKey);
  }, [onClose]);

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

  const toggleFullscreen = () => {
    const el = document.getElementById('camera-modal-video');
    if (!document.fullscreenElement) {
      el?.requestFullscreen();
      setFullscreen(true);
    } else {
      document.exitFullscreen();
      setFullscreen(false);
    }
  };

    return (
        <div 
            className="fixed inset-0 z-50 flex modal-backdrop items-center justify-center bg-slate-900/80" 
            onClick={(e) => e.target === e.currentTarget ? onClose() : (() => {})()}
        >
            <div 
                className="
                    relative inset-0 m-auto bg-slate-900 
                    border border-rounded rounded-xl overflow-hidden
                    w-full max-w-5xl shadow-2xl h-fit
                "
            >
                <div className="flex text-slate-400 p-[.375rem] lg:p-[.5rem] justify-between">
                    <button className="bg-red-700 border-red-900 border rounded-[100%] hover:bg-red-500 w-fit h-fi my-auto" onClick={() => onClose()}><X/>
                    </button>
                    <span className="font-bold text-lg md:text-xl lg:text-2xl"> {camera.name} </span>
                    <button className="hover:bg-slate-500/50 border rounded-[100%] border-none w-fit h-fit my-auto"><EllipsisVertical/>
                    </button>
                </div>
                <div className={``} style={{ aspectRatio: '16/9' }}>
                    <video 
                        ref={videoRef} 
                        className={`w-full h-full`} 
                        autoPlay 
                        playsInline
                        controls
                    />
                </div>
            </div>
        </div>
    )
}