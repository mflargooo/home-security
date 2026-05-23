import { Modal } from "./Modal";
import type { Camera } from "../types/Camera";
import { useHLSClip } from "../hooks/useHLSClip";
import { createSnapshot, deleteSnapshot, getSnapshotUrl } from "../utils/api";
import { useState, useRef, useEffect } from "react"
import { ArrowDownToLine } from "lucide-react";

type CameraSnapshotModalProps = {
    camera: Camera;
    onClose: () => void;
    hidden?: boolean
}

export function CameraSnapshotModal({ camera, onClose, hidden=false } : CameraSnapshotModalProps ) {
    const videoRef = useRef<HTMLVideoElement>(null);
    const [snapshotUrl, setSnapshotUrl] = useState<string | null>(null);
    const [sessionID, setSessionID] = useState<string | null>(null);

    useEffect(() => {
        const loadSnapshot = async () => {
            const data = await createSnapshot(camera)
            setSessionID(data.session_id)
            setSnapshotUrl(getSnapshotUrl(data.playlist_url))
        };
        loadSnapshot()
    }, [])

    const { clipState, setClipState, duration, seek } = useHLSClip(videoRef, snapshotUrl)

    return (
        <>
            <Modal hidden={hidden} className="w-full max-w-5xl" title={camera.name} onClose={() => { sessionID ? deleteSnapshot(sessionID) : {}; onClose() }}>
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
                onClick={async () => {
                    
                }}
                >
                <ArrowDownToLine className="mr-1" size={"1.25em"}/>
                <span>Save</span>  
                </button>
            </div>
            </Modal>
        </>
    )
}