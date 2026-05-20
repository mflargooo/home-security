import type { Camera } from "../types/Camera";

type CameraProps = {
    camera: Camera
    onClose: (camera: Camera) => void
}

export function CameraModal({ camera, onClose } : CameraProps) {
    return (
        <div id="cover" className="w-screen h-screen bg-slate-900 opacity-70">
            <p className="text-white">
                {camera.id}, {camera.name}, {camera.rtsp_url} 
            </p>
        </div>
    )
}