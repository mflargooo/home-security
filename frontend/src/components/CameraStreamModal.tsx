import { useRef, useEffect, useState } from "react"
import type { Camera } from "../types/Camera";
import { getStreamUrl } from "../utils/api"
import Hls from "hls.js"
import { Modal } from "./Modal";
import { CameraSettingsModal } from "./CameraSettingsModal";
import { CameraSnapshotModal } from "./CameraSnapshotModal";
import { Plus } from "lucide-react";
import { useHLSStream } from "../hooks/useHLSStream";

type CameraProps = {
    camera: Camera
    onClose: () => void
}

export function CameraStreamModal({ camera, onClose } : CameraProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const [settingsIsOpen, setSettingsIsOpen] = useState<boolean>(false); 
  const [clipperIsOpen, setClipperIsOpen] = useState<boolean>(false); 

  const [streamUrl, setStreamUrl] = useState<string | null>(null)

  const loadStream = () => {
      setStreamUrl(getStreamUrl(camera.id) + "?t=" + new Date().getTime())
  };

  useEffect(() => {
      loadStream()
  }, [camera.rtsp_url])

  const { streamState } = useHLSStream(videoRef, streamUrl)

  return (
    <>
      {!clipperIsOpen &&
        <Modal className="w-full max-w-5xl" title={camera.name} onClose={() => { setSettingsIsOpen(false); onClose() }} ellipsisOnClick={() => setSettingsIsOpen(true)}>
          <div style={{ aspectRatio: '16/9' }}>
            <video 
                ref={videoRef} 
                className={`w-full h-full`} 
                autoPlay 
                playsInline
                muted
                controls
            />
          </div>
          <div className="flex justify-end w-full max-w-5xl h-fit ">
            <button id="create" className="text-base md:text-xl py-[.25em] pl-[.375em] pr-[.625em] m-4"
              onClick={() => {
                setClipperIsOpen(true)
              }}
            >
              <Plus className="mr-1" size={"1.25em"}/>
              <span>New Clip</span>  
            </button>
          </div>
        </Modal>
      }
      { settingsIsOpen &&
        <CameraSettingsModal camera={camera} onClose={() => setSettingsIsOpen(false)} />
      }
      { clipperIsOpen &&
        <CameraSnapshotModal camera={camera} onClose={() => { setClipperIsOpen(false); loadStream() }} />
      }
    </>
  )
}